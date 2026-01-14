package repository

import (
	"context"

	"github.com/SKD-fastcampus/bot-management/services/bot-mgmt-server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormTaskRepository struct {
	db         *gorm.DB
	maxRetries int
}

// NewGormTaskRepository creates a new gormTaskRepository
func NewGormTaskRepository(db *gorm.DB, maxRetries int) domain.TaskRepository {
	return &gormTaskRepository{
		db:         db,
		maxRetries: maxRetries,
	}
}

func (r *gormTaskRepository) Create(ctx context.Context, task *domain.AnalysisTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *gormTaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AnalysisTask, error) {
	var task domain.AnalysisTask
	if err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *gormTaskRepository) Update(ctx context.Context, task *domain.AnalysisTask) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *gormTaskRepository) GetPendingTasks(ctx context.Context) ([]*domain.AnalysisTask, error) {
	var tasks []*domain.AnalysisTask
	if err := r.db.WithContext(ctx).Where("status = ?", domain.TaskStatusPending).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *gormTaskRepository) GetFailedTasks(ctx context.Context) ([]*domain.AnalysisTask, error) {
	var tasks []*domain.AnalysisTask
	// RetryCount < maxRetries and Status = FAILED
	if err := r.db.WithContext(ctx).Where("status = ? AND retry_count < ?", domain.TaskStatusFailed, r.maxRetries).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *gormTaskRepository) GetRunningTasks(ctx context.Context) ([]*domain.AnalysisTask, error) {
	var tasks []*domain.AnalysisTask
	if err := r.db.WithContext(ctx).Where("status = ?", domain.TaskStatusRunning).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *gormTaskRepository) GetActiveTaskByURL(ctx context.Context, url string) (*domain.AnalysisTask, error) {
	var task domain.AnalysisTask
	// Active status: PENDING or RUNNING or (FAILED and retry_count < maxRetries)
	if err := r.db.WithContext(ctx).
		Where("url = ? AND (status IN ? OR (status = ? AND retry_count < ?))",
			url,
			[]domain.TaskStatus{domain.TaskStatusPending, domain.TaskStatusRunning},
			domain.TaskStatusFailed,
			r.maxRetries).
		First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No active task found
		}
		return nil, err
	}
	return &task, nil
}
