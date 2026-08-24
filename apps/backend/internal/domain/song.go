package domain

import (
	"time"
)

type Song struct {
	ID           string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Title        string    `gorm:"type:varchar(255)" json:"title"`
	ChannelTitle string    `gorm:"type:varchar(255)" json:"channelTitle"`
	ThumbnailURL string    `gorm:"type:text" json:"thumbnailUrl"`
	ViewCount    string    `gorm:"type:varchar(50)" json:"viewCount"`
	CountryCode  string    `gorm:"primaryKey;type:varchar(2)" json:"countryCode"`
	CategoryID   string    `gorm:"type:varchar(50)" json:"categoryId"`
	PublishedAt  time.Time `json:"publishedAt"`
	VideoType    string    `gorm:"type:varchar(50)" json:"videoType"`
}

type SongRepository interface {
	UpsertSongs(songs []Song) error
	GetTrends(countryCode string) ([]Song, error)
	GetNewReleases(countryCode string) ([]Song, error)
	GetMVs(countryCode string) ([]Song, error)
}

type YouTubeRepository interface {
	FetchTrending(countryCode string) ([]Song, error)
}

type SongUsecase interface {
	SyncTrendingMusic(countryCode string) error
	GetTrends(countryCode string) ([]Song, error)
	GetNewReleases(countryCode string) ([]Song, error)
	GetMVs(countryCode string) ([]Song, error)
}
