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

func (m *MockYouTubeRepository) FetchTrendingByCategory(countryCode, categoryID string) ([]domain.CategoryVideo, error) {
	args := m.Called(countryCode, categoryID)
	return args.Get(0).([]domain.CategoryVideo), args.Error(1)
}

func (m *MockYouTubeRepository) FetchVideoCategories(countryCode string) ([]domain.VideoCategory, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.VideoCategory), args.Error(1)
}

type MockCategoryVideoRepository struct {
	mock.Mock
}

func (m *MockCategoryVideoRepository) UpsertVideos(videos []domain.CategoryVideo) error {
	args := m.Called(videos)
	return args.Error(0)
}

func (m *MockCategoryVideoRepository) GetVideos(countryCode string) ([]domain.CategoryVideo, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.CategoryVideo), args.Error(1)
}

func (m *MockCategoryVideoRepository) GetNewVideos(countryCode string) ([]domain.CategoryVideo, error) {
	args := m.Called(countryCode)
	return args.Get(0).([]domain.CategoryVideo), args.Error(1)
}

func TestGetTrends_Success(t *testing.T) {
	mockDB := new(MockSongRepository)
	mockYT := new(MockYouTubeRepository)

	mockSongs := []domain.Song{
		{ID: "1", Title: "เพลงฮิต 1", CountryCode: "TH"},
		{ID: "2", Title: "เพลงฮิต 2", CountryCode: "TH"},
	}

	mockDB.On("GetTrends", "TH").Return(mockSongs, nil)

	uc := usecase.NewSongUsecase(mockDB, mockYT, nil)

	result, err := uc.GetTrends("TH", "")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "เพลงฮิต 1", result[0].Title)
	mockDB.AssertExpectations(t)
}

func TestGetTrends_WithCategory_Success(t *testing.T) {
	mockDB := new(MockSongRepository)
	mockYT := new(MockYouTubeRepository)
	mockCategoryRepo := new(MockCategoryVideoRepository)

	mockVideos := []domain.CategoryVideo{
		{ID: "1", Title: "Gaming Video 1", CountryCode: "TH"},
	}
	mockCategoryRepo.On("GetVideos", "TH").Return(mockVideos, nil)

	uc := usecase.NewSongUsecase(mockDB, mockYT, map[string]domain.CategoryVideoRepository{
		"20": mockCategoryRepo,
	})

	result, err := uc.GetTrends("TH", "20")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Gaming Video 1", result[0].Title)
	assert.Equal(t, "20", result[0].CategoryID)
	mockCategoryRepo.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetTrends", mock.Anything)
}

func TestGetTrends_UnknownCategory_ReturnsError(t *testing.T) {
	mockDB := new(MockSongRepository)
	mockYT := new(MockYouTubeRepository)

	uc := usecase.NewSongUsecase(mockDB, mockYT, map[string]domain.CategoryVideoRepository{})

	_, err := uc.GetTrends("TH", "999")

	assert.ErrorIs(t, err, domain.ErrUnknownCategory)
}

func TestGetNewReleases_WithCategory_Success(t *testing.T) {
	mockDB := new(MockSongRepository)
	mockYT := new(MockYouTubeRepository)
	mockCategoryRepo := new(MockCategoryVideoRepository)

	mockVideos := []domain.CategoryVideo{
		{ID: "1", Title: "New Gaming Video", CountryCode: "TH"},
	}
	mockCategoryRepo.On("GetNewVideos", "TH").Return(mockVideos, nil)

	uc := usecase.NewSongUsecase(mockDB, mockYT, map[string]domain.CategoryVideoRepository{
		"20": mockCategoryRepo,
	})

	result, err := uc.GetNewReleases("TH", "20")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "New Gaming Video", result[0].Title)
	assert.Equal(t, "20", result[0].CategoryID)
	mockCategoryRepo.AssertExpectations(t)
	mockDB.AssertNotCalled(t, "GetNewReleases", mock.Anything)
}

func TestGetNewReleases_UnknownCategory_ReturnsError(t *testing.T) {
	mockDB := new(MockSongRepository)
	mockYT := new(MockYouTubeRepository)

	uc := usecase.NewSongUsecase(mockDB, mockYT, map[string]domain.CategoryVideoRepository{})

	_, err := uc.GetNewReleases("TH", "999")

	assert.ErrorIs(t, err, domain.ErrUnknownCategory)
}
