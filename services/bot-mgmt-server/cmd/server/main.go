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

	"github.com/labstack/echo/v4"
	"github.com/SKD-fastcampus/bot-management/pkg/config"
	"github.com/SKD-fastcampus/bot-management/pkg/logger"
	httpHandler "github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/adapter/handler/http"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/adapter/repository"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	awsInfra "github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/infrastructure/aws"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/usecase"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	log := logger.DefaultZapLogger()
	defer log.Sync()
	log.Info("Starting Bot Management Server")

	// 2. Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		cfg.GetString("db.host"), 
		cfg.GetString("db.user"), 
		cfg.GetString("db.password"), 
		cfg.GetString("db.name"), 
		cfg.GetString("db.port"))
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto Migration 
	if err := db.AutoMigrate(&domain.AnalysisTask{}); err != nil {
		log.Fatal("Failed to migrate database", zap.Error(err))
	}


	// 3. AWS Config
	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithRegion(cfg.GetString("aws.region")))
	if err != nil {
		log.Fatal("Failed to load AWS config", zap.Error(err))
	}

	// 4. Infrastructure & Repositories
	taskRepo := repository.NewPostgresTaskRepository(db)
	ecsClient := awsInfra.NewECSClient(awsCfg,
		cfg.GetString("ecs.cluster"),
		cfg.GetString("ecs.task_def"),
		cfg.GetStringSlice("ecs.subnets"),
		cfg.GetString("ecs.sec_group"),
	)


	// 5. Usecase
	taskUC := usecase.NewTaskUsecase(taskRepo, ecsClient)

	// 6. Handlers
	h := httpHandler.NewTaskHandler(taskUC)

	// 7. Echo Server
	e := echo.New()
	apiGroup := e.Group("/api/v1")
	h.RegisterRoutes(apiGroup)

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
