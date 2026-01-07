package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/labstack/echo/v4"
	"github.com/semo-backend-monorepo/pkg/logger"
	httpHandler "github.com/semo-backend-monorepo/services/bot-mgmt-server/internal/adapter/handler/http"
	"github.com/semo-backend-monorepo/services/bot-mgmt-server/internal/adapter/repository"
	awsInfra "github.com/semo-backend-monorepo/services/bot-mgmt-server/internal/infrastructure/aws"
	"github.com/semo-backend-monorepo/services/bot-mgmt-server/internal/usecase"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Logger
	log := logger.DefaultZapLogger()
	defer log.Sync()
	log.Info("Starting Bot Management Server")

	// 2. Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Seoul",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto Migration (Use with caution in prod)
	// db.AutoMigrate(&domain.AnalysisTask{}) // Assuming domain is imported if we want to migrate

	// 3. AWS Config
	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatal("Failed to load AWS config", zap.Error(err))
	}

	// 4. Infrastructure & Repositories
	taskRepo := repository.NewPostgresTaskRepository(db)
	ecsClient := awsInfra.NewECSClient(awsCfg,
		os.Getenv("ECS_CLUSTER"),
		os.Getenv("ECS_TASK_DEF"),
		os.Getenv("ECS_SUBNET"),
		os.Getenv("ECS_SEC_GROUP"),
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
