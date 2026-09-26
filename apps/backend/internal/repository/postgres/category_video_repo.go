package postgres

import (
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type categoryVideoRepository struct {
	db        *gorm.DB
	tableName string
}

func NewCategoryVideoRepository(db *gorm.DB, tableName string) domain.CategoryVideoRepository {
	return &categoryVideoRepository{db: db, tableName: tableName}
}

func (r *categoryVideoRepository) UpsertVideos(videos []domain.CategoryVideo) error {
	if len(videos) == 0 {
		return nil
	}

	return r.db.Table(r.tableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "country_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "channel_title", "thumbnail_url", "view_count", "published_at"}),
	}).Create(&videos).Error
}

func (r *categoryVideoRepository) GetVideos(countryCode string) ([]domain.CategoryVideo, error) {
	var videos []domain.CategoryVideo
	err := r.db.Table(r.tableName).
		Where("country_code = ?", countryCode).
		Order("CAST(view_count AS BIGINT) DESC").
		Limit(maxListResults).
		Find(&videos).Error
	return videos, err
}

func (r *categoryVideoRepository) GetNewVideos(countryCode string) ([]domain.CategoryVideo, error) {
	var videos []domain.CategoryVideo
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	err := r.db.Table(r.tableName).
		Where("country_code = ? AND published_at >= ?", countryCode, sevenDaysAgo).
		Order("CAST(view_count AS BIGINT) DESC").
		Limit(maxListResults).
		Find(&videos).Error
	return videos, err
}

// r.tableName is safe to concatenate: database.Connect validates it against a
// regex and a reserved-name deny-list at boot.
func (r *categoryVideoRepository) GetTopAcrossCountries(countries []string, perCountryLimit int) ([]domain.CategoryVideo, error) {
	var videos []domain.CategoryVideo
	query := `
		SELECT id, title, channel_title, thumbnail_url, view_count, country_code, published_at FROM (
			SELECT *, ROW_NUMBER() OVER (
				PARTITION BY country_code ORDER BY CAST(view_count AS BIGINT) DESC
			) AS rn
			FROM ` + r.tableName + `
			WHERE country_code IN ?
		) ranked
		WHERE rn <= ?
		ORDER BY CAST(view_count AS BIGINT) DESC
	`
	err := r.db.Raw(query, countries, perCountryLimit).Scan(&videos).Error
	return videos, err
}
