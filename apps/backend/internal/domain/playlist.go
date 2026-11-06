package domain

import (
	"errors"
	"time"
)

const (
	MaxPlaylistsPerUser     = 50
	MaxItemsPerPlaylist     = 200
	MaxPlaylistNameRunes    = 60
	MaxPlaylistTitleRunes   = 200
	MaxPlaylistChannelRunes = 100
	MaxThumbnailURLBytes    = 300
	AllowedThumbnailPrefix  = "https://i.ytimg.com/"
)

var (
	ErrPlaylistNotFound         = errors.New("playlist not found")
	ErrPlaylistItemNotFound     = errors.New("playlist item not found")
	ErrPlaylistLimitReached     = errors.New("playlist limit reached")
	ErrPlaylistItemLimitReached = errors.New("playlist item limit reached")
	ErrPlaylistInvalidOrder     = errors.New("videoIds must contain every item in the playlist exactly once")
	ErrVideoAlreadyInPlaylist   = errors.New("video already in playlist")
)

type Playlist struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"-"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	ItemCount int64          `gorm:"->;-:migration" json:"itemCount"`
	Items     []PlaylistItem `gorm:"-" json:"items,omitempty"`
}

type PlaylistItem struct {
	PlaylistID   string    `gorm:"type:uuid;primaryKey" json:"-"`
	VideoID      string    `gorm:"type:varchar(11);primaryKey" json:"videoId"`
	Position     int       `gorm:"not null;index" json:"position"`
	Title        string    `gorm:"type:varchar(255)" json:"title"`
	ChannelTitle string    `gorm:"type:varchar(255)" json:"channelTitle"`
	ThumbnailURL string    `gorm:"type:text" json:"thumbnailUrl"`
	AddedAt      time.Time `json:"addedAt"`
}

type PlaylistItemInput struct {
	VideoID      string
	Title        string
	ChannelTitle string
	ThumbnailURL string
}

type PlaylistRepository interface {
	ListByUser(userID string) ([]Playlist, error)
	CreateWithLimit(userID, name string, maxPerUser int) (*Playlist, error)
	GetWithItems(userID, playlistID string) (*Playlist, error)
	Rename(userID, playlistID, name string) (*Playlist, error)
	Delete(userID, playlistID string) error
	AddItemWithLimit(userID, playlistID string, item PlaylistItem, maxItems int) (*PlaylistItem, error)
	RemoveItem(userID, playlistID, videoID string) error
	Reorder(userID, playlistID string, videoIDs []string) error
}

type PlaylistUsecase interface {
	List(userID string) ([]Playlist, error)
	Create(userID, name string) (*Playlist, error)
	Get(userID, playlistID string) (*Playlist, error)
	Rename(userID, playlistID, name string) (*Playlist, error)
	Delete(userID, playlistID string) error
	AddItem(userID, playlistID string, input PlaylistItemInput) (*PlaylistItem, error)
	RemoveItem(userID, playlistID, videoID string) error
	Reorder(userID, playlistID string, videoIDs []string) (*Playlist, error)
}
