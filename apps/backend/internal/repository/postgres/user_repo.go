package postgres

import (
	"errors"
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) UpsertByGoogleSub(identity domain.GoogleIdentity, loginAt time.Time) (*domain.User, error) {
	var user domain.User
	err := r.db.Raw(`
		INSERT INTO users (google_sub, email, name, avatar_url, created_at, last_login_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (google_sub) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			avatar_url = EXCLUDED.avatar_url,
			last_login_at = EXCLUDED.last_login_at
		RETURNING *`,
		identity.Subject, identity.Email, identity.Name, identity.Picture, loginAt, loginAt,
	).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DeleteWithData(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("id = ?", id).Take(&domain.User{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		if err != nil {
			return err
		}
		var playlistIDs []string
		if err := tx.Model(&domain.Playlist{}).Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", id).Pluck("id", &playlistIDs).Error; err != nil {
			return err
		}
		if len(playlistIDs) > 0 {
			if err := tx.Where("playlist_id IN ?", playlistIDs).Delete(&domain.PlaylistItem{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("user_id = ?", id).Delete(&domain.Playlist{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&domain.User{}).Error
	})
}
