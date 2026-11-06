package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"tunetrend-backend/internal/domain"
)

type authUsecase struct {
	verifier domain.GoogleVerifier
	tokens   domain.SessionTokens
	users    domain.UserRepository
}

func NewAuthUsecase(verifier domain.GoogleVerifier, tokens domain.SessionTokens, users domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{verifier: verifier, tokens: tokens, users: users}
}

func (u *authUsecase) LoginWithGoogle(ctx context.Context, idToken string) (*domain.AuthSession, error) {
	idToken = strings.TrimSpace(idToken)
	if idToken == "" {
		return nil, domain.ErrUnauthorized
	}
	identity, err := u.verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, err
	}
	user, err := u.users.UpsertByGoogleSub(*identity, time.Now())
	if err != nil {
		return nil, err
	}
	token, err := u.tokens.Issue(user.ID)
	if err != nil {
		return nil, err
	}
	return &domain.AuthSession{Token: token, User: user}, nil
}

func (u *authUsecase) Authenticate(token string) (string, error) {
	userID, err := u.tokens.Parse(token)
	if err != nil {
		return "", err
	}
	if !isUUID(userID) {
		return "", domain.ErrUnauthorized
	}
	if _, err := u.users.FindByID(userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrUnauthorized
		}
		return "", err
	}
	return userID, nil
}

func (u *authUsecase) GetUser(userID string) (*domain.User, error) {
	return u.users.FindByID(userID)
}

func (u *authUsecase) DeleteUser(userID string) error {
	return u.users.DeleteWithData(userID)
}
