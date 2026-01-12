package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of an analysis task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusRunning   TaskStatus = "RUNNING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusFailed    TaskStatus = "FAILED"
)

// AnalysisTask represents a smishing analysis task
type AnalysisTask struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	RequestUUID string     `gorm:"index" json:"request_uuid"` // External User/Request UUID
	ExternalID  string     `gorm:"index" json:"external_id"`  // AWS Task ARN or similar
	URL         string     `gorm:"not null" json:"url"`
	Status      TaskStatus `gorm:"default:'PENDING'" json:"status"`
	RetryCount  int        `gorm:"default:0" json:"retry_count"`
	Result      string     `gorm:"type:text" json:"result,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	Create(ctx context.Context, task *AnalysisTask) error
	GetByID(ctx context.Context, id uuid.UUID) (*AnalysisTask, error)
	Update(ctx context.Context, task *AnalysisTask) error
	GetPendingTasks(ctx context.Context) ([]*AnalysisTask, error)
	GetFailedTasks(ctx context.Context) ([]*AnalysisTask, error)
	GetRunningTasks(ctx context.Context) ([]*AnalysisTask, error)
	GetActiveTaskByURL(ctx context.Context, url string) (*AnalysisTask, error)
}

// BotExecutor defines the interface for running and checking bot tasks
type BotExecutor interface {
	RunBot(ctx context.Context, task *AnalysisTask) (string, error) // Returns external task ID (e.g., ARN)
	GetBotStatus(ctx context.Context, externalID string) (TaskStatus, error)
}
