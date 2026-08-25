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

func (r *workerLogRepository) LatestStatuses(jobNames []string, countryCode string) (map[string]string, error) {
	result := make(map[string]string, len(jobNames))
	if len(jobNames) == 0 {
		return result, nil
	}

	type row struct {
		JobName string
		Status  string
	}
	var rows []row

	err := r.db.
		Model(&domain.WorkerLog{}).
		Select("DISTINCT ON (job_name) job_name, status").
		Where("job_name IN ? AND country_code = ?", jobNames, countryCode).
		Order("job_name, started_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, rw := range rows {
		result[rw.JobName] = rw.Status
	}
	return result, nil
}

func (r *workerLogRepository) JobsWithSuccessSince(jobNames []string, countryCode string, since time.Time) (map[string]bool, error) {
	result := make(map[string]bool, len(jobNames))
	if len(jobNames) == 0 {
		return result, nil
	}

	var matched []string
	err := r.db.
		Model(&domain.WorkerLog{}).
		Distinct("job_name").
		Where("job_name IN ? AND country_code = ? AND status = ? AND started_at >= ?",
			jobNames, countryCode, domain.WorkerLogStatusSuccess, since).
		Pluck("job_name", &matched).Error
	if err != nil {
		return nil, err
	}

	for _, name := range matched {
		result[name] = true
	}
	return result, nil
}
