package postgres

import (
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
