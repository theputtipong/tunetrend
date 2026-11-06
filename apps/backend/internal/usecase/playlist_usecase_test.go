package usecase_test

import (
	"strings"
	"testing"

	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testUserID     = "11111111-1111-4111-8111-111111111111"
	testPlaylistID = "22222222-2222-4222-8222-222222222222"
)

type MockPlaylistRepository struct {
	mock.Mock
}

func (m *MockPlaylistRepository) ListByUser(userID string) ([]domain.Playlist, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Playlist), args.Error(1)
}

func (m *MockPlaylistRepository) CreateWithLimit(userID, name string, max int) (*domain.Playlist, error) {
	args := m.Called(userID, name, max)
	p, _ := args.Get(0).(*domain.Playlist)
	return p, args.Error(1)
}

func (m *MockPlaylistRepository) GetWithItems(userID, playlistID string) (*domain.Playlist, error) {
	args := m.Called(userID, playlistID)
	p, _ := args.Get(0).(*domain.Playlist)
	return p, args.Error(1)
}

func (m *MockPlaylistRepository) Rename(userID, playlistID, name string) (*domain.Playlist, error) {
	args := m.Called(userID, playlistID, name)
	p, _ := args.Get(0).(*domain.Playlist)
	return p, args.Error(1)
}

func (m *MockPlaylistRepository) Delete(userID, playlistID string) error {
	return m.Called(userID, playlistID).Error(0)
}

func (m *MockPlaylistRepository) AddItemWithLimit(userID, playlistID string, item domain.PlaylistItem, max int) (*domain.PlaylistItem, error) {
	args := m.Called(userID, playlistID, item, max)
	p, _ := args.Get(0).(*domain.PlaylistItem)
	return p, args.Error(1)
}

func (m *MockPlaylistRepository) RemoveItem(userID, playlistID, videoID string) error {
	return m.Called(userID, playlistID, videoID).Error(0)
}

func (m *MockPlaylistRepository) Reorder(userID, playlistID string, videoIDs []string) error {
	return m.Called(userID, playlistID, videoIDs).Error(0)
}

func isValidation(t *testing.T, err error) {
	t.Helper()
	var verr *domain.ValidationError
	assert.ErrorAs(t, err, &verr)
}

func TestPlaylistUsecase_CreateTrimsName(t *testing.T) {
	repo := new(MockPlaylistRepository)
	repo.On("CreateWithLimit", testUserID, "เพลงขับรถ", domain.MaxPlaylistsPerUser).
		Return(&domain.Playlist{ID: testPlaylistID, Name: "เพลงขับรถ"}, nil)

	p, err := usecase.NewPlaylistUsecase(repo).Create(testUserID, "  เพลงขับรถ \n")
	require.NoError(t, err)
	assert.Equal(t, "เพลงขับรถ", p.Name)
	repo.AssertExpectations(t)
}

func TestPL02_PlaylistUsecase_NameLimits(t *testing.T) {
	uc := usecase.NewPlaylistUsecase(new(MockPlaylistRepository))

	for _, name := range []string{"", "   ", strings.Repeat("ก", 61), "bad\x00name"} {
		_, err := uc.Create(testUserID, name)
		isValidation(t, err)
		_, err = uc.Rename(testUserID, testPlaylistID, name)
		isValidation(t, err)
	}

	repo := new(MockPlaylistRepository)
	sixty := strings.Repeat("ก", 60)
	repo.On("CreateWithLimit", testUserID, sixty, domain.MaxPlaylistsPerUser).Return(&domain.Playlist{Name: sixty}, nil)
	_, err := usecase.NewPlaylistUsecase(repo).Create(testUserID, sixty)
	assert.NoError(t, err)
}

func TestPL02_PlaylistUsecase_ItemValidation(t *testing.T) {
	uc := usecase.NewPlaylistUsecase(new(MockPlaylistRepository))
	valid := domain.PlaylistItemInput{
		VideoID:      "dQw4w9WgXcQ",
		Title:        "Song",
		ChannelTitle: "Channel",
		ThumbnailURL: "https://i.ytimg.com/vi/dQw4w9WgXcQ/hqdefault.jpg",
	}

	cases := map[string]func(in *domain.PlaylistItemInput){
		"short videoId":     func(in *domain.PlaylistItemInput) { in.VideoID = "abc" },
		"bad char videoId":  func(in *domain.PlaylistItemInput) { in.VideoID = "dQw4w9WgXc!" },
		"long videoId":      func(in *domain.PlaylistItemInput) { in.VideoID = "dQw4w9WgXcQQ" },
		"empty title":       func(in *domain.PlaylistItemInput) { in.Title = "  " },
		"long title":        func(in *domain.PlaylistItemInput) { in.Title = strings.Repeat("a", 201) },
		"long channel":      func(in *domain.PlaylistItemInput) { in.ChannelTitle = strings.Repeat("a", 101) },
		"http thumbnail":    func(in *domain.PlaylistItemInput) { in.ThumbnailURL = "http://i.ytimg.com/vi/x.jpg" },
		"foreign thumbnail": func(in *domain.PlaylistItemInput) { in.ThumbnailURL = "https://evil.com/https://i.ytimg.com/" },
		"lookalike host":    func(in *domain.PlaylistItemInput) { in.ThumbnailURL = "https://i.ytimg.com.evil.com/x.jpg" },
		"long thumbnail": func(in *domain.PlaylistItemInput) {
			in.ThumbnailURL = "https://i.ytimg.com/" + strings.Repeat("a", 300)
		},
	}
	for name, mutate := range cases {
		t.Run("PL-02 "+name, func(t *testing.T) {
			in := valid
			mutate(&in)
			_, err := uc.AddItem(testUserID, testPlaylistID, in)
			isValidation(t, err)
		})
	}

	repo := new(MockPlaylistRepository)
	repo.On("AddItemWithLimit", testUserID, testPlaylistID, mock.MatchedBy(func(it domain.PlaylistItem) bool {
		return it.VideoID == "dQw4w9WgXcQ" && it.Title == "Song" && it.ThumbnailURL == ""
	}), domain.MaxItemsPerPlaylist).Return(&domain.PlaylistItem{VideoID: "dQw4w9WgXcQ"}, nil)
	in := valid
	in.Title = " Song "
	in.ThumbnailURL = ""
	_, err := usecase.NewPlaylistUsecase(repo).AddItem(testUserID, testPlaylistID, in)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPL01_PlaylistUsecase_MalformedIDIsNotFound(t *testing.T) {
	uc := usecase.NewPlaylistUsecase(new(MockPlaylistRepository))

	_, err := uc.Get(testUserID, "not-a-uuid")
	assert.ErrorIs(t, err, domain.ErrPlaylistNotFound)
	assert.ErrorIs(t, uc.Delete(testUserID, "1 OR 1=1"), domain.ErrPlaylistNotFound)
	_, err = uc.Rename(testUserID, "x", "name")
	assert.ErrorIs(t, err, domain.ErrPlaylistNotFound)
}

func TestPL02_PlaylistUsecase_ReorderRejectsBadInput(t *testing.T) {
	uc := usecase.NewPlaylistUsecase(new(MockPlaylistRepository))

	_, err := uc.Reorder(testUserID, testPlaylistID, []string{"dQw4w9WgXcQ", "dQw4w9WgXcQ"})
	assert.ErrorIs(t, err, domain.ErrPlaylistInvalidOrder)

	_, err = uc.Reorder(testUserID, testPlaylistID, []string{"bad"})
	assert.ErrorIs(t, err, domain.ErrPlaylistInvalidOrder)

	_, err = uc.Reorder(testUserID, testPlaylistID, make([]string, domain.MaxItemsPerPlaylist+1))
	assert.ErrorIs(t, err, domain.ErrPlaylistInvalidOrder)
}

func TestPlaylistUsecase_RemoveItemBadVideoIDIsNotFound(t *testing.T) {
	uc := usecase.NewPlaylistUsecase(new(MockPlaylistRepository))
	assert.ErrorIs(t, uc.RemoveItem(testUserID, testPlaylistID, "bad"), domain.ErrPlaylistItemNotFound)
}
