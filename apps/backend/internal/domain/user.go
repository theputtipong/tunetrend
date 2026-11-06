package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnauthorized      = errors.New("unauthorized")
	ErrAuthNotConfigured = errors.New("auth not configured")
	ErrUserNotFound      = errors.New("user not found")
	ErrAuthUnavailable   = errors.New("auth_unavailable")
)

type User struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	GoogleSub   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"`
	Email       string    `gorm:"type:varchar(320)" json:"email"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	AvatarURL   string    `gorm:"type:text" json:"avatarUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	LastLoginAt time.Time `json:"lastLoginAt"`
}

type GoogleIdentity struct {
	Subject string
	Email   string
	Name    string
	Picture string
}

type GoogleVerifier interface {
	Verify(ctx context.Context, idToken string) (*GoogleIdentity, error)
}

type SessionTokens interface {
	Issue(userID string) (string, error)
	Parse(token string) (string, error)
}

type UserRepository interface {
	UpsertByGoogleSub(identity GoogleIdentity, loginAt time.Time) (*User, error)
	FindByID(id string) (*User, error)
	DeleteWithData(id string) error
}

type AuthSession struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type AuthUsecase interface {
	LoginWithGoogle(ctx context.Context, idToken string) (*AuthSession, error)
	Authenticate(token string) (string, error)
	GetUser(userID string) (*User, error)
	DeleteUser(userID string) error
}
