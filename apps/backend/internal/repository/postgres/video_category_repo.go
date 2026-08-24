package postgres

import (
	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type videoCategoryRepository struct {
	db *gorm.DB
}

func NewVideoCategoryRepository(db *gorm.DB) domain.VideoCategoryRepository {
	return &videoCategoryRepository{db: db}
}

func (r *videoCategoryRepository) UpsertCategories(categories []domain.VideoCategory) error {
	if len(categories) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "country_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "assignable"}),
	}).Create(&categories).Error
}
