package usecase_test

import (
	"testing"

	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) UpsertSongs(songs []domain.Song) error {
	args := m.Called(songs)
	return args.Error(0)
}

func (m *MockSongRepository) GetTrends(countryCode string) ([]domain.Song, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.Song), args.Error(1)
}

func (m *MockSongRepository) GetNewReleases(countryCode string) ([]domain.Song, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.Song), args.Error(1)
}

func (m *MockSongRepository) GetMVs(countryCode string) ([]domain.Song, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.Song), args.Error(1)
}

type MockYouTubeRepository struct {
	mock.Mock
}

func (m *MockYouTubeRepository) FetchTrending(countryCode string) ([]domain.Song, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.Song), args.Error(1)
}

func TestGetTrends_Success(t *testing.T) {
	mockDB := new(MockSongRepository)
	mockYT := new(MockYouTubeRepository)

	mockSongs := []domain.Song{
		{ID: "1", Title: "เพลงฮิต 1", CountryCode: "TH"},
		{ID: "2", Title: "เพลงฮิต 2", CountryCode: "TH"},
	}

	mockDB.On("GetTrends", "TH").Return(mockSongs, nil)

	uc := usecase.NewSongUsecase(mockDB, mockYT)

	result, err := uc.GetTrends("TH")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "เพลงฮิต 1", result[0].Title)
	mockDB.AssertExpectations(t)
}
