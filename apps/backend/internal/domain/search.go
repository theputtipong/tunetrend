package domain

import (
	"errors"
	"time"
)

const (
	SearchSourceChart    = "chart"
	SearchSourceSearched = "searched"
	SearchSourceYouTube  = "youtube"

	ChartTabTrending = "trending"
	ChartTabNew      = "new"

	SearchVideoRefreshJobName = "search_video_refresh"
	YouTubeQuotaAlertJobName  = "youtube_quota_alert"

	QuotaErrorExhausted   = "quota_exhausted"
	QuotaErrorUserLimit   = "user_limit"
	QuotaErrorUnavailable = "search_unavailable"
)

var (
	ErrQuotaExhausted   = errors.New(QuotaErrorExhausted)
	ErrUserLimit        = errors.New(QuotaErrorUserLimit)
	ErrQuotaUnavailable = errors.New(QuotaErrorUnavailable)
	ErrUnknownVideo     = errors.New("unknown video")
)

type QuotaError struct {
	Err       error
	RetryAt   time.Time
	Remaining int64
}

func (e *QuotaError) Error() string { return e.Err.Error() }
func (e *QuotaError) Unwrap() error { return e.Err }

type SearchVideo struct {
	ID             string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Title          string    `gorm:"type:varchar(255)" json:"title"`
	ChannelTitle   string    `gorm:"type:varchar(255)" json:"channelTitle"`
	ThumbnailURL   string    `gorm:"type:text" json:"thumbnailUrl"`
	ViewCount      string    `gorm:"type:varchar(50)" json:"viewCount"`
	PublishedAt    time.Time `json:"publishedAt"`
	VideoType      string    `gorm:"type:varchar(50)" json:"videoType"`
	CategoryID     string    `gorm:"type:varchar(10);index" json:"categoryId"`
	FirstSeenAt    time.Time `json:"firstSeenAt"`
	StatsUpdatedAt time.Time `json:"statsUpdatedAt"`
}

type SearchHit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	VideoID     string    `gorm:"type:varchar(50);index;not null" json:"videoId"`
	CountryCode string    `gorm:"type:varchar(2);not null;index:idx_search_hits_country_hit_at,priority:1" json:"countryCode"`
	QueryNorm   string    `gorm:"type:varchar(255)" json:"queryNorm"`
	HitAt       time.Time `gorm:"not null;index:idx_search_hits_country_hit_at,priority:2" json:"hitAt"`
}

type SearchItem struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	ChannelTitle string    `json:"channelTitle"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	ViewCount    string    `json:"viewCount"`
	PublishedAt  time.Time `json:"publishedAt"`
	VideoType    string    `json:"videoType"`
	CategoryID   string    `json:"categoryId"`
	Source       string    `json:"source"`
	ChartTab     string    `json:"chartTab,omitempty"`
	ChartRank    int       `json:"chartRank,omitempty"`
}

type SearchedTrend struct {
	SearchItem
	Rank       int   `json:"rank"`
	Hits       int64 `json:"hits"`
	PrevRank   *int  `json:"prevRank"`
	RankChange int   `json:"rankChange"`
	IsNew      bool  `json:"isNew"`
}

type LocalSearchResult struct {
	Items          []SearchItem `json:"items"`
	SuggestYouTube bool         `json:"suggestYouTube"`
}

type YouTubeSearchResult struct {
	Items     []SearchItem `json:"items"`
	Remaining int64        `json:"remaining"`
	Cached    bool         `json:"cached"`
}

type SearchRepository interface {
	SearchLocal(countryCode, categoryID, query string, limit int) ([]SearchItem, error)
	UpsertSearchVideos(videos []SearchVideo) error
	SearchVideoIDsWithHitsSince(since time.Time) ([]string, error)
	VideoExists(videoID string) (bool, error)
	CreateHit(hit SearchHit) error
	SearchedTrends(countryCode string, now time.Time, limit int) ([]SearchedTrend, error)
}

type YouTubeQuotaGuard interface {
	Reserve(userID string) (remaining int64, err error)
	Release(userID string)
	Remaining(userID string) (int64, error)
	RetryAt() time.Time
}

type SearchCache interface {
	Get(key string) ([]SearchItem, bool, error)
	Set(key string, items []SearchItem) error
}

type HitDeduper interface {
	FirstHit(ip, videoID string) (bool, error)
}

type SearchUsecase interface {
	SearchLocal(query, countryCode, categoryID string) (LocalSearchResult, error)
	SearchYouTube(userID, query, countryCode, categoryID string) (YouTubeSearchResult, error)
	RecordHit(ip, videoID, countryCode, query string) (bool, error)
	GetSearchedTrends(countryCode string) ([]SearchedTrend, error)
	RefreshSearchVideoStats() (int, error)
}
