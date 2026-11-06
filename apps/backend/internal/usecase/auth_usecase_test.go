package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockVerifier struct{ mock.Mock }

func (m *MockVerifier) Verify(ctx context.Context, token string) (*domain.GoogleIdentity, error) {
	args := m.Called(token)
	id, _ := args.Get(0).(*domain.GoogleIdentity)
	return id, args.Error(1)
}

type MockTokens struct{ mock.Mock }

func (m *MockTokens) Issue(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokens) Parse(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) UpsertByGoogleSub(id domain.GoogleIdentity, at time.Time) (*domain.User, error) {
	args := m.Called(id, at)
	u, _ := args.Get(0).(*domain.User)
	return u, args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	u, _ := args.Get(0).(*domain.User)
	return u, args.Error(1)
}

func (m *MockUserRepository) DeleteWithData(id string) error {
	return m.Called(id).Error(0)
}

func TestAuthUsecase_LoginIssuesToken(t *testing.T) {
	v, tok, repo := new(MockVerifier), new(MockTokens), new(MockUserRepository)
	identity := &domain.GoogleIdentity{Subject: "sub-1", Email: "a@example.com"}
	user := &domain.User{ID: testUserID, GoogleSub: "sub-1"}
	v.On("Verify", "id-token").Return(identity, nil)
	repo.On("UpsertByGoogleSub", *identity, mock.Anything).Return(user, nil)
	tok.On("Issue", testUserID).Return("session", nil)

	s, err := usecase.NewAuthUsecase(v, tok, repo).LoginWithGoogle(context.Background(), "id-token")
	require.NoError(t, err)
	assert.Equal(t, "session", s.Token)
	assert.Equal(t, user, s.User)
}

func TestAUTH01_AuthUsecase_LoginRejectsBadGoogleToken(t *testing.T) {
	v, tok, repo := new(MockVerifier), new(MockTokens), new(MockUserRepository)
	v.On("Verify", "bad").Return(nil, domain.ErrUnauthorized)

	_, err := usecase.NewAuthUsecase(v, tok, repo).LoginWithGoogle(context.Background(), "bad")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)

	_, err = usecase.NewAuthUsecase(v, tok, repo).LoginWithGoogle(context.Background(), "  ")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	repo.AssertNotCalled(t, "UpsertByGoogleSub", mock.Anything, mock.Anything)
}

func TestAUTH01_AuthUsecase_AuthenticateDeletedUser(t *testing.T) {
	v, tok, repo := new(MockVerifier), new(MockTokens), new(MockUserRepository)
	tok.On("Parse", "session").Return(testUserID, nil)
	repo.On("FindByID", testUserID).Return(nil, domain.ErrUserNotFound)

	_, err := usecase.NewAuthUsecase(v, tok, repo).Authenticate("session")
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
}

func TestAuthUsecase_AuthenticatePropagatesConfigAndDBErrors(t *testing.T) {
	v, tok, repo := new(MockVerifier), new(MockTokens), new(MockUserRepository)
	tok.On("Parse", "a").Return("", domain.ErrAuthNotConfigured)
	tok.On("Parse", "b").Return(testUserID, nil)
	repo.On("FindByID", testUserID).Return(nil, errors.New("db down"))
	uc := usecase.NewAuthUsecase(v, tok, repo)

	_, err := uc.Authenticate("a")
	assert.ErrorIs(t, err, domain.ErrAuthNotConfigured)
	_, err = uc.Authenticate("b")
	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrUnauthorized)
}
