package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"tunetrend-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/idtoken"
)

func fakeValidate(validAud string, payload *idtoken.Payload) validateFunc {
	return func(_ context.Context, token, audience string) (*idtoken.Payload, error) {
		if token != "good" {
			return nil, errors.New("idtoken: invalid token")
		}
		if audience != validAud {
			return nil, errors.New("idtoken: audience provided does not match aud claim in the JWT")
		}
		return payload, nil
	}
}

func googlePayload(iss string, emailVerified any) *idtoken.Payload {
	return &idtoken.Payload{
		Issuer:   iss,
		Audience: "android-id",
		Subject:  "sub-123",
		Claims: map[string]any{
			"email":          "a@example.com",
			"email_verified": emailVerified,
			"name":           "Alice",
			"picture":        "https://lh3.googleusercontent.com/a",
		},
	}
}

func TestGoogleVerifier_AcceptsAnyConfiguredAudience(t *testing.T) {
	v := NewGoogleVerifier("web-id, android-id ,ios-id")
	v.validate = fakeValidate("android-id", googlePayload("https://accounts.google.com", true))

	id, err := v.Verify(context.Background(), "good")
	require.NoError(t, err)
	assert.Equal(t, &domain.GoogleIdentity{
		Subject: "sub-123",
		Email:   "a@example.com",
		Name:    "Alice",
		Picture: "https://lh3.googleusercontent.com/a",
	}, id)
}

func TestGoogleVerifier_AcceptsShortIssuerAndStringEmailVerified(t *testing.T) {
	v := NewGoogleVerifier("android-id")
	v.validate = fakeValidate("android-id", googlePayload("accounts.google.com", "true"))
	_, err := v.Verify(context.Background(), "good")
	require.NoError(t, err)
}

func TestAUTH01_GoogleVerifier_Rejects(t *testing.T) {
	t.Run("AUTH-01 forged token", func(t *testing.T) {
		v := NewGoogleVerifier("web-id,android-id")
		v.validate = fakeValidate("android-id", googlePayload("https://accounts.google.com", true))
		_, err := v.Verify(context.Background(), "forged")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 wrong audience", func(t *testing.T) {
		v := NewGoogleVerifier("web-id,ios-id")
		v.validate = fakeValidate("android-id", googlePayload("https://accounts.google.com", true))
		_, err := v.Verify(context.Background(), "good")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 wrong issuer", func(t *testing.T) {
		v := NewGoogleVerifier("android-id")
		v.validate = fakeValidate("android-id", googlePayload("https://evil.example.com", true))
		_, err := v.Verify(context.Background(), "good")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 email not verified", func(t *testing.T) {
		v := NewGoogleVerifier("android-id")
		v.validate = fakeValidate("android-id", googlePayload("https://accounts.google.com", false))
		_, err := v.Verify(context.Background(), "good")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 missing subject", func(t *testing.T) {
		p := googlePayload("https://accounts.google.com", true)
		p.Subject = ""
		v := NewGoogleVerifier("android-id")
		v.validate = fakeValidate("android-id", p)
		_, err := v.Verify(context.Background(), "good")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("AUTH-01 real validator rejects garbage", func(t *testing.T) {
		v := NewGoogleVerifier("android-id")
		_, err := v.Verify(context.Background(), "garbage")
		assert.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}

func TestGoogleVerifier_NotConfigured(t *testing.T) {
	v := NewGoogleVerifier(" , ")
	assert.False(t, v.Configured())
	_, err := v.Verify(context.Background(), "good")
	assert.ErrorIs(t, err, domain.ErrAuthNotConfigured)
}

func TestGoogleVerifier_AppliesFiveSecondTimeout(t *testing.T) {
	v := NewGoogleVerifier("android-id")
	var deadline time.Time
	var hasDeadline bool
	v.validate = func(ctx context.Context, _, _ string) (*idtoken.Payload, error) {
		deadline, hasDeadline = ctx.Deadline()
		return googlePayload("https://accounts.google.com", true), nil
	}
	start := time.Now()
	_, err := v.Verify(context.Background(), "good")
	require.NoError(t, err)
	require.True(t, hasDeadline)
	assert.WithinDuration(t, start.Add(5*time.Second), deadline, time.Second)
}

func TestGoogleVerifier_TransportErrorsAreUnavailable(t *testing.T) {
	cases := map[string]error{
		"url error":        &url.Error{Op: "Get", URL: "https://www.googleapis.com/oauth2/v3/certs", Err: errors.New("dial tcp: connection refused")},
		"deadline":         context.DeadlineExceeded,
		"cert status":      errors.New("idtoken: unable to retrieve cert, got status code 500"),
		"net op error":     &net.OpError{Op: "dial", Err: errors.New("no route to host")},
		"wrapped deadline": fmt.Errorf("fetch: %w", context.DeadlineExceeded),
	}
	for name, cause := range cases {
		t.Run(name, func(t *testing.T) {
			v := NewGoogleVerifier("web-id,android-id")
			v.validate = func(context.Context, string, string) (*idtoken.Payload, error) { return nil, cause }
			_, err := v.Verify(context.Background(), "good")
			assert.ErrorIs(t, err, domain.ErrAuthUnavailable)
			assert.NotErrorIs(t, err, domain.ErrUnauthorized)
		})
	}
}

func TestGoogleVerifier_HungValidatorReturnsUnavailable(t *testing.T) {
	v := NewGoogleVerifier("android-id")
	release := make(chan struct{})
	defer close(release)
	v.validate = func(context.Context, string, string) (*idtoken.Payload, error) {
		<-release
		return nil, errors.New("late")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := v.Verify(ctx, "good")
	assert.ErrorIs(t, err, domain.ErrAuthUnavailable)
	assert.Less(t, time.Since(start), 2*time.Second)
}
