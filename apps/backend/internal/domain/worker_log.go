package domain

import (
	"time"
)

const (
	WorkerLogStatusSuccess = "SUCCESS"
	WorkerLogStatusFailed  = "FAILED"
)

type WorkerLog struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	JobName         string    `gorm:"type:varchar(100);index" json:"jobName"`
	CountryCode     string    `gorm:"type:varchar(10)" json:"countryCode"`
	Status          string    `gorm:"type:varchar(20);index" json:"status"`
	Message         string    `gorm:"type:text" json:"message"`
	StartedAt       time.Time `json:"startedAt"`
	FinishedAt      time.Time `json:"finishedAt"`
	DurationMs      int64     `json:"durationMs"`
	IntervalMinutes int       `json:"intervalMinutes"`
	CreatedAt       time.Time `gorm:"index" json:"createdAt"`
}

type WorkerLogRepository interface {
	CreateLog(logEntry WorkerLog) error
	WithLock(lockKey int64, fn func() error) (acquired bool, err error)
	DeleteOlderThan(cutoff time.Time) error
	// LatestStatuses คืนค่าสถานะล่าสุดของแต่ละ job
	LatestStatuses(jobNames []string, countryCode string) (map[string]string, error)
	// JobsWithSuccessSince คืนค่า job ที่มี log สำเร็จตั้งแต่เวลา since
	JobsWithSuccessSince(jobNames []string, countryCode string, since time.Time) (map[string]bool, error)
}
