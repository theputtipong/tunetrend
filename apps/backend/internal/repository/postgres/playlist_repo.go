package postgres

import (
	"errors"
	"strings"
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type playlistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) domain.PlaylistRepository {
	return &playlistRepository{db: db}
}

const playlistWithCount = "playlists.*, (SELECT COUNT(*) FROM playlist_items WHERE playlist_items.playlist_id = playlists.id) AS item_count"

func (r *playlistRepository) ListByUser(userID string) ([]domain.Playlist, error) {
	playlists := []domain.Playlist{}
	err := r.db.Model(&domain.Playlist{}).
		Select(playlistWithCount).
		Where("user_id = ?", userID).
		Order("updated_at DESC, id").
		Find(&playlists).Error
	return playlists, err
}

func (r *playlistRepository) CreateWithLimit(userID, name string, maxPerUser int) (*domain.Playlist, error) {
	var created *domain.Playlist
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").
			Where("id = ?", userID).Take(&domain.User{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrUserNotFound
			}
			return err
		}
		var count int64
		if err := tx.Model(&domain.Playlist{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(maxPerUser) {
			return domain.ErrPlaylistLimitReached
		}
		p := &domain.Playlist{UserID: userID, Name: name}
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		created = p
		return nil
	})
	return created, err
}

func (r *playlistRepository) GetWithItems(userID, playlistID string) (*domain.Playlist, error) {
	p, err := findOwnedPlaylist(r.db, userID, playlistID, false)
	if err != nil {
		return nil, err
	}
	items := []domain.PlaylistItem{}
	if err := r.db.Where("playlist_id = ?", p.ID).Order("position, added_at").Find(&items).Error; err != nil {
		return nil, err
	}
	p.Items = items
	p.ItemCount = int64(len(items))
	return p, nil
}

func (r *playlistRepository) Rename(userID, playlistID, name string) (*domain.Playlist, error) {
	res := r.db.Model(&domain.Playlist{}).
		Where("id = ? AND user_id = ?", playlistID, userID).
		Updates(map[string]any{"name": name, "updated_at": time.Now()})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, domain.ErrPlaylistNotFound
	}
	var p domain.Playlist
	err := r.db.Model(&domain.Playlist{}).Select(playlistWithCount).
		Where("id = ? AND user_id = ?", playlistID, userID).Take(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *playlistRepository) Delete(userID, playlistID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Where("id = ? AND user_id = ?", playlistID, userID).Delete(&domain.Playlist{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domain.ErrPlaylistNotFound
		}
		return tx.Where("playlist_id = ?", playlistID).Delete(&domain.PlaylistItem{}).Error
	})
}

func (r *playlistRepository) AddItemWithLimit(userID, playlistID string, item domain.PlaylistItem, maxItems int) (*domain.PlaylistItem, error) {
	var added *domain.PlaylistItem
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := findOwnedPlaylist(tx, userID, playlistID, true); err != nil {
			return err
		}
		var exists int64
		if err := tx.Model(&domain.PlaylistItem{}).
			Where("playlist_id = ? AND video_id = ?", playlistID, item.VideoID).Count(&exists).Error; err != nil {
			return err
		}
		if exists > 0 {
			return domain.ErrVideoAlreadyInPlaylist
		}
		var count int64
		if err := tx.Model(&domain.PlaylistItem{}).Where("playlist_id = ?", playlistID).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(maxItems) {
			return domain.ErrPlaylistItemLimitReached
		}
		now := time.Now()
		item.PlaylistID = playlistID
		item.Position = int(count)
		item.AddedAt = now
		if err := tx.Create(&item).Error; err != nil {
			return err
		}
		added = &item
		return touchPlaylist(tx, playlistID, now)
	})
	return added, err
}

func (r *playlistRepository) RemoveItem(userID, playlistID, videoID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := findOwnedPlaylist(tx, userID, playlistID, true); err != nil {
			return err
		}
		var item domain.PlaylistItem
		err := tx.Where("playlist_id = ? AND video_id = ?", playlistID, videoID).Take(&item).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrPlaylistItemNotFound
		}
		if err != nil {
			return err
		}
		if err := tx.Where("playlist_id = ? AND video_id = ?", playlistID, videoID).Delete(&domain.PlaylistItem{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&domain.PlaylistItem{}).
			Where("playlist_id = ? AND position > ?", playlistID, item.Position).
			UpdateColumn("position", gorm.Expr("position - 1")).Error; err != nil {
			return err
		}
		return touchPlaylist(tx, playlistID, time.Now())
	})
}

func (r *playlistRepository) Reorder(userID, playlistID string, videoIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if _, err := findOwnedPlaylist(tx, userID, playlistID, true); err != nil {
			return err
		}
		var current []string
		if err := tx.Model(&domain.PlaylistItem{}).Where("playlist_id = ?", playlistID).
			Pluck("video_id", &current).Error; err != nil {
			return err
		}
		if !isPermutation(current, videoIDs) {
			return domain.ErrPlaylistInvalidOrder
		}
		if len(videoIDs) == 0 {
			return nil
		}
		if err := tx.Model(&domain.PlaylistItem{}).Where("playlist_id = ?", playlistID).
			UpdateColumn("position", gorm.Expr("array_position(?::text[], video_id::text) - 1", "{"+strings.Join(videoIDs, ",")+"}")).Error; err != nil {
			return err
		}
		return touchPlaylist(tx, playlistID, time.Now())
	})
}

func findOwnedPlaylist(db *gorm.DB, userID, playlistID string, lock bool) (*domain.Playlist, error) {
	q := db.Where("id = ? AND user_id = ?", playlistID, userID)
	if lock {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	var p domain.Playlist
	err := q.Take(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPlaylistNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func touchPlaylist(tx *gorm.DB, playlistID string, at time.Time) error {
	return tx.Model(&domain.Playlist{}).Where("id = ?", playlistID).UpdateColumn("updated_at", at).Error
}

func isPermutation(current, proposed []string) bool {
	if len(current) != len(proposed) {
		return false
	}
	remaining := make(map[string]int, len(current))
	for _, id := range current {
		remaining[id]++
	}
	for _, id := range proposed {
		if remaining[id] == 0 {
			return false
		}
		remaining[id]--
	}
	return true
}
