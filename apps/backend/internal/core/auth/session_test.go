package auth

import (
	"errors"
	"testing"
	"time"

	"tunetrend-backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func TestSessionTokens_RoundTrip(t *testing.T) {
	s := NewSessionTokens(testSecret, SessionTTL)
	require.True(t, s.Configured())

	tok, err := s.Issue("user-1")
	require.NoError(t, err)

	sub, err := s.Parse(tok)
	require.NoError(t, err)
	assert.Equal(t, "user-1", sub)
}

func TestSessionTokens_ExpiresIn30Days(t *testing.T) {
	s := NewSessionTokens(testSecret, SessionTTL)
	tok, err := s.Issue("user-1")
	require.NoError(t, err)

	claims := jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tok, &claims, func(*jwt.Token) (any, error) { return []byte(testSecret), nil })
	require.NoError(t, err)
	assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), claims.ExpiresAt.Time, time.Minute)
}

func TestSessionTokens_NotConfigured(t *testing.T) {
	for _, secret := range []string{"", "short-secret", testSecret[:31]} {
		s := NewSessionTokens(secret, SessionTTL)
		assert.False(t, s.Configured())
		_, err := s.Issue("user-1")
		assert.ErrorIs(t, err, domain.ErrAuthNotConfigured)
		_, err = s.Parse("anything")
		assert.ErrorIs(t, err, domain.ErrAuthNotConfigured)
	}
}

func TestAUTH01_SessionTokens_Rejects(t *testing.T) {
	s := NewSessionTokens(testSecret, SessionTTL)

	t.Run("AUTH-01 expired token", func(t *testing.T) {
		expired := NewSessionTokens(testSecret, -time.Hour)
		tok, err := expired.Issue("user-1")
		require.NoError(t, err)
		_, err = s.Parse(tok)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 alg none", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
			Subject:   "user-1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		})
		raw, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)
		_, err = s.Parse(raw)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 other HMAC alg", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
			Subject:   "user-1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		})
		raw, err := tok.SignedString([]byte(testSecret))
		require.NoError(t, err)
		_, err = s.Parse(raw)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 wrong secret", func(t *testing.T) {
		other := NewSessionTokens("ffffffffffffffffffffffffffffffff", SessionTTL)
		tok, err := other.Issue("user-1")
		require.NoError(t, err)
		_, err = s.Parse(tok)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 missing exp", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Subject: "user-1"})
		raw, err := tok.SignedString([]byte(testSecret))
		require.NoError(t, err)
		_, err = s.Parse(raw)
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 garbage", func(t *testing.T) {
		_, err := s.Parse("not.a.jwt")
		assert.True(t, errors.Is(err, domain.ErrUnauthorized))
	})
}
