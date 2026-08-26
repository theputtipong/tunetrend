package postgres

import (
	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
)

type categoryVideoConfigRepository struct {
	db *gorm.DB
}

func NewCategoryVideoConfigRepository(db *gorm.DB) domain.CategoryVideoConfigRepository {
	return &categoryVideoConfigRepository{db: db}
}

func (r *categoryVideoConfigRepository) GetAll() ([]domain.CategoryVideoConfig, error) {
	var configs []domain.CategoryVideoConfig
	err := r.db.Order("category_id ASC").Find(&configs).Error
	return configs, err
}
