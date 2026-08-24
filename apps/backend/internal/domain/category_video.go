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
}

type CategoryVideoUsecase interface {
	SyncVideos(countryCode string) error
}

type CategoryVideoConfig struct {
	CategoryID string
	TableName  string
	Label      string
}

var CategoryVideoConfigs = []CategoryVideoConfig{
	{CategoryID: "1", TableName: "film_animation_videos", Label: "Film & Animation"},
	{CategoryID: "2", TableName: "autos_vehicles_videos", Label: "Autos & Vehicles"},
	{CategoryID: "20", TableName: "gaming_videos", Label: "Gaming"},
	{CategoryID: "24", TableName: "entertainment_videos", Label: "Entertainment"},
	{CategoryID: "25", TableName: "news_politics_videos", Label: "News & Politics"},
	{CategoryID: "27", TableName: "education_videos", Label: "Education"},
}
