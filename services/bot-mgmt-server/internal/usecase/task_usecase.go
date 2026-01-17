package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/infrastructure/firebase"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TaskUsecase interface {
	CreateTask(ctx context.Context, url, requestUUID, firebaseToken, analysisID string) (*domain.AnalysisTask, error)
	GetTaskStatus(ctx context.Context, id uuid.UUID) (*domain.AnalysisTask, error)
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus, result string) error
	RetryFailedTasks(ctx context.Context) error
	CheckRunningTasks(ctx context.Context) error
}

type taskUsecase struct {
	repo     domain.TaskRepository
	executor domain.BotExecutor
	verifier firebase.TokenVerifier
	logger   *zap.Logger
}

func NewTaskUsecase(repo domain.TaskRepository, executor domain.BotExecutor, verifier firebase.TokenVerifier, logger *zap.Logger) TaskUsecase {
	return &taskUsecase{
		repo:     repo,
		executor: executor,
		verifier: verifier,
		logger:   logger,
	}
}

func (u *taskUsecase) CreateTask(ctx context.Context, url, requestUUID, firebaseToken, analysisID string) (*domain.AnalysisTask, error) {
	// Verify Firebase Token if provided
	if firebaseToken != "" {
		if _, err := u.verifier.VerifyIDToken(ctx, firebaseToken); err != nil {
			u.logger.Warn("Invalid firebase token", zap.Error(err))
			return nil, fmt.Errorf("invalid firebase token: %w", err)
		}
	} else {
		// Enforce token presence?
		// For now, let's make it optional or strictly required as per "verify FirebaseToken" request.
		// User said "receive ... and verify", implies it is required.
		// However, I'll log warning for now or return error if strict.
		// Let's assume strict for security if token is meant to be verified.
		// BUT the user didn't say it's MANDATORY. Let's return error if empty for better security.
		return nil, fmt.Errorf("firebase token is required")
	}

	task := &domain.AnalysisTask{
		ID:            uuid.New(),
		RequestUUID:   requestUUID,
		URL:           url,
		FirebaseToken: firebaseToken,
		AnalysisID:    analysisID,
		Status:        domain.TaskStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Check if there is already an active task for this URL
	if existingTask, err := u.repo.GetActiveTaskByURL(ctx, url); err == nil && existingTask != nil {
		u.logger.Info("Returning existing active task for URL",
			zap.String("url", url),
			zap.String("task_id", existingTask.ID.String()),
			zap.String("status", string(existingTask.Status)))
		return existingTask, nil
	}

	if err := u.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Trigger Bot
	// Note: In a real system, we might want to do this asynchronously or via a queue.
	// For now, we launch it immediately.
	go func() {
		bgCtx := context.Background()
		extID, err := u.executor.RunBot(bgCtx, task)
		if err != nil {
			u.logger.Error("Failed to run bot", zap.String("task_id", task.ID.String()), zap.Error(err))
			u.UpdateTaskStatus(bgCtx, task.ID, domain.TaskStatusFailed, err.Error())
		} else {

			// Update with External ID
			task.ExternalID = extID
			task.Status = domain.TaskStatusRunning
			task.UpdatedAt = time.Now()
			u.repo.Update(bgCtx, task)
		}
	}()

	return task, nil
}

func (u *taskUsecase) GetTaskStatus(ctx context.Context, id uuid.UUID) (*domain.AnalysisTask, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *taskUsecase) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus, result string) error {
	task, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	task.Status = status
	if result != "" {
		task.Result = result
	}
	task.UpdatedAt = time.Now()

	return u.repo.Update(ctx, task)
}

func (u *taskUsecase) RetryFailedTasks(ctx context.Context) error {
	tasks, err := u.repo.GetFailedTasks(ctx)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		task.RetryCount++
		task.Status = domain.TaskStatusPending // Reset to Pending to be picked up or run immediately
		task.UpdatedAt = time.Now()

		if err := u.repo.Update(ctx, task); err != nil {
			u.logger.Error("Failed to update task retry count", zap.String("task_id", task.ID.String()), zap.Error(err))
			continue
		}

		// Launch immediately
		go func(t *domain.AnalysisTask) {
			extID, err := u.executor.RunBot(context.Background(), t)
			if err != nil {
				// If fails again, it will be marked FAILED again by the err check or subsequent check
				u.UpdateTaskStatus(context.Background(), t.ID, domain.TaskStatusFailed, err.Error())
			} else {
				t.ExternalID = extID
				t.Status = domain.TaskStatusRunning
				t.UpdatedAt = time.Now()
				u.repo.Update(context.Background(), t)
			}
		}(task)
	}
	return nil
}

func (u *taskUsecase) CheckRunningTasks(ctx context.Context) error {
	tasks, err := u.repo.GetRunningTasks(ctx)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.ExternalID == "" {
			continue
		}

		u.logger.Debug("Checking task status", zap.String("task_id", task.ID.String()), zap.String("external_id", task.ExternalID))

		status, err := u.executor.GetBotStatus(ctx, task.ExternalID)

		if err != nil {
			u.logger.Error("Failed to check status", zap.String("task_id", task.ID.String()), zap.Error(err))
			continue
		}

		if status != task.Status {
			u.logger.Info("Updating task status",
				zap.String("task_id", task.ID.String()),
				zap.String("old_status", string(task.Status)),
				zap.String("new_status", string(status)))
			task.Status = status

			task.UpdatedAt = time.Now()
			u.repo.Update(ctx, task)
		}
	}
	return nil
}
