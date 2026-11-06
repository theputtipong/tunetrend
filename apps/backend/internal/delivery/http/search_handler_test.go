package http_test

import (
	"encoding/json"
	"errors"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"tunetrend-backend/internal/core/ratelimit"
	"tunetrend-backend/internal/core/upstash/upstashtest"
	deliveryhttp "tunetrend-backend/internal/delivery/http"
	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubSearchUsecase struct {
	localCalls   int
	youtubeCalls int
	hitCalls     int
	localErr     error
	youtubeErr   error
	hitErr       error
	hitCounted   bool
	lastHitIP    string
}

func (s *stubSearchUsecase) SearchLocal(q, country, category string) (domain.LocalSearchResult, error) {
	s.localCalls++
	if s.localErr != nil {
		return domain.LocalSearchResult{}, s.localErr
	}
	return domain.LocalSearchResult{Items: []domain.SearchItem{{ID: "a"}}, SuggestYouTube: true}, nil
}

func (s *stubSearchUsecase) SearchYouTube(userID, q, country, category string) (domain.YouTubeSearchResult, error) {
	s.youtubeCalls++
	if s.youtubeErr != nil {
		return domain.YouTubeSearchResult{}, s.youtubeErr
	}
	return domain.YouTubeSearchResult{Items: []domain.SearchItem{{ID: "y"}}, Remaining: 7}, nil
}

func (s *stubSearchUsecase) RecordHit(ip, videoID, country, q string) (bool, error) {
	s.hitCalls++
	s.lastHitIP = ip
	return s.hitCounted, s.hitErr
}

func (s *stubSearchUsecase) GetSearchedTrends(country string) ([]domain.SearchedTrend, error) {
	return []domain.SearchedTrend{{Rank: 1}}, nil
}

func (s *stubSearchUsecase) RefreshSearchVideoStats() (int, error) { return 0, nil }

func withUser(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if userID != "" {
			c.Locals(domain.LocalsUserID, userID)
		}
		return c.Next()
	}
}

func newSearchApp(uc domain.SearchUsecase, userID string, limiter fiber.Handler) *fiber.App {
	app := fiber.New()
	h := deliveryhttp.NewSearchHandler(uc)
	if limiter == nil {
		limiter = func(c *fiber.Ctx) error { return c.Next() }
	}
	app.Get("/search", limiter, h.Search)
	app.Get("/search/youtube", withUser(userID), h.SearchYouTube)
	app.Post("/search/hit", h.RecordHit)
	app.Get("/trends/searched", h.GetSearchedTrends)
	return app
}

type searchResp struct {
	status int
	body   map[string]any
}

func searchCall(t *testing.T, app *fiber.App, method, target, body string) searchResp {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("X-Forwarded-For", "9.9.9.9")
	res, err := app.Test(req, -1)
	require.NoError(t, err)
	defer res.Body.Close()
	out := searchResp{status: res.StatusCode, body: map[string]any{}}
	_ = json.NewDecoder(res.Body).Decode(&out.body)
	return out
}

func TestSearchHandler_LocalSuccess(t *testing.T) {
	uc := &stubSearchUsecase{}
	r := searchCall(t, newSearchApp(uc, "", nil), nethttp.MethodGet, "/search?q=rain&country=TH", "")
	assert.Equal(t, 200, r.status)
	assert.Equal(t, true, r.body["success"])
	data := r.body["data"].(map[string]any)
	assert.Equal(t, true, data["suggestYouTube"])
}

func TestSearchHandler_SRCH04_ValidationErrorIs400(t *testing.T) {
	uc := &stubSearchUsecase{localErr: &domain.ValidationError{Message: "q must be 2-100 characters"}, youtubeErr: &domain.ValidationError{Message: "unsupported country"}}
	app := newSearchApp(uc, "user-1", nil)

	r := searchCall(t, app, nethttp.MethodGet, "/search?q=a&country=TH", "")
	assert.Equal(t, 400, r.status)
	assert.Equal(t, false, r.body["success"])
	assert.Equal(t, "q must be 2-100 characters", r.body["error"])

	r = searchCall(t, app, nethttp.MethodGet, "/search/youtube?q=rain&country=FR", "")
	assert.Equal(t, 400, r.status)
}

func TestSearchHandler_SRCH05_RateLimitReturns429WithoutSearching(t *testing.T) {
	srv := upstashtest.NewServer(t)
	t.Setenv("UPSTASH_REDIS_REST_URL", srv.URL)
	t.Setenv("UPSTASH_REDIS_REST_TOKEN", srv.Token())
	uc := &stubSearchUsecase{}
	app := newSearchApp(uc, "", ratelimit.New(30, time.Minute).Named("search").Middleware())

	for i := 0; i < 30; i++ {
		require.Equal(t, 200, searchCall(t, app, nethttp.MethodGet, "/search?q=rain&country=TH", "").status)
	}
	r := searchCall(t, app, nethttp.MethodGet, "/search?q=rain&country=TH", "")
	assert.Equal(t, 429, r.status)
	assert.Equal(t, 30, uc.localCalls)
	assert.Equal(t, 0, uc.youtubeCalls)
}

func TestSearchHandler_YouTubeSuccessIncludesRemaining(t *testing.T) {
	uc := &stubSearchUsecase{}
	r := searchCall(t, newSearchApp(uc, "user-1", nil), nethttp.MethodGet, "/search/youtube?q=rain&country=TH", "")
	assert.Equal(t, 200, r.status)
	data := r.body["data"].(map[string]any)
	assert.Equal(t, float64(7), data["remaining"])
	assert.Equal(t, false, data["cached"])
}

func TestSearchHandler_YouTubeQuotaErrors(t *testing.T) {
	retryAt := time.Date(2026, 11, 7, 8, 0, 0, 0, time.UTC)
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{&domain.QuotaError{Err: domain.ErrUserLimit, RetryAt: retryAt}, 429, "user_limit"},
		{&domain.QuotaError{Err: domain.ErrQuotaExhausted, RetryAt: retryAt}, 503, "quota_exhausted"},
		{&domain.QuotaError{Err: domain.ErrQuotaUnavailable, RetryAt: retryAt}, 503, "search_unavailable"},
		{errors.New("youtube search: 500"), 502, "youtube_error"},
	}
	for _, tc := range cases {
		t.Run(tc.code, func(t *testing.T) {
			uc := &stubSearchUsecase{youtubeErr: tc.err}
			r := searchCall(t, newSearchApp(uc, "user-1", nil), nethttp.MethodGet, "/search/youtube?q=rain&country=TH", "")
			assert.Equal(t, tc.status, r.status)
			assert.Equal(t, tc.code, r.body["error"])
			if tc.code != "youtube_error" {
				assert.Equal(t, "2026-11-07T08:00:00Z", r.body["retryAt"])
				assert.Equal(t, float64(0), r.body["remaining"])
			}
		})
	}
}

func TestSearchHandler_RecordHit(t *testing.T) {
	uc := &stubSearchUsecase{hitCounted: true}
	app := newSearchApp(uc, "", nil)
	r := searchCall(t, app, nethttp.MethodPost, "/search/hit", `{"videoId":"dQw4w9WgXcQ","country":"TH","q":"rain"}`)
	assert.Equal(t, 200, r.status)
	assert.Equal(t, true, r.body["data"].(map[string]any)["counted"])
	assert.Equal(t, "9.9.9.9", uc.lastHitIP)

	r = searchCall(t, app, nethttp.MethodPost, "/search/hit", `{bad json`)
	assert.Equal(t, 400, r.status)

	uc.hitErr = domain.ErrUnknownVideo
	r = searchCall(t, app, nethttp.MethodPost, "/search/hit", `{"videoId":"AAAAAAAAAAA","country":"TH"}`)
	assert.Equal(t, 404, r.status)
}

func TestSearchHandler_SearchedTrends(t *testing.T) {
	r := searchCall(t, newSearchApp(&stubSearchUsecase{}, "", nil), nethttp.MethodGet, "/trends/searched?country=TH", "")
	assert.Equal(t, 200, r.status)
	assert.Len(t, r.body["data"].([]any), 1)
}

type stubAuthUsecase struct {
	domain.AuthUsecase
	tokens map[string]string
}

func (s stubAuthUsecase) Authenticate(token string) (string, error) {
	if id, ok := s.tokens[token]; ok {
		return id, nil
	}
	return "", domain.ErrUnauthorized
}

func newRoutedSearchApp(t *testing.T, uc domain.SearchUsecase) *fiber.App {
	t.Helper()
	srv := upstashtest.NewServer(t)
	t.Setenv("UPSTASH_REDIS_REST_URL", srv.URL)
	t.Setenv("UPSTASH_REDIS_REST_TOKEN", srv.Token())
	app := fiber.New()
	authUC := stubAuthUsecase{tokens: map[string]string{"valid-token": "user-1"}}
	deliveryhttp.RegisterSearchRoutes(app, deliveryhttp.NewSearchHandler(uc), authUC, func(c *fiber.Ctx) error { return c.Next() })
	return app
}

func authedSearchCall(t *testing.T, app *fiber.App, target, token string) searchResp {
	t.Helper()
	req := httptest.NewRequest(nethttp.MethodGet, target, nil)
	req.Header.Set("X-Forwarded-For", "9.9.9.9")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := app.Test(req, -1)
	require.NoError(t, err)
	defer res.Body.Close()
	out := searchResp{status: res.StatusCode, body: map[string]any{}}
	_ = json.NewDecoder(res.Body).Decode(&out.body)
	return out
}

func TestSearchRoutes_SRCH06_YouTubeRequiresAuthThroughRealMiddleware(t *testing.T) {
	uc := &stubSearchUsecase{}
	app := newRoutedSearchApp(t, uc)

	for _, token := range []string{"", "forged-token"} {
		r := authedSearchCall(t, app, "/search/youtube?q=rain&country=TH", token)
		assert.Equal(t, 401, r.status, "token %q", token)
		assert.Equal(t, false, r.body["success"])
	}
	assert.Equal(t, 0, uc.youtubeCalls)

	r := authedSearchCall(t, app, "/search/youtube?q=rain&country=TH", "valid-token")
	assert.Equal(t, 200, r.status)
	assert.Equal(t, 1, uc.youtubeCalls)
}

func TestSearchRoutes_SRCH05_YouTubeRateLimitReturns429WithoutSearching(t *testing.T) {
	uc := &stubSearchUsecase{}
	app := newRoutedSearchApp(t, uc)

	for i := 0; i < 30; i++ {
		require.Equal(t, 200, authedSearchCall(t, app, "/search/youtube?q=rain&country=TH", "valid-token").status)
	}
	r := authedSearchCall(t, app, "/search/youtube?q=rain&country=TH", "valid-token")
	assert.Equal(t, 429, r.status)
	assert.Equal(t, 30, uc.youtubeCalls)

	require.Equal(t, 200, authedSearchCall(t, app, "/search?q=rain&country=TH", "").status, "limiters are scoped per route")
}
