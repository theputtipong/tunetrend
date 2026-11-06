package auth

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tunetrend-backend/internal/domain"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

const verifyTimeout = 5 * time.Second

type validateFunc func(ctx context.Context, token, audience string) (*idtoken.Payload, error)

type GoogleVerifier struct {
	audiences []string
	validate  validateFunc
}

func NewGoogleVerifier(clientIDs string) *GoogleVerifier {
	var audiences []string
	for _, id := range strings.Split(clientIDs, ",") {
		if id = strings.TrimSpace(id); id != "" {
			audiences = append(audiences, id)
		}
	}
	return &GoogleVerifier{audiences: audiences, validate: newValidateFunc()}
}

func newValidateFunc() validateFunc {
	v, err := idtoken.NewValidator(context.Background(), option.WithHTTPClient(&http.Client{Timeout: verifyTimeout}))
	if err != nil {
		return idtoken.Validate
	}
	return v.Validate
}

func (v *GoogleVerifier) Configured() bool {
	return len(v.audiences) > 0
}

func (v *GoogleVerifier) Verify(ctx context.Context, token string) (*domain.GoogleIdentity, error) {
	if !v.Configured() {
		return nil, domain.ErrAuthNotConfigured
	}
	ctx, cancel := context.WithTimeout(ctx, verifyTimeout)
	defer cancel()

	unavailable := false
	for _, aud := range v.audiences {
		payload, err := v.validateWithContext(ctx, token, aud)
		if err != nil {
			if isTransportError(err) {
				unavailable = true
			}
			continue
		}
		return identityFromPayload(payload)
	}
	if unavailable {
		return nil, domain.ErrAuthUnavailable
	}
	return nil, domain.ErrUnauthorized
}

type validateResult struct {
	payload *idtoken.Payload
	err     error
}

func (v *GoogleVerifier) validateWithContext(ctx context.Context, token, aud string) (*idtoken.Payload, error) {
	done := make(chan validateResult, 1)
	go func() {
		p, err := v.validate(ctx, token, aud)
		done <- validateResult{payload: p, err: err}
	}()
	select {
	case r := <-done:
		return r.payload, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func isTransportError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	return strings.Contains(err.Error(), "unable to retrieve cert")
}

func identityFromPayload(p *idtoken.Payload) (*domain.GoogleIdentity, error) {
	if p == nil || p.Subject == "" {
		return nil, domain.ErrUnauthorized
	}
	if p.Issuer != "accounts.google.com" && p.Issuer != "https://accounts.google.com" {
		return nil, domain.ErrUnauthorized
	}
	if !claimTrue(p.Claims["email_verified"]) {
		return nil, domain.ErrUnauthorized
	}
	return &domain.GoogleIdentity{
		Subject: p.Subject,
		Email:   claimString(p.Claims["email"]),
		Name:    claimString(p.Claims["name"]),
		Picture: claimString(p.Claims["picture"]),
	}, nil
}

func claimTrue(v any) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == "true"
	}
	return false
}

func claimString(v any) string {
	s, _ := v.(string)
	return s
}
