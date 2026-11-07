package usecase_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"

	"tunetrend-backend/internal/core/searchcache"
	"tunetrend-backend/internal/core/upstash"
	"tunetrend-backend/internal/core/upstash/upstashtest"
	"tunetrend-backend/internal/core/ytquota"
	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) SearchLocal(countryCode, categoryID, query string, limit int) ([]domain.SearchItem, error) {
	args := m.Called(countryCode, categoryID, query, limit)
	return args.Get(0).([]domain.SearchItem), args.Error(1)
}

func (m *MockSearchRepository) UpsertSearchVideos(videos []domain.SearchVideo) error {
	return m.Called(videos).Error(0)
}

func (m *MockSearchRepository) SearchVideoIDsWithHitsSince(since time.Time) ([]string, error) {
	args := m.Called(since)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockSearchRepository) VideoExists(videoID string) (bool, error) {
	args := m.Called(videoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSearchRepository) CreateHit(hit domain.SearchHit) error {
	return m.Called(hit).Error(0)
}

func (m *MockSearchRepository) CreateQueryHit(hit domain.SearchQueryHit) error {
	return m.Called(hit).Error(0)
}

func (m *MockSearchRepository) SaveQueryResults(countryCode, queryNorm string, items []domain.SearchItem, at time.Time) error {
	return m.Called(countryCode, queryNorm, items, at).Error(0)
}

func (m *MockSearchRepository) SearchedQueries(countryCode string, since time.Time, minVisitors, limit int) ([]domain.SearchedQuery, error) {
	args := m.Called(countryCode, since, minVisitors, limit)
	return args.Get(0).([]domain.SearchedQuery), args.Error(1)
}

func (m *MockSearchRepository) DeleteQueryHitsBefore(before time.Time) (int64, error) {
	args := m.Called(before)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSearchRepository) SearchedTrends(countryCode string, now time.Time, limit int) ([]domain.SearchedTrend, error) {
	args := m.Called(countryCode, now, limit)
	return args.Get(0).([]domain.SearchedTrend), args.Error(1)
}

type stubCategoryRepo struct {
	domain.VideoCategoryRepository
	active map[string][]domain.VideoCategory
}

func (s *stubCategoryRepo) GetActiveCategories(countryCode string) ([]domain.VideoCategory, error) {
	return s.active[countryCode], nil
}

type recordingMailer struct{ sent []domain.MailMessage }

func (m *recordingMailer) Send(msg domain.MailMessage) error {
	m.sent = append(m.sent, msg)
	return nil
}

type searchHarness struct {
	srv    *upstashtest.Server
	repo   *MockSearchRepository
	yt     *MockYouTubeRepository
	mailer *recordingMailer
	uc     domain.SearchUsecase
	now    time.Time
}

var fixedNow = time.Date(2026, 11, 6, 18, 0, 0, 0, time.UTC)

func newSearchHarness(t *testing.T) *searchHarness {
	h := &searchHarness{
		srv:    upstashtest.NewServer(t),
		repo:   new(MockSearchRepository),
		yt:     new(MockYouTubeRepository),
		mailer: &recordingMailer{},
		now:    fixedNow,
	}
	h.repo.On("CreateQueryHit", mock.Anything).Return(nil).Maybe()
	h.repo.On("SaveQueryResults", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	h.repo.On("SearchLocal", mock.Anything, mock.Anything, mock.Anything, 5).Return([]domain.SearchItem{}, nil).Maybe()
	h.uc = h.build(upstash.New(h.srv.URL, h.srv.Token()))
	return h
}

func (h *searchHarness) build(store *upstash.Client) domain.SearchUsecase {
	nowFn := func() time.Time { return h.now }
	guard := ytquota.New(store, h.mailer, nil, ytquota.Options{GuardUnits: 9000, PerUserDaily: 10, Now: nowFn})
	return usecase.NewSearchUsecase(usecase.SearchDeps{
		Repo:         h.repo,
		YouTube:      ytquota.WrapYouTube(h.yt, guard),
		Guard:        guard,
		Cache:        searchcache.New(store),
		Deduper:      searchcache.NewHitDeduper(store),
		CategoryRepo: &stubCategoryRepo{active: map[string][]domain.VideoCategory{"TH": {{ID: "20", CountryCode: "TH", IsActive: true}, {ID: "24", CountryCode: "TH", IsActive: true}}}},
		CategoryConfigs: []domain.CategoryVideoConfig{
			{CategoryID: "20", TableName: "gaming_videos"},
			{CategoryID: "23", TableName: "comedy_videos"},
		},
		Countries:  []string{"TH", "US"},
		VisitorKey: []byte("test-visitor-key"),
		Now:        nowFn,
	})
}

func ytVideos(n int) []domain.SearchVideo {
	out := make([]domain.SearchVideo, n)
	for i := range out {
		out[i] = domain.SearchVideo{ID: "vid" + strings.Repeat("x", 7) + string(rune('a'+i)), Title: "song", CategoryID: "10", ViewCount: "100"}
	}
	return out
}

func TestSearchLocal_SRCH01_ExistingSongsNeverCallYouTube(t *testing.T) {
	h := newSearchHarness(t)
	items := []domain.SearchItem{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}, {ID: "e"}}
	h.repo.On("SearchLocal", "TH", "10", "ฝนตก", 30).Return(items, nil)

	res, err := h.uc.SearchLocal("9.9.9.9", "  ฝนตก ", "th", "")
	require.NoError(t, err)
	assert.Len(t, res.Items, 5)
	assert.False(t, res.SuggestYouTube)
	h.yt.AssertNotCalled(t, "SearchVideos", mock.Anything, mock.Anything, mock.Anything)
	assert.Len(t, h.yt.Calls, 0)
	_, ok := h.srv.Value("yt:quota:2026-11-06")
	assert.False(t, ok)
}

func TestSearchLocal_SuggestsYouTubeWhenFewerThanFive(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchLocal", "TH", "10", "rain", 30).Return([]domain.SearchItem{}, nil)
	res, err := h.uc.SearchLocal("9.9.9.9", "rain", "TH", "10")
	require.NoError(t, err)
	assert.NotNil(t, res.Items)
	assert.True(t, res.SuggestYouTube)
}

func TestSearchYouTube_SRCH02_CallsYouTubeAndWritesSearchVideos(t *testing.T) {
	h := newSearchHarness(t)
	videos := ytVideos(3)
	videos[2].CategoryID = ""
	h.yt.On("SearchVideos", "ฝนตก", "TH", "10").Return(videos, nil).Once()
	h.repo.On("UpsertSearchVideos", mock.MatchedBy(func(v []domain.SearchVideo) bool {
		return len(v) == 3 && v[0].FirstSeenAt.Equal(fixedNow) && v[0].StatsUpdatedAt.Equal(fixedNow) && v[2].CategoryID == "10"
	})).Return(nil).Once()

	res, err := h.uc.SearchYouTube("user-1", "ฝนตก", "TH", "10")
	require.NoError(t, err)
	assert.Len(t, res.Items, 3)
	assert.Equal(t, domain.SearchSourceYouTube, res.Items[0].Source)
	assert.Equal(t, int64(9), res.Remaining)
	assert.False(t, res.Cached)
	h.repo.AssertExpectations(t)

	v, _ := h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "101", v)
	_, cached := h.srv.Value("yt:search:TH:10:ฝนตก")
	assert.True(t, cached)
	assert.Equal(t, int64(86400), h.srv.TTL("yt:search:TH:10:ฝนตก"))
}

func TestSearchYouTube_SRCH03_RepeatWithin24hServedFromCache(t *testing.T) {
	h := newSearchHarness(t)
	h.yt.On("SearchVideos", "Rain  Song", "TH", "10").Return(ytVideos(2), nil).Once()
	h.repo.On("UpsertSearchVideos", mock.Anything).Return(nil).Once()

	_, err := h.uc.SearchYouTube("user-1", "Rain  Song", "TH", "10")
	require.NoError(t, err)

	res, err := h.uc.SearchYouTube("user-2", " rain song ", "TH", "10")
	require.NoError(t, err)
	assert.True(t, res.Cached)
	assert.Len(t, res.Items, 2)
	assert.Equal(t, int64(10), res.Remaining, "cache hit must not consume user credit")
	h.yt.AssertNumberOfCalls(t, "SearchVideos", 1)
	v, _ := h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "101", v, "cache hit must not consume system quota")
}

func TestSearch_SRCH04_InvalidInputReturnsValidationError(t *testing.T) {
	h := newSearchHarness(t)
	cases := []struct{ name, q, country string }{
		{"q too short", "a", "TH"},
		{"q blank after trim", "   x   ", "TH"},
		{"q too long", strings.Repeat("ก", 101), "TH"},
		{"country unsupported", "rain", "FR"},
		{"country empty", "rain", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := h.uc.SearchLocal("9.9.9.9", tc.q, tc.country, "10")
			var ve *domain.ValidationError
			assert.ErrorAs(t, err, &ve)
			_, err = h.uc.SearchYouTube("user-1", tc.q, tc.country, "10")
			assert.ErrorAs(t, err, &ve)
		})
	}
	h.repo.On("SearchLocal", "TH", "10", strings.Repeat("ก", 100), 30).Return([]domain.SearchItem{}, nil)
	_, err := h.uc.SearchLocal("9.9.9.9", strings.Repeat("ก", 100), "TH", "10")
	assert.NoError(t, err, "100 runes is allowed")
	assert.Len(t, h.yt.Calls, 0)
}

func TestSearchYouTube_SRCH07_UserLimitButCacheStillServed(t *testing.T) {
	h := newSearchHarness(t)
	h.srv.Set("yt:user:user-1:2026-11-06", "10")
	h.srv.Set("yt:search:TH:10:cached song", `[{"id":"abc","title":"Cached","source":"youtube"}]`)

	_, err := h.uc.SearchYouTube("user-1", "new song", "TH", "10")
	var qe *domain.QuotaError
	require.ErrorAs(t, err, &qe)
	assert.ErrorIs(t, err, domain.ErrUserLimit)
	assert.True(t, qe.RetryAt.Equal(time.Date(2026, 11, 7, 8, 0, 0, 0, time.UTC)))
	assert.Equal(t, int64(0), qe.Remaining)

	res, err := h.uc.SearchYouTube("user-1", "Cached Song", "TH", "10")
	require.NoError(t, err)
	assert.True(t, res.Cached)
	assert.Equal(t, "abc", res.Items[0].ID)
	assert.Equal(t, int64(0), res.Remaining)
	assert.Len(t, h.yt.Calls, 0)
}

func TestSearch_SRCH08_CategoryRouting(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchLocal", "TH", "20", "minecraft", 30).Return([]domain.SearchItem{}, nil).Once()
	_, err := h.uc.SearchLocal("9.9.9.9", "minecraft", "TH", "20")
	require.NoError(t, err)
	h.repo.AssertExpectations(t)

	var ve *domain.ValidationError
	for _, cat := range []string{"23", "24", "999", "abc"} {
		_, err = h.uc.SearchLocal("9.9.9.9", "minecraft", "TH", cat)
		assert.ErrorAs(t, err, &ve, "category %s", cat)
		_, err = h.uc.SearchYouTube("user-1", "minecraft", "TH", cat)
		assert.ErrorAs(t, err, &ve, "category %s", cat)
	}
	_, err = h.uc.SearchLocal("9.9.9.9", "minecraft", "US", "20")
	assert.ErrorAs(t, err, &ve, "category inactive for country")
}

func TestSearchYouTube_QUOTA01_CounterNearGuardBlocksWithoutCallingYouTube(t *testing.T) {
	h := newSearchHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "8900")

	_, err := h.uc.SearchYouTube("user-1", "rain", "TH", "10")
	var qe *domain.QuotaError
	require.ErrorAs(t, err, &qe)
	assert.ErrorIs(t, err, domain.ErrQuotaExhausted)
	assert.False(t, qe.RetryAt.IsZero())
	assert.Len(t, h.yt.Calls, 0)
	assert.Len(t, h.mailer.sent, 1)
	v, _ := h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "8900", v)
}

func TestSearchYouTube_QuotaErrorsCarryRealRemainingCredits(t *testing.T) {
	t.Run("quota_exhausted", func(t *testing.T) {
		h := newSearchHarness(t)
		h.srv.Set("yt:user:user-1:2026-11-06", "3")
		h.srv.Set("yt:quota:2026-11-06", "8900")

		_, err := h.uc.SearchYouTube("user-1", "rain", "TH", "10")
		var qe *domain.QuotaError
		require.ErrorAs(t, err, &qe)
		assert.ErrorIs(t, err, domain.ErrQuotaExhausted)
		assert.Equal(t, int64(7), qe.Remaining)
		v, _ := h.srv.Value("yt:user:user-1:2026-11-06")
		assert.Equal(t, "3", v, "reading remaining must not consume credit")
	})
	t.Run("search_unavailable on quota counter", func(t *testing.T) {
		h := newSearchHarness(t)
		h.srv.Set("yt:user:user-1:2026-11-06", "3")
		h.srv.FailKeyPrefix("yt:quota:")

		_, err := h.uc.SearchYouTube("user-1", "rain", "TH", "10")
		var qe *domain.QuotaError
		require.ErrorAs(t, err, &qe)
		assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
		assert.Equal(t, int64(7), qe.Remaining)
		v, _ := h.srv.Value("yt:user:user-1:2026-11-06")
		assert.Equal(t, "3", v)
	})
	t.Run("search_unavailable on cache", func(t *testing.T) {
		h := newSearchHarness(t)
		h.srv.Set("yt:user:user-1:2026-11-06", "3")
		h.srv.FailKeyPrefix("yt:search:")

		_, err := h.uc.SearchYouTube("user-1", "rain", "TH", "10")
		var qe *domain.QuotaError
		require.ErrorAs(t, err, &qe)
		assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
		assert.Equal(t, int64(7), qe.Remaining)
	})
}

func TestSearch_QUOTA04_RedisDownFailsClosedForYouTubeOnly(t *testing.T) {
	h := newSearchHarness(t)
	h.srv.SetFail(true)
	h.repo.On("SearchLocal", "TH", "10", "rain", 30).Return([]domain.SearchItem{{ID: "a"}}, nil)

	res, err := h.uc.SearchLocal("9.9.9.9", "rain", "TH", "10")
	require.NoError(t, err)
	assert.Len(t, res.Items, 1)

	_, err = h.uc.SearchYouTube("user-1", "rain", "TH", "10")
	assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)

	uc := h.build(upstash.New("", ""))
	_, err = uc.SearchYouTube("user-1", "rain", "TH", "10")
	assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
	assert.Len(t, h.yt.Calls, 0)
	h.repo.AssertNotCalled(t, "CreateQueryHit", mock.Anything)
}

func TestSearchYouTube_ReleasesUserCreditOnYouTubeError(t *testing.T) {
	h := newSearchHarness(t)
	h.yt.On("SearchVideos", "rain", "TH", "10").Return([]domain.SearchVideo(nil), errors.New("403 quotaExceeded"))

	_, err := h.uc.SearchYouTube("user-1", "rain", "TH", "10")
	assert.Error(t, err)
	v, _ := h.srv.Value("yt:user:user-1:2026-11-06")
	assert.Equal(t, "0", v)
	_, cached := h.srv.Value("yt:search:TH:10:rain")
	assert.False(t, cached)
}

func TestRecordHit_TRND03_SameIPAndVideoCountedOncePerHour(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("VideoExists", "dQw4w9WgXcQ").Return(true, nil)
	h.repo.On("CreateHit", domain.SearchHit{VideoID: "dQw4w9WgXcQ", CountryCode: "TH", QueryNorm: "rain song", HitAt: fixedNow}).Return(nil).Once()

	counted, err := h.uc.RecordHit("1.2.3.4", "dQw4w9WgXcQ", "th", " Rain   Song ")
	require.NoError(t, err)
	assert.True(t, counted)

	counted, err = h.uc.RecordHit("1.2.3.4", "dQw4w9WgXcQ", "TH", "rain song")
	require.NoError(t, err)
	assert.False(t, counted)

	h.repo.AssertNumberOfCalls(t, "CreateHit", 1)
	assert.Equal(t, int64(3600), h.srv.TTL("search:hit:1.2.3.4:dQw4w9WgXcQ"))
}

func TestRecordHit_Validation(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("VideoExists", "AAAAAAAAAAA").Return(false, nil)
	var ve *domain.ValidationError

	_, err := h.uc.RecordHit("ip", "bad id!", "TH", "q")
	assert.ErrorAs(t, err, &ve)
	_, err = h.uc.RecordHit("ip", "dQw4w9WgXcQ", "FR", "q")
	assert.ErrorAs(t, err, &ve)
	_, err = h.uc.RecordHit("ip", "AAAAAAAAAAA", "TH", "q")
	assert.ErrorIs(t, err, domain.ErrUnknownVideo)
}

func TestRecordHit_RedisDownSkipsCounting(t *testing.T) {
	h := newSearchHarness(t)
	h.srv.SetFail(true)
	h.repo.On("VideoExists", "dQw4w9WgXcQ").Return(true, nil)
	counted, err := h.uc.RecordHit("ip", "dQw4w9WgXcQ", "TH", "q")
	require.NoError(t, err)
	assert.False(t, counted)
	h.repo.AssertNotCalled(t, "CreateHit", mock.Anything)
}

func TestGetSearchedTrends_ComputesRankChange(t *testing.T) {
	h := newSearchHarness(t)
	prev := 4
	same := 2
	h.repo.On("SearchedTrends", "TH", fixedNow, 30).Return([]domain.SearchedTrend{
		{Rank: 1, PrevRank: &prev},
		{Rank: 2, PrevRank: &same},
		{Rank: 3},
	}, nil)

	got, err := h.uc.GetSearchedTrends("th")
	require.NoError(t, err)
	assert.Equal(t, 3, got[0].RankChange)
	assert.False(t, got[0].IsNew)
	assert.Equal(t, 0, got[1].RankChange)
	assert.True(t, got[2].IsNew)

	_, err = h.uc.GetSearchedTrends("FR")
	var ve *domain.ValidationError
	assert.ErrorAs(t, err, &ve)
}

func TestRefreshSearchVideoStats_BatchesAndUpserts(t *testing.T) {
	h := newSearchHarness(t)
	ids := []string{"a", "b"}
	h.repo.On("SearchVideoIDsWithHitsSince", fixedNow.Add(-7*24*time.Hour)).Return(ids, nil)
	h.yt.On("FetchVideosByIDs", ids).Return(ytVideos(2), nil)
	h.repo.On("UpsertSearchVideos", mock.MatchedBy(func(v []domain.SearchVideo) bool {
		return len(v) == 2 && v[0].StatsUpdatedAt.Equal(fixedNow)
	})).Return(nil)

	n, err := h.uc.RefreshSearchVideoStats()
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	v, _ := h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "1", v)
}

func TestRefreshSearchVideoStats_NoHitsSkipsYouTube(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchVideoIDsWithHitsSince", mock.Anything).Return([]string{}, nil)
	n, err := h.uc.RefreshSearchVideoStats()
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Len(t, h.yt.Calls, 0)
}

func itemsWithIDs(ids ...string) []domain.SearchItem {
	out := make([]domain.SearchItem, len(ids))
	for i, id := range ids {
		out[i] = domain.SearchItem{ID: id}
	}
	return out
}

func dropExpectations(calls []*mock.Call, method string) []*mock.Call {
	var kept []*mock.Call
	for _, c := range calls {
		if c.Method != method {
			kept = append(kept, c)
		}
	}
	return kept
}

func queryHits(h *searchHarness) []domain.SearchQueryHit {
	var out []domain.SearchQueryHit
	for _, c := range h.repo.Calls {
		if c.Method == "CreateQueryHit" {
			out = append(out, c.Arguments.Get(0).(domain.SearchQueryHit))
		}
	}
	return out
}

func TestSearchLocal_TRND06_SearchCountsQueryAndSavesTopFiveResults(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchLocal", "TH", "10", "Rain  Song", 30).Return(itemsWithIDs("a", "b", "c", "d", "e", "f", "g"), nil)

	_, err := h.uc.SearchLocal("1.2.3.4", "Rain  Song", "th", "")
	require.NoError(t, err)

	hits := queryHits(h)
	require.Len(t, hits, 1)
	assert.Equal(t, "TH", hits[0].CountryCode)
	assert.Equal(t, "rain song", hits[0].QueryNorm)
	assert.Equal(t, fixedNow, hits[0].HitAt)
	mac := hmac.New(sha256.New, []byte("test-visitor-key"))
	mac.Write([]byte("1.2.3.4"))
	wantHash := hex.EncodeToString(mac.Sum(nil))[:32]
	assert.Equal(t, wantHash, hits[0].VisitorHash)
	h.repo.AssertCalled(t, "SaveQueryResults", "TH", "rain song", itemsWithIDs("a", "b", "c", "d", "e"), fixedNow)
	assert.Equal(t, int64(3600), h.srv.TTL("search:q:TH:"+wantHash+":rain song"))
	_, rawIPKey := h.srv.Value("search:q:TH:1.2.3.4:rain song")
	assert.False(t, rawIPKey, "raw IP must not be stored next to the search term")
}

func TestSearchLocal_TRND13_UntrackedSearchFromGroupLinkNotCounted(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchLocal", "TH", "10", "rain song", 30).Return(itemsWithIDs("a"), nil)

	res, err := h.uc.SearchLocal("", "rain song", "TH", "")
	require.NoError(t, err)
	assert.Len(t, res.Items, 1)
	h.repo.AssertNotCalled(t, "CreateQueryHit", mock.Anything)
	h.repo.AssertNotCalled(t, "SaveQueryResults", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSearchLocal_TRND07_SameIPAndQueryCountedOncePerHour(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchLocal", "TH", "10", mock.Anything, 30).Return(itemsWithIDs("a"), nil)

	for _, q := range []string{"rain song", " RAIN   song "} {
		_, err := h.uc.SearchLocal("1.2.3.4", q, "TH", "")
		require.NoError(t, err)
	}
	require.Len(t, queryHits(h), 1)

	_, err := h.uc.SearchLocal("5.6.7.8", "rain song", "TH", "")
	require.NoError(t, err)
	hits := queryHits(h)
	require.Len(t, hits, 2)
	assert.NotEqual(t, hits[0].VisitorHash, hits[1].VisitorHash)
}

func TestSearchLocal_TRND08_NonMusicOrEmptyResultsNotCounted(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("SearchLocal", "TH", "20", "rain song", 30).Return(itemsWithIDs("g"), nil)
	h.repo.On("SearchLocal", "TH", "10", "nothing here", 30).Return([]domain.SearchItem{}, nil)

	_, err := h.uc.SearchLocal("1.2.3.4", "rain song", "TH", "20")
	require.NoError(t, err)
	_, err = h.uc.SearchLocal("1.2.3.4", "nothing here", "TH", "")
	require.NoError(t, err)

	h.repo.AssertNotCalled(t, "CreateQueryHit", mock.Anything)
	h.repo.AssertNotCalled(t, "SaveQueryResults", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSearchYouTube_TRND09_FreshResultsRefreshGroupWithoutCounting(t *testing.T) {
	h := newSearchHarness(t)
	h.yt.On("SearchVideos", "Rain Song", "TH", "10").Return(ytVideos(2), nil).Once()
	h.repo.On("UpsertSearchVideos", mock.Anything).Return(nil).Once()
	h.repo.ExpectedCalls = dropExpectations(h.repo.ExpectedCalls, "SearchLocal")
	h.repo.On("SearchLocal", "TH", "10", "Rain Song", 5).Return(itemsWithIDs("a", "yt1"), nil).Once()

	_, err := h.uc.SearchYouTube("user-1", "Rain Song", "TH", "10")
	require.NoError(t, err)
	res, err := h.uc.SearchYouTube("user-2", "rain song", "TH", "10")
	require.NoError(t, err)
	require.True(t, res.Cached)

	h.repo.AssertNumberOfCalls(t, "SaveQueryResults", 1)
	h.repo.AssertCalled(t, "SaveQueryResults", "TH", "rain song", itemsWithIDs("a", "yt1"), fixedNow)
	assert.Empty(t, queryHits(h))
}

func TestSearchYouTube_NonMusicDoesNotTouchGroups(t *testing.T) {
	h := newSearchHarness(t)
	h.yt.On("SearchVideos", "speedrun", "TH", "20").Return(ytVideos(1), nil).Once()
	h.repo.On("UpsertSearchVideos", mock.Anything).Return(nil).Once()

	_, err := h.uc.SearchYouTube("user-1", "speedrun", "TH", "20")
	require.NoError(t, err)
	h.repo.AssertNotCalled(t, "SaveQueryResults", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestSearch_TRND10_GroupWriteFailureDoesNotFailSearch(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.ExpectedCalls = nil
	h.repo.On("CreateQueryHit", mock.Anything).Return(errors.New("db down"))
	h.repo.On("SaveQueryResults", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db down"))
	h.repo.On("SearchLocal", "TH", "10", "rain song", 30).Return(itemsWithIDs("a"), nil)

	res, err := h.uc.SearchLocal("1.2.3.4", "rain song", "TH", "")
	require.NoError(t, err)
	assert.Len(t, res.Items, 1)
}

func TestGetSearchedQueries_TRND11_LastSevenDaysAtLeastTwoVisitors(t *testing.T) {
	h := newSearchHarness(t)
	want := []domain.SearchedQuery{{Query: "rain song", Rank: 1, Searches: 4, Items: itemsWithIDs("a")}}
	h.repo.On("SearchedQueries", "TH", fixedNow.Add(-7*24*time.Hour), 2, 30).Return(want, nil)

	got, err := h.uc.GetSearchedQueries("th")
	require.NoError(t, err)
	assert.Equal(t, want, got)

	var ve *domain.ValidationError
	_, err = h.uc.GetSearchedQueries("FR")
	assert.ErrorAs(t, err, &ve)
}

func TestPruneSearchQueryHits_TRND12_DeletesHitsOlderThanSevenDays(t *testing.T) {
	h := newSearchHarness(t)
	h.repo.On("DeleteQueryHitsBefore", fixedNow.Add(-7*24*time.Hour)).Return(int64(12), nil)

	n, err := h.uc.PruneSearchQueryHits()
	require.NoError(t, err)
	assert.Equal(t, int64(12), n)
}
