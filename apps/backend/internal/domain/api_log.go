package domain

import "time"

type ApiLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RequestID    string    `gorm:"type:varchar(64);index" json:"requestId"`
	Method       string    `gorm:"type:varchar(10)" json:"method"`
	Path         string    `gorm:"type:varchar(255);index" json:"path"`
	Query        string    `gorm:"type:text" json:"query,omitempty"`
	StatusCode   int       `gorm:"index" json:"statusCode"`
	Success      bool      `gorm:"index" json:"success"`
	ErrorMessage string    `gorm:"type:text" json:"errorMessage,omitempty"`
	ClientIP     string    `gorm:"type:varchar(64)" json:"clientIp"`
	UserAgent    string    `gorm:"type:varchar(255)" json:"userAgent"`
	DurationMs   int64     `json:"durationMs"`
	CreatedAt    time.Time `gorm:"index" json:"createdAt"`
}

type ApiLogRepository interface {
	Create(logEntry *ApiLog) error
	DeleteOlderThan(cutoff time.Time) error
}
