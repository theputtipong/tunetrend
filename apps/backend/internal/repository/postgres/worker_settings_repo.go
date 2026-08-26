package postgres

import (
	"errors"
	"log"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
)

type workerSettingsRepository struct {
	db *gorm.DB
}

func NewWorkerSettingsRepository(db *gorm.DB) domain.WorkerSettingsRepository {
	return &workerSettingsRepository{db: db}
}

func (r *workerSettingsRepository) Get() (domain.WorkerSettings, error) {
	var rows []domain.WorkerSettings
	if err := r.db.Order("id ASC").Limit(2).Find(&rows).Error; err != nil {
		return domain.WorkerSettings{}, err
	}

	switch len(rows) {
	case 0:
		return domain.WorkerSettings{}, errors.New("worker_settings: no rows found (seeding did not run?)")
	case 1:
		return rows[0], nil
	default:
		log.Printf("⚠️ worker_settings มีมากกว่า 1 แถว ใช้แถว id=%d เป็นค่าที่ใช้งานจริง", rows[0].ID)
		return rows[0], nil
	}
}
