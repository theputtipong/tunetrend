package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"tunetrend-backend/internal/domain"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

const apiKeyHeader = "X-goog-api-key"

type youtubeRepository struct {
	apiKey string
}

func NewYouTubeRepository() domain.YouTubeRepository {
	return &youtubeRepository{
		apiKey: os.Getenv("YOUTUBE_API_KEY"),
	}
}

func (r *youtubeRepository) FetchTrending(countryCode string) ([]domain.Song, error) {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics,id&chart=mostPopular&videoCategoryId=10&regionCode=%s&maxResults=50",
		countryCode,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, r.apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtube api error status: %d, body: %s", resp.StatusCode, string(body))
	}

	var ytResp struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title        string    `json:"title"`
				ChannelTitle string    `json:"channelTitle"`
				CategoryId   string    `json:"categoryId"`
				PublishedAt  time.Time `json:"publishedAt"`
				Thumbnails   struct {
					High struct {
						URL string `json:"url"`
					} `json:"high"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			Statistics struct {
				ViewCount string `json:"viewCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		return nil, err
	}

	var songs []domain.Song
	for _, item := range ytResp.Items {
		songs = append(songs, domain.Song{
			ID:           item.ID,
			Title:        item.Snippet.Title,
			ChannelTitle: item.Snippet.ChannelTitle,
			ThumbnailURL: item.Snippet.Thumbnails.High.URL,
			ViewCount:    item.Statistics.ViewCount,
			CountryCode:  countryCode,
			CategoryID:   item.Snippet.CategoryId,
			PublishedAt:  item.Snippet.PublishedAt,
			VideoType:    detectVideoType(item.Snippet.Title),
		})
	}

	return songs, nil
}

func (r *youtubeRepository) FetchVideoCategories(countryCode string) ([]domain.VideoCategory, error) {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videoCategories?part=snippet&regionCode=%s",
		countryCode,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, r.apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtube api error status: %d, body: %s", resp.StatusCode, string(body))
	}

	var ytResp struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title      string `json:"title"`
				Assignable bool   `json:"assignable"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		return nil, err
	}

	var categories []domain.VideoCategory
	for _, item := range ytResp.Items {
		categories = append(categories, domain.VideoCategory{
			ID:          item.ID,
			CountryCode: countryCode,
			Title:       item.Snippet.Title,
			Assignable:  item.Snippet.Assignable,
		})
	}

	return categories, nil
}

func (r *youtubeRepository) FetchTrendingByCategory(countryCode, categoryID string) ([]domain.CategoryVideo, error) {
	url := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics,id&chart=mostPopular&videoCategoryId=%s&regionCode=%s&maxResults=50",
		categoryID, countryCode,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, r.apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("youtube api error status: %d, body: %s", resp.StatusCode, string(body))
	}

	var ytResp struct {
		Items []struct {
			ID      string `json:"id"`
			Snippet struct {
				Title        string    `json:"title"`
				ChannelTitle string    `json:"channelTitle"`
				PublishedAt  time.Time `json:"publishedAt"`
				Thumbnails   struct {
					High struct {
						URL string `json:"url"`
					} `json:"high"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			Statistics struct {
				ViewCount string `json:"viewCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ytResp); err != nil {
		return nil, err
	}

	var videos []domain.CategoryVideo
	for _, item := range ytResp.Items {
		videos = append(videos, domain.CategoryVideo{
			ID:           item.ID,
			Title:        item.Snippet.Title,
			ChannelTitle: item.Snippet.ChannelTitle,
			ThumbnailURL: item.Snippet.Thumbnails.High.URL,
			ViewCount:    item.Statistics.ViewCount,
			CountryCode:  countryCode,
			PublishedAt:  item.Snippet.PublishedAt,
		})
	}

	return videos, nil
}

func detectVideoType(title string) string {
	tokens := strings.FieldsFunc(strings.ToLower(title), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	hasWord := func(word string) bool {
		for _, tok := range tokens {
			if tok == word {
				return true
			}
		}
		return false
	}
	hasAdjacent := func(a, b string) bool {
		for i := 0; i+1 < len(tokens); i++ {
			if tokens[i] == a && tokens[i+1] == b {
				return true
			}
		}
		return false
	}

	switch {
	case hasWord("mv"), hasAdjacent("m", "v"), hasAdjacent("music", "video"), hasAdjacent("official", "video"):
		return "MV"
	case hasWord("lyric"), hasWord("lyrics"):
		return "Lyric"
	case hasWord("audio"):
		return "Audio Track"
	case hasWord("cover"):
		return "Cover"
	case hasWord("live"), hasWord("performance"), hasWord("concert"):
		return "Live Performance"
	default:
		return "General"
	}
}
