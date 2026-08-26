package postgres

import (
	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
)

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) domain.ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) Create(msg *domain.ContactMessage) error {
	return r.db.Create(msg).Error
}

func (r *contactRepository) UpdateEmailStatus(id uint, sent bool, errMsg string) error {
	return r.db.Model(&domain.ContactMessage{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"email_sent": sent, "email_error": errMsg}).Error
}
