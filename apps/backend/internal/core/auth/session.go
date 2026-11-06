package auth

import (
	"time"

	"tunetrend-backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SessionTTL         = 30 * 24 * time.Hour
	MinSecretBytes     = 32
	sessionIssuer      = "tunetrend"
	sessionSigningName = "HS256"
)

type SessionTokens struct {
	secret []byte
	ttl    time.Duration
}

func NewSessionTokens(secret string, ttl time.Duration) *SessionTokens {
	return &SessionTokens{secret: []byte(secret), ttl: ttl}
}

func (s *SessionTokens) Configured() bool {
	return len(s.secret) >= MinSecretBytes
}

func (s *SessionTokens) Issue(userID string) (string, error) {
	if !s.Configured() {
		return "", domain.ErrAuthNotConfigured
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    sessionIssuer,
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *SessionTokens) Parse(raw string) (string, error) {
	if !s.Configured() {
		return "", domain.ErrAuthNotConfigured
	}
	claims := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(raw, &claims,
		func(*jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{sessionSigningName}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(sessionIssuer),
	)
	if err != nil || claims.Subject == "" {
		return "", domain.ErrUnauthorized
	}
	return claims.Subject, nil
}
