package usecase_test

import (
	"errors"
	"testing"
	"time"

	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func waitForAsyncMailSend(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for async email status update")
	}
}

type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) Create(msg *domain.ContactMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockContactRepository) UpdateEmailStatus(id uint, sent bool, errMsg string) error {
	args := m.Called(id, sent, errMsg)
	return args.Error(0)
}

type MockMailSender struct {
	mock.Mock
}

func (m *MockMailSender) Send(msg domain.MailMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func TestSubmitContactMessage_Success_WithEmail(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	done := make(chan struct{})
	mockRepo.On("Create", mock.Anything).Return(nil)
	mockMailer.On("Send", mock.Anything).Return(nil)
	mockRepo.On("UpdateEmailStatus", mock.Anything, true, "").Run(func(mock.Arguments) { close(done) }).Return(nil)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Name:         "Alice",
		Message:      "Hello there, interested in collaborating!",
		ContactEmail: "alice@example.com",
	})

	assert.NoError(t, err)
	waitForAsyncMailSend(t, done)
	mockRepo.AssertExpectations(t)
	mockMailer.AssertExpectations(t)
}

func TestSubmitContactMessage_Success_WithPhone(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	done := make(chan struct{})
	mockRepo.On("Create", mock.Anything).Return(nil)
	mockMailer.On("Send", mock.Anything).Return(nil)
	mockRepo.On("UpdateEmailStatus", mock.Anything, true, "").Run(func(mock.Arguments) { close(done) }).Return(nil)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:      "Hello there, interested in collaborating!",
		ContactPhone: "081-234-5678",
	})

	assert.NoError(t, err)
	waitForAsyncMailSend(t, done)
	mockRepo.AssertExpectations(t)
}

func TestSubmitContactMessage_MissingBothContacts(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message: "Hello there, interested in collaborating!",
	})

	assert.Error(t, err)
	var verr *domain.ValidationError
	assert.True(t, errors.As(err, &verr))
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestSubmitContactMessage_MessageTooShort(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:      "hi",
		ContactEmail: "alice@example.com",
	})

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestSubmitContactMessage_InvalidEmailFormat(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:      "Hello there, interested in collaborating!",
		ContactEmail: "not-an-email",
	})

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestSubmitContactMessage_InvalidThaiPhone(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:      "Hello there, interested in collaborating!",
		ContactPhone: "0123456789",
	})

	assert.Error(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestSubmitContactMessage_HoneypotFilled_SilentSuccess(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:  "",
		Honeypot: "i-am-a-bot",
	})

	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockMailer.AssertNotCalled(t, "Send", mock.Anything)
}

func TestSubmitContactMessage_EmailSendFails_StillSuccess(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	done := make(chan struct{})
	mockRepo.On("Create", mock.Anything).Return(nil)
	mockMailer.On("Send", mock.Anything).Return(errors.New("smtp: connection refused"))
	mockRepo.On("UpdateEmailStatus", mock.Anything, false, "smtp: connection refused").Run(func(mock.Arguments) { close(done) }).Return(nil)

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:      "Hello there, interested in collaborating!",
		ContactEmail: "alice@example.com",
	})

	assert.NoError(t, err)
	waitForAsyncMailSend(t, done)
	mockRepo.AssertExpectations(t)
}

func TestSubmitContactMessage_RepoCreateFails_ReturnsError(t *testing.T) {
	mockRepo := new(MockContactRepository)
	mockMailer := new(MockMailSender)

	mockRepo.On("Create", mock.Anything).Return(errors.New("db connection lost"))

	uc := usecase.NewContactUsecase(mockRepo, mockMailer)
	err := uc.SubmitContactMessage(domain.ContactSubmission{
		Message:      "Hello there, interested in collaborating!",
		ContactEmail: "alice@example.com",
	})

	assert.Error(t, err)
	mockMailer.AssertNotCalled(t, "Send", mock.Anything)
}
