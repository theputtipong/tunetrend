package youtube

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func videoItem(id, title, category string) string {
	return fmt.Sprintf(`{"id":%q,"snippet":{"title":%q,"channelTitle":"Ch","categoryId":%q,"publishedAt":"2026-11-01T00:00:00Z","thumbnails":{"high":{"url":"https://img/%s"}}},"statistics":{"viewCount":"123"}}`, id, title, category, id)
}

func TestSearchVideos_CallsSearchListThenVideosList(t *testing.T) {
	var paths []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		assert.Equal(t, "test-key", r.Header.Get(apiKeyHeader))
		q := r.URL.Query()
		switch r.URL.Path {
		case "/search":
			assert.Equal(t, "video", q.Get("type"))
			assert.Equal(t, "ฝนตก & rain", q.Get("q"))
			assert.Equal(t, "TH", q.Get("regionCode"))
			assert.Equal(t, "10", q.Get("videoCategoryId"))
			assert.Equal(t, "25", q.Get("maxResults"))
			_, _ = w.Write([]byte(`{"items":[{"id":{"videoId":"bbbbbbbbbbb"}},{"id":{"videoId":"aaaaaaaaaaa"}}]}`))
		case "/videos":
			assert.Equal(t, "bbbbbbbbbbb,aaaaaaaaaaa", q.Get("id"))
			assert.Contains(t, q.Get("part"), "statistics")
			_, _ = w.Write([]byte(`{"items":[` + videoItem("aaaaaaaaaaa", "Song (Official MV)", "10") + `,` + videoItem("bbbbbbbbbbb", "Song Lyrics", "10") + `]}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	repo := &youtubeRepository{apiKey: "test-key", baseURL: srv.URL}
	videos, err := repo.SearchVideos("ฝนตก & rain", "TH", "10")
	require.NoError(t, err)
	require.Len(t, videos, 2)
	assert.Equal(t, "bbbbbbbbbbb", videos[0].ID, "search ranking order preserved")
	assert.Equal(t, "Lyric", videos[0].VideoType)
	assert.Equal(t, "MV", videos[1].VideoType)
	assert.Equal(t, "123", videos[1].ViewCount)
	assert.Equal(t, "10", videos[1].CategoryID)
	assert.Equal(t, "https://img/aaaaaaaaaaa", videos[1].ThumbnailURL)
	assert.Equal(t, []string{"/search", "/videos"}, paths)
}

func TestSearchVideos_NoResultsSkipsVideosList(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer srv.Close()

	repo := &youtubeRepository{apiKey: "k", baseURL: srv.URL}
	videos, err := repo.SearchVideos("nothing", "TH", "10")
	require.NoError(t, err)
	assert.Empty(t, videos)
	assert.Equal(t, 1, calls)
}

func TestSearchVideos_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`quotaExceeded`))
	}))
	defer srv.Close()

	repo := &youtubeRepository{apiKey: "k", baseURL: srv.URL}
	_, err := repo.SearchVideos("x", "TH", "10")
	assert.ErrorContains(t, err, "403")
}

func TestFetchVideosByIDs_BatchesOf50(t *testing.T) {
	var batchSizes []int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := strings.Split(r.URL.Query().Get("id"), ",")
		batchSizes = append(batchSizes, len(ids))
		items := make([]string, len(ids))
		for i, id := range ids {
			items[i] = videoItem(id, "t", "10")
		}
		_, _ = w.Write([]byte(`{"items":[` + strings.Join(items, ",") + `]}`))
	}))
	defer srv.Close()

	ids := make([]string, 120)
	for i := range ids {
		ids[i] = fmt.Sprintf("id%09d", i)
	}
	repo := &youtubeRepository{apiKey: "k", baseURL: srv.URL}
	videos, err := repo.FetchVideosByIDs(ids)
	require.NoError(t, err)
	assert.Len(t, videos, 120)
	assert.Equal(t, []int{50, 50, 20}, batchSizes)

	videos, err = repo.FetchVideosByIDs(nil)
	require.NoError(t, err)
	assert.Empty(t, videos)
	assert.Len(t, batchSizes, 3)
}
