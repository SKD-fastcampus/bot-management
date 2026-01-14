package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/usecase"
	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateTask_DuplicateURL(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	mockExecutor := new(mocks.MockBotExecutor)
	logger := zap.NewNop()

	u := usecase.NewTaskUsecase(mockRepo, mockExecutor, nil, logger)

	ctx := context.Background()
	url := "http://example.com"
	reqUUID := "req-123"

	// Mock GetActiveTaskByURL to return an existing task
	existingTask := &domain.AnalysisTask{
		ID:          uuid.New(),
		RequestUUID: "req-existing",
		URL:         url,
		Status:      domain.TaskStatusRunning,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetActiveTaskByURL", ctx, url).Return(existingTask, nil)

	// Call CreateTask
	result, err := u.CreateTask(ctx, url, reqUUID, "dummy-token", "dummy-analysis-id")

	// Verify
	assert.NoError(t, err)
	assert.Equal(t, existingTask.ID, result.ID)
	assert.Equal(t, existingTask.Status, result.Status)

	// Verify that Create was NOT called
	mockRepo.AssertNotCalled(t, "Create", ctx, mock.Anything)
}

func TestCreateTask_NewURL(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	mockExecutor := new(mocks.MockBotExecutor)
	logger := zap.NewNop()

	u := usecase.NewTaskUsecase(mockRepo, mockExecutor, nil, logger)

	ctx := context.Background()
	url := "http://example.com/new"
	reqUUID := "req-new"

	// Mock GetActiveTaskByURL to return nil (no active task)
	mockRepo.On("GetActiveTaskByURL", ctx, url).Return((*domain.AnalysisTask)(nil), nil)

	// Mock Create to succeed
	mockRepo.On("Create", ctx, mock.MatchedBy(func(task *domain.AnalysisTask) bool {
		return task.URL == url && task.RequestUUID == reqUUID
	})).Return(nil)

	// Mock Bot execution
	mockExecutor.On("RunBot", mock.Anything, mock.Anything).Return("ext-new", nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	// Call CreateTask
	result, err := u.CreateTask(ctx, url, reqUUID, "dummy-token", "dummy-analysis-id")

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, url, result.URL)
	assert.Equal(t, domain.TaskStatusPending, result.Status)

	// Verify that Create WAS called
	mockRepo.AssertCalled(t, "Create", ctx, mock.Anything)

	// Wait a bit for goroutine
	time.Sleep(100 * time.Millisecond)
}
