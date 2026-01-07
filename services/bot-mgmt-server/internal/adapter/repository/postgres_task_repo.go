package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/semo-backend-monorepo/services/bot-mgmt-server/internal/domain"
	"gorm.io/gorm"
)

type postgresTaskRepository struct {
	db *gorm.DB
}

// NewPostgresTaskRepository creates a new postgresTaskRepository
func NewPostgresTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) Create(ctx context.Context, task *domain.AnalysisTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *postgresTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AnalysisTask, error) {
	var task domain.AnalysisTask
	if err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *postgresTaskRepository) Update(ctx context.Context, task *domain.AnalysisTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *postgresTaskRepository) GetPendingTasks(ctx context.Context) ([]*domain.AnalysisTask, error) {
	var tasks []*domain.AnalysisTask
	if err := r.db.WithContext(ctx).Where("status = ?", domain.TaskStatusPending).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *postgresTaskRepository) GetFailedTasks(ctx context.Context) ([]*domain.AnalysisTask, error) {
	var tasks []*domain.AnalysisTask
	// RetryCount < 3 and Status = FAILED
	if err := r.db.WithContext(ctx).Where("status = ? AND retry_count < ?", domain.TaskStatusFailed, 3).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *postgresTaskRepository) GetRunningTasks(ctx context.Context) ([]*domain.AnalysisTask, error) {
	var tasks []*domain.AnalysisTask
	if err := r.db.WithContext(ctx).Where("status = ?", domain.TaskStatusRunning).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

