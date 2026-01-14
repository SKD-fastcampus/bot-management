package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/SKD-fastcampus/bot-management/pkg/config"
	"github.com/SKD-fastcampus/bot-management/pkg/logger"
	httpHandler "github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/adapter/handler/http"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/adapter/repository"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	awsInfra "github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/infrastructure/aws"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/infrastructure/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/usecase"

	_ "github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"

	"go.uber.org/zap"
)

// @title Bot Management Server API
// @version 1.0
// @description API for managing smishing analysis bots.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// 1-1. Config Load
	cfg, err := config.Load("bot-mgmt-server")
	if err != nil {
		// Just log error but continue if env vars are set?
		// Or panic? strict startup is better.
		// Note: pkg/config attempts to load file, if not found returns error.
		// If we want to support ONLY env vars, we might need to adjust pkg/config or handle error gracefully.
		// For now, let's panic as per geo service example.
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 1-2. Logger
	logConfig := logger.Config{
		Level:       cfg.GetString("logger.level"),
		Format:      cfg.GetString("logger.format"),
		Development: cfg.GetString("app.env") == "dev",
	}
	if logConfig.Level == "" {
		logConfig.Level = "info" // Default
	}
	if logConfig.Format == "" {
		logConfig.Format = "json" // Default
	}

	log, err := logger.NewZapLogger(logConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer log.Sync()

	log.Info("Starting Bot Management Server")

	// 2. Database
	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto Migration
	if err := database.AutoMigrate(&domain.AnalysisTask{}); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}

	// 3. AWS Config
	awsOpts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(cfg.GetString("aws.region")),
	}
	if profile := cfg.GetString("aws.profile"); profile != "" {
		awsOpts = append(awsOpts, awsConfig.WithSharedConfigProfile(profile))
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), awsOpts...)
	if err != nil {
		log.Fatal("Failed to load AWS config", zap.Error(err))
	}

	// 4. Infrastructure & Repositories
	maxRetries := cfg.GetInt("task.max_retries")
	if maxRetries == 0 {
		maxRetries = 3 // Default
	}

	taskRepo := repository.NewGormTaskRepository(database, maxRetries)
	ecsClient := awsInfra.NewECSClient(awsCfg,
		cfg.GetString("ecs.cluster"),
		cfg.GetString("ecs.task_def"),
		cfg.GetString("ecs.container_name"),
		cfg.GetStringSlice("ecs.subnets"),
		cfg.GetString("ecs.sec_group"),
		log,
	)

	// 5. Usecase
	taskUC := usecase.NewTaskUsecase(taskRepo, ecsClient, log)

	// 6. Handlers
	h := httpHandler.NewTaskHandler(taskUC)

	// 7. Echo Server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	apiGroup := e.Group("/api/v1")
	h.RegisterRoutes(apiGroup)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 8. Background Workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Retry Worker
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := taskUC.RetryFailedTasks(ctx); err != nil {
					log.Error("Failed to retry tasks", zap.Error(err))
				}
			}
		}
	}()

	// Polling Worker
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := taskUC.CheckRunningTasks(ctx); err != nil {
					log.Error("Failed to check running tasks", zap.Error(err))
				}
			}
		}
	}()

	// 9. Start Server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatal("Shutting down server", zap.Error(err))
		}
	}()

	// 10. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := e.Shutdown(ctxShutdown); err != nil {
		e.Logger.Fatal(err)
	}
}
