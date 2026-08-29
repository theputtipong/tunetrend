package domain

import (
	"errors"
	"time"
)

var ErrUnknownCategory = errors.New("unknown category")

type CategoryVideo struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Title        string    `gorm:"type:varchar(255)" json:"title"`
	ChannelTitle string    `gorm:"type:varchar(255)" json:"channelTitle"`
	ThumbnailURL string    `gorm:"type:text" json:"thumbnailUrl"`
	ViewCount    string    `gorm:"type:varchar(50)" json:"viewCount"`
	CountryCode  string    `gorm:"primaryKey;type:varchar(2)" json:"countryCode"`
	PublishedAt  time.Time `json:"publishedAt"`
}

type CategoryVideoRepository interface {
	UpsertVideos(videos []CategoryVideo) error
	GetVideos(countryCode string) ([]CategoryVideo, error)
	GetNewVideos(countryCode string) ([]CategoryVideo, error)
	// GetTopAcrossCountries คืนค่าวิดีโอ top-N ต่อประเทศ (ตาม perCountryLimit)
	// รวมทุกประเทศใน countries เข้าด้วยกัน เรียงจากยอดวิวมากไปน้อย
	GetTopAcrossCountries(countries []string, perCountryLimit int) ([]CategoryVideo, error)
}

type CategoryVideoUsecase interface {
	SyncVideos(countryCode string) error
}
