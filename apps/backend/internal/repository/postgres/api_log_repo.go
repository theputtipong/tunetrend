package postgres

import (
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
)

type apiLogRepository struct {
	db *gorm.DB
}

func NewApiLogRepository(db *gorm.DB) domain.ApiLogRepository {
	return &apiLogRepository{db: db}
}

func (r *apiLogRepository) Create(logEntry *domain.ApiLog) error {
	return r.db.Create(logEntry).Error
}

func (r *apiLogRepository) DeleteOlderThan(cutoff time.Time) error {
	return r.db.Where("created_at < ?", cutoff).Delete(&domain.ApiLog{}).Error
}
