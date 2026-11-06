package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"tunetrend-backend/internal/domain"
)

const (
	searchMaxResults = 25
	videosBatchSize  = 50
)

func (r *youtubeRepository) getJSON(endpoint string, params url.Values, out any) error {
	req, err := http.NewRequest(http.MethodGet, r.baseURL+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set(apiKeyHeader, r.apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("youtube api error status: %d, body: %s", resp.StatusCode, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (r *youtubeRepository) SearchVideos(query, countryCode, categoryID string) ([]domain.SearchVideo, error) {
	params := url.Values{}
	params.Set("part", "id")
	params.Set("type", "video")
	params.Set("q", query)
	params.Set("regionCode", countryCode)
	params.Set("videoCategoryId", categoryID)
	params.Set("maxResults", fmt.Sprint(searchMaxResults))

	var searchResp struct {
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}
	if err := r.getJSON("/search", params, &searchResp); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		if item.ID.VideoID != "" {
			ids = append(ids, item.ID.VideoID)
		}
	}
	if len(ids) == 0 {
		return []domain.SearchVideo{}, nil
	}

	videos, err := r.FetchVideosByIDs(ids)
	if err != nil {
		return nil, err
	}

	byID := make(map[string]domain.SearchVideo, len(videos))
	for _, v := range videos {
		byID[v.ID] = v
	}
	ordered := make([]domain.SearchVideo, 0, len(videos))
	for _, id := range ids {
		if v, ok := byID[id]; ok {
			ordered = append(ordered, v)
		}
	}
	return ordered, nil
}

func (r *youtubeRepository) FetchVideosByIDs(ids []string) ([]domain.SearchVideo, error) {
	videos := make([]domain.SearchVideo, 0, len(ids))
	for start := 0; start < len(ids); start += videosBatchSize {
		end := min(start+videosBatchSize, len(ids))

		params := url.Values{}
		params.Set("part", "snippet,statistics")
		params.Set("id", strings.Join(ids[start:end], ","))
		params.Set("maxResults", fmt.Sprint(videosBatchSize))

		var resp struct {
			Items []struct {
				ID      string `json:"id"`
				Snippet struct {
					Title        string    `json:"title"`
					ChannelTitle string    `json:"channelTitle"`
					CategoryID   string    `json:"categoryId"`
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
		if err := r.getJSON("/videos", params, &resp); err != nil {
			return nil, err
		}

		for _, item := range resp.Items {
			videos = append(videos, domain.SearchVideo{
				ID:           item.ID,
				Title:        item.Snippet.Title,
				ChannelTitle: item.Snippet.ChannelTitle,
				ThumbnailURL: item.Snippet.Thumbnails.High.URL,
				ViewCount:    item.Statistics.ViewCount,
				PublishedAt:  item.Snippet.PublishedAt,
				VideoType:    detectVideoType(item.Snippet.Title),
				CategoryID:   item.Snippet.CategoryID,
			})
		}
	}
	return videos, nil
}
