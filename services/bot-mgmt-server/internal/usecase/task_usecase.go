package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/semo-backend-monorepo/services/bot-mgmt-server/internal/domain"
)

type TaskUsecase interface {
	CreateTask(ctx context.Context, url, requestUUID string) (*domain.AnalysisTask, error)
	GetTaskStatus(ctx context.Context, id uuid.UUID) (*domain.AnalysisTask, error)
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus, result string) error
	RetryFailedTasks(ctx context.Context) error
	CheckRunningTasks(ctx context.Context) error
}

type taskUsecase struct {
	repo     domain.TaskRepository
	executor domain.BotExecutor
}

func NewTaskUsecase(repo domain.TaskRepository, executor domain.BotExecutor) TaskUsecase {
	return &taskUsecase{
		repo:     repo,
		executor: executor,
	}
}

func (u *taskUsecase) CreateTask(ctx context.Context, url, requestUUID string) (*domain.AnalysisTask, error) {
	task := &domain.AnalysisTask{
		ID:          uuid.New(),
		RequestUUID: requestUUID,
		URL:         url,
		Status:      domain.TaskStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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
			fmt.Printf("Failed to run bot for task %s: %v\n", task.ID, err)
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
			fmt.Printf("Failed to update task retry count %s: %v\n", task.ID, err)
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
		
		status, err := u.executor.GetBotStatus(ctx, task.ExternalID)
		if err != nil {
			fmt.Printf("Failed to check status for task %s: %v\n", task.ID, err)
			continue
		}

		if status != task.Status {
			fmt.Printf("Updating status for task %s: %s -> %s\n", task.ID, task.Status, status)
			task.Status = status
			task.UpdatedAt = time.Now()
			u.repo.Update(ctx, task)
		}
	}
	return nil
}

