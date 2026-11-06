package usecase

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"tunetrend-backend/internal/domain"
)

var (
	videoIDRegex = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)
	uuidRegex    = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
)

func isUUID(s string) bool { return uuidRegex.MatchString(s) }

type playlistUsecase struct {
	repo domain.PlaylistRepository
}

func NewPlaylistUsecase(repo domain.PlaylistRepository) domain.PlaylistUsecase {
	return &playlistUsecase{repo: repo}
}

func (u *playlistUsecase) List(userID string) ([]domain.Playlist, error) {
	return u.repo.ListByUser(userID)
}

func (u *playlistUsecase) Create(userID, name string) (*domain.Playlist, error) {
	clean, err := cleanPlaylistName(name)
	if err != nil {
		return nil, err
	}
	return u.repo.CreateWithLimit(userID, clean, domain.MaxPlaylistsPerUser)
}

func (u *playlistUsecase) Get(userID, playlistID string) (*domain.Playlist, error) {
	if !isUUID(playlistID) {
		return nil, domain.ErrPlaylistNotFound
	}
	return u.repo.GetWithItems(userID, playlistID)
}

func (u *playlistUsecase) Rename(userID, playlistID, name string) (*domain.Playlist, error) {
	if !isUUID(playlistID) {
		return nil, domain.ErrPlaylistNotFound
	}
	clean, err := cleanPlaylistName(name)
	if err != nil {
		return nil, err
	}
	return u.repo.Rename(userID, playlistID, clean)
}

func (u *playlistUsecase) Delete(userID, playlistID string) error {
	if !isUUID(playlistID) {
		return domain.ErrPlaylistNotFound
	}
	return u.repo.Delete(userID, playlistID)
}

func (u *playlistUsecase) AddItem(userID, playlistID string, input domain.PlaylistItemInput) (*domain.PlaylistItem, error) {
	if !isUUID(playlistID) {
		return nil, domain.ErrPlaylistNotFound
	}
	item, err := cleanPlaylistItem(input)
	if err != nil {
		return nil, err
	}
	return u.repo.AddItemWithLimit(userID, playlistID, *item, domain.MaxItemsPerPlaylist)
}

func (u *playlistUsecase) RemoveItem(userID, playlistID, videoID string) error {
	if !isUUID(playlistID) {
		return domain.ErrPlaylistNotFound
	}
	if !videoIDRegex.MatchString(videoID) {
		return domain.ErrPlaylistItemNotFound
	}
	return u.repo.RemoveItem(userID, playlistID, videoID)
}

func (u *playlistUsecase) Reorder(userID, playlistID string, videoIDs []string) (*domain.Playlist, error) {
	if !isUUID(playlistID) {
		return nil, domain.ErrPlaylistNotFound
	}
	if len(videoIDs) > domain.MaxItemsPerPlaylist {
		return nil, domain.ErrPlaylistInvalidOrder
	}
	seen := make(map[string]bool, len(videoIDs))
	for _, id := range videoIDs {
		if !videoIDRegex.MatchString(id) || seen[id] {
			return nil, domain.ErrPlaylistInvalidOrder
		}
		seen[id] = true
	}
	if err := u.repo.Reorder(userID, playlistID, videoIDs); err != nil {
		return nil, err
	}
	return u.repo.GetWithItems(userID, playlistID)
}

func cleanPlaylistName(name string) (string, error) {
	clean := strings.TrimSpace(name)
	n := utf8.RuneCountInString(clean)
	if n < 1 || n > domain.MaxPlaylistNameRunes {
		return "", &domain.ValidationError{Message: "name must be 1-60 characters"}
	}
	if !utf8.ValidString(clean) || hasControlChars(clean) {
		return "", &domain.ValidationError{Message: "name contains invalid characters"}
	}
	return clean, nil
}

func cleanPlaylistItem(in domain.PlaylistItemInput) (*domain.PlaylistItem, error) {
	if !videoIDRegex.MatchString(in.VideoID) {
		return nil, &domain.ValidationError{Message: "invalid videoId"}
	}
	title := strings.TrimSpace(in.Title)
	if n := utf8.RuneCountInString(title); n < 1 || n > domain.MaxPlaylistTitleRunes || !utf8.ValidString(title) || hasControlChars(title) {
		return nil, &domain.ValidationError{Message: "title must be 1-200 characters"}
	}
	channel := strings.TrimSpace(in.ChannelTitle)
	if utf8.RuneCountInString(channel) > domain.MaxPlaylistChannelRunes || !utf8.ValidString(channel) || hasControlChars(channel) {
		return nil, &domain.ValidationError{Message: "channelTitle must be at most 100 characters"}
	}
	thumb := strings.TrimSpace(in.ThumbnailURL)
	if thumb != "" && (!strings.HasPrefix(thumb, domain.AllowedThumbnailPrefix) || len(thumb) > domain.MaxThumbnailURLBytes || hasControlChars(thumb)) {
		return nil, &domain.ValidationError{Message: "thumbnailUrl must start with https://i.ytimg.com/"}
	}
	return &domain.PlaylistItem{
		VideoID:      in.VideoID,
		Title:        title,
		ChannelTitle: channel,
		ThumbnailURL: thumb,
	}, nil
}

func hasControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
