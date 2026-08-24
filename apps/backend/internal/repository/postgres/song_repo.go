package postgres

import (
	"time"
	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) domain.SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) UpsertSongs(songs []domain.Song) error {
	if len(songs) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "country_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"view_count", "title", "thumbnail_url", "channel_title", "category_id", "published_at", "video_type"}),
	}).Create(&songs).Error
}

func (r *songRepository) GetTrends(countryCode string) ([]domain.Song, error) {
	var songs []domain.Song
	err := r.db.Where("country_code = ?", countryCode).
		Order("CAST(view_count AS BIGINT) DESC").
		Limit(50).Find(&songs).Error
	return songs, err
}

func (r *songRepository) GetNewReleases(countryCode string) ([]domain.Song, error) {
	var songs []domain.Song
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	err := r.db.Where("country_code = ? AND published_at >= ?", countryCode, sevenDaysAgo).
		Order("CAST(view_count AS BIGINT) DESC").
		Limit(50).
		Find(&songs).Error
	return songs, err
}

func (r *songRepository) GetMVs(countryCode string) ([]domain.Song, error) {
	var songs []domain.Song
	err := r.db.Where("country_code = ? AND video_type = ?", countryCode, "MV").
		Order("CAST(view_count AS BIGINT) DESC").
		Limit(50).
		Find(&songs).Error
	return songs, err
}
