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

// UpsertCategoriesSetActive upsert title/assignable/is_active/deactivated_reason ทั้งหมด —
// ใช้ตอนที่ usecase ตัดสินใจ is_active ของแถวนั้นแล้วแน่ชัด (ไม่ใช่ grace period/manual override)
func (r *videoCategoryRepository) UpsertCategoriesSetActive(categories []domain.VideoCategory) error {
	if len(categories) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "country_code"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "assignable", "is_active", "deactivated_reason"}),
	}).Create(&categories).Error
}

// UpsertCategoriesPreserveActive upsert แค่ title/assignable — ไม่แตะ is_active/deactivated_reason เลย
// เพราะแถวนี้ยังอยู่ใน grace period ของ fetch failure หรือถูก manual-override ไว้
func (r *videoCategoryRepository) UpsertCategoriesPreserveActive(categories []domain.VideoCategory) error {
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

func (r *videoCategoryRepository) GetDeactivatedReasons(countryCode string, ids []string) (map[string]string, error) {
	result := make(map[string]string, len(ids))
	if len(ids) == 0 {
		return result, nil
	}

	type row struct {
		ID                string
		DeactivatedReason string
	}
	var rows []row

	err := r.db.Model(&domain.VideoCategory{}).
		Select("id, deactivated_reason").
		Where("country_code = ? AND id IN ?", countryCode, ids).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, rw := range rows {
		result[rw.ID] = rw.DeactivatedReason
	}
	return result, nil
}
