package postgres

import (
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
)

type workerLogRepository struct {
	db *gorm.DB
}

func NewWorkerLogRepository(db *gorm.DB) domain.WorkerLogRepository {
	return &workerLogRepository{db: db}
}

func (r *workerLogRepository) CreateLog(logEntry domain.WorkerLog) error {
	return r.db.Create(&logEntry).Error
}

func (r *workerLogRepository) WithLock(lockKey int64, fn func() error) (bool, error) {
	acquired := false

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Raw("SELECT pg_try_advisory_xact_lock(?)", lockKey).Scan(&acquired).Error; err != nil {
			return err
		}
		if !acquired {
			return nil
		}
		return fn()
	})

	return acquired, err
}

func (r *workerLogRepository) DeleteOlderThan(cutoff time.Time) error {
	return r.db.Where("created_at < ?", cutoff).Delete(&domain.WorkerLog{}).Error
}
