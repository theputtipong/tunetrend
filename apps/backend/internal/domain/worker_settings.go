package domain

type WorkerSettings struct {
	ID                               uint   `gorm:"primaryKey" json:"id"`
	Countries                        string `gorm:"type:varchar(255);not null" json:"countries"`
	MusicSyncIntervalMinutes         int    `gorm:"not null" json:"musicSyncIntervalMinutes"`
	CategoryVideoSyncIntervalMinutes int    `gorm:"not null" json:"categoryVideoSyncIntervalMinutes"`
	CategoryResumeIntervalMinutes    int    `gorm:"not null" json:"categoryResumeIntervalMinutes"`
	CategorySyncIntervalMinutes      int    `gorm:"not null" json:"categorySyncIntervalMinutes"`
	CategoryFetchFailureGraceMinutes int    `gorm:"not null" json:"categoryFetchFailureGraceMinutes"`
}

type WorkerSettingsRepository interface {
	Get() (WorkerSettings, error)
}
