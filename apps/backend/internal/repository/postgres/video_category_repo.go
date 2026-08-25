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

// UpsertCategories อัปเดตแค่ title/assignable จาก YouTube เท่านั้น — ไม่แตะ is_active/note
// เพราะสอง field นี้เป็นค่าที่แอดมินตั้งเอง ไม่ใช่ค่าที่ sync มาจาก YouTube
func (r *videoCategoryRepository) UpsertCategories(categories []domain.VideoCategory) error {
	if len(categories) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "country_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "assignable"}),
	}).Create(&categories).Error
}

func (r *videoCategoryRepository) GetActiveCategories(countryCode string) ([]domain.VideoCategory, error) {
	var categories []domain.VideoCategory
	err := r.db.Where("country_code = ? AND is_active = ?", countryCode, true).
		Order("title ASC").
		Find(&categories).Error
	return categories, err
}
