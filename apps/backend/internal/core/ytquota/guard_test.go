package ytquota_test

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"tunetrend-backend/internal/core/upstash"
	"tunetrend-backend/internal/core/upstash/upstashtest"
	"tunetrend-backend/internal/core/ytquota"
	"tunetrend-backend/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeMailer struct {
	mu   sync.Mutex
	sent []domain.MailMessage
	fail error
}

func (m *fakeMailer) Send(msg domain.MailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fail != nil {
		return m.fail
	}
	m.sent = append(m.sent, msg)
	return nil
}

type fakeLogRepo struct {
	domain.WorkerLogRepository
	mu   sync.Mutex
	logs []domain.WorkerLog
}

func (r *fakeLogRepo) CreateLog(l domain.WorkerLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, l)
	return nil
}

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func pt(t *testing.T, value string) time.Time {
	loc, err := time.LoadLocation("America/Los_Angeles")
	require.NoError(t, err)
	ts, err := time.ParseInLocation("2006-01-02 15:04", value, loc)
	require.NoError(t, err)
	return ts
}

type harness struct {
	srv    *upstashtest.Server
	guard  *ytquota.Guard
	mailer *fakeMailer
	logs   *fakeLogRepo
	clk    *clock
}

func newHarness(t *testing.T) *harness {
	srv := upstashtest.NewServer(t)
	h := &harness{srv: srv, mailer: &fakeMailer{}, logs: &fakeLogRepo{}, clk: &clock{t: pt(t, "2026-11-06 10:00")}}
	h.guard = ytquota.New(upstash.New(srv.URL, srv.Token()), h.mailer, h.logs, ytquota.Options{
		GuardUnits:   9000,
		PerUserDaily: 10,
		Now:          h.clk.now,
	})
	return h
}

func TestGuard_DateKeyUsesPacificTime(t *testing.T) {
	h := newHarness(t)
	h.clk.t = time.Date(2026, 11, 7, 7, 59, 0, 0, time.UTC)
	h.guard.Record(1, ytquota.SourceWorker)
	v, ok := h.srv.Value("yt:quota:2026-11-06")
	require.True(t, ok)
	assert.Equal(t, "1", v)
	assert.Equal(t, int64(48*3600), h.srv.TTL("yt:quota:2026-11-06"))

	h.clk.t = time.Date(2026, 11, 7, 8, 0, 0, 0, time.UTC)
	h.guard.Record(1, ytquota.SourceWorker)
	v, ok = h.srv.Value("yt:quota:2026-11-07")
	require.True(t, ok)
	assert.Equal(t, "1", v)
}

func TestGuard_RetryAtIsNextPacificMidnight(t *testing.T) {
	h := newHarness(t)
	assert.True(t, h.guard.RetryAt().Equal(pt(t, "2026-11-07 00:00")))
}

func TestGuard_QUOTA01_ReserveRejectsWhenCounterPlusSearchCostExceedsGuard(t *testing.T) {
	h := newHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "8900")

	_, err := h.guard.Reserve("user-1")
	assert.ErrorIs(t, err, domain.ErrQuotaExhausted)

	used, _ := h.srv.Value("yt:user:user-1:2026-11-06")
	assert.Equal(t, "0", used, "user credit must be refunded when system quota blocks")

	v, _ := h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "8900", v, "rejected reservation must be rolled back")

	h.srv.Set("yt:quota:2026-11-06", "8899")
	remaining, err := h.guard.Reserve("user-1")
	require.NoError(t, err)
	assert.Equal(t, int64(9), remaining)
	v, _ = h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "9000", v, "successful reservation must count system units up front")
	s, _ := h.srv.Value("yt:quota:2026-11-06:search")
	assert.Equal(t, "101", s)
	assert.Equal(t, int64(48*3600), h.srv.TTL("yt:quota:2026-11-06:search"))
}

func TestGuard_QUOTA01_ConcurrentReservationsNeverOvershootGuard(t *testing.T) {
	h := newHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "8800")

	var wg sync.WaitGroup
	var mu sync.Mutex
	succeeded := 0
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := h.guard.Reserve("user-" + strconv.Itoa(i))
			if err == nil {
				mu.Lock()
				succeeded++
				mu.Unlock()
				return
			}
			assert.ErrorIs(t, err, domain.ErrQuotaExhausted)
		}(i)
	}
	wg.Wait()

	assert.LessOrEqual(t, succeeded, 1)
	v, _ := h.srv.Value("yt:quota:2026-11-06")
	total, err := strconv.ParseInt(v, 10, 64)
	require.NoError(t, err)
	assert.LessOrEqual(t, total, int64(9000))
	assert.Equal(t, int64(8800+101*succeeded), total)
}

func TestGuard_QUOTA04_SystemCounterWriteFailureFailsClosedAndRefundsUser(t *testing.T) {
	h := newHarness(t)
	h.srv.FailKeyPrefix("yt:quota:")

	_, err := h.guard.Reserve("user-1")
	assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
	used, _ := h.srv.Value("yt:user:user-1:2026-11-06")
	assert.Equal(t, "0", used)
}

func TestGuard_QUOTA03_AlertSentOncePerPacificDayAndCounterResets(t *testing.T) {
	h := newHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "8950")
	h.srv.Set("yt:quota:2026-11-06:worker", "550")
	h.srv.Set("yt:quota:2026-11-06:search", "8400")

	for i := 0; i < 3; i++ {
		_, err := h.guard.Reserve("user-1")
		assert.ErrorIs(t, err, domain.ErrQuotaExhausted)
	}

	require.Len(t, h.mailer.sent, 1)
	assert.Contains(t, h.mailer.sent[0].Subject, "9,000")
	assert.Contains(t, h.mailer.sent[0].Body, "worker 550")
	assert.Contains(t, h.mailer.sent[0].Body, "search 8,400")
	require.Len(t, h.logs.logs, 1)
	assert.Equal(t, domain.YouTubeQuotaAlertJobName, h.logs.logs[0].JobName)
	assert.Equal(t, int64(172800), h.srv.TTL("yt:quota:alerted:2026-11-06"))

	h.clk.t = pt(t, "2026-11-07 00:01")
	remaining, err := h.guard.Reserve("user-1")
	require.NoError(t, err)
	assert.Equal(t, int64(9), remaining)
	assert.Len(t, h.mailer.sent, 1)
}

func TestGuard_QUOTA03_WorkerRecordReachingGuardTriggersAlert(t *testing.T) {
	h := newHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "8998")

	h.guard.Record(1, ytquota.SourceWorker)
	assert.Len(t, h.mailer.sent, 0)

	h.guard.Record(1, ytquota.SourceWorker)
	require.Len(t, h.mailer.sent, 1)
	assert.Contains(t, h.mailer.sent[0].Body, "9,000")
	require.Len(t, h.logs.logs, 1)

	h.guard.Record(1, ytquota.SourceWorker)
	assert.Len(t, h.mailer.sent, 1)
}

func TestGuard_QUOTA03_FailedAlertMailIsRetried(t *testing.T) {
	h := newHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "8950")
	h.mailer.fail = errBoom

	_, err := h.guard.Reserve("user-1")
	assert.ErrorIs(t, err, domain.ErrQuotaExhausted)
	_, alerted := h.srv.Value("yt:quota:alerted:2026-11-06")
	assert.False(t, alerted, "alerted flag must be cleared when mail fails")

	h.mailer.fail = nil
	_, err = h.guard.Reserve("user-1")
	assert.ErrorIs(t, err, domain.ErrQuotaExhausted)
	assert.Len(t, h.mailer.sent, 1)
	assert.Len(t, h.logs.logs, 1)
	_, alerted = h.srv.Value("yt:quota:alerted:2026-11-06")
	assert.True(t, alerted)
}

func TestGuard_QUOTA04_FailClosedWhenRedisUnavailable(t *testing.T) {
	t.Run("unconfigured", func(t *testing.T) {
		g := ytquota.New(upstash.New("", ""), nil, nil, ytquota.Options{GuardUnits: 9000, PerUserDaily: 10})
		_, err := g.Reserve("user-1")
		assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
		_, err = g.Remaining("user-1")
		assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
	})
	t.Run("server error", func(t *testing.T) {
		h := newHarness(t)
		h.srv.SetFail(true)
		_, err := h.guard.Reserve("user-1")
		assert.ErrorIs(t, err, domain.ErrQuotaUnavailable)
	})
}

func TestGuard_UserLimit(t *testing.T) {
	h := newHarness(t)
	for i := 1; i <= 10; i++ {
		remaining, err := h.guard.Reserve("user-1")
		require.NoError(t, err)
		assert.Equal(t, int64(10-i), remaining)
	}
	_, err := h.guard.Reserve("user-1")
	assert.ErrorIs(t, err, domain.ErrUserLimit)

	remaining, err := h.guard.Remaining("user-1")
	require.NoError(t, err)
	assert.Equal(t, int64(0), remaining)

	_, err = h.guard.Reserve("user-2")
	assert.NoError(t, err)
}

func TestGuard_ReleaseRefundsUserCredit(t *testing.T) {
	h := newHarness(t)
	_, err := h.guard.Reserve("user-1")
	require.NoError(t, err)
	h.guard.Release("user-1")
	remaining, err := h.guard.Remaining("user-1")
	require.NoError(t, err)
	assert.Equal(t, int64(10), remaining)
}

func TestOptionsFromEnv(t *testing.T) {
	t.Setenv("YOUTUBE_QUOTA_GUARD_UNITS", "")
	t.Setenv("YOUTUBE_SEARCH_PER_USER_DAILY", "")
	o := ytquota.OptionsFromEnv()
	assert.Equal(t, int64(9000), o.GuardUnits)
	assert.Equal(t, int64(10), o.PerUserDaily)

	t.Setenv("YOUTUBE_QUOTA_GUARD_UNITS", "5000")
	t.Setenv("YOUTUBE_SEARCH_PER_USER_DAILY", "abc")
	o = ytquota.OptionsFromEnv()
	assert.Equal(t, int64(5000), o.GuardUnits)
	assert.Equal(t, int64(10), o.PerUserDaily)
}

var errBoom = errors.New("boom")

type countingYT struct {
	calls []string
}

func (c *countingYT) FetchTrending(string) ([]domain.Song, error) {
	c.calls = append(c.calls, "FetchTrending")
	return []domain.Song{{ID: "a"}}, nil
}
func (c *countingYT) FetchVideoCategories(string) ([]domain.VideoCategory, error) {
	c.calls = append(c.calls, "FetchVideoCategories")
	return nil, nil
}
func (c *countingYT) FetchTrendingByCategory(string, string) ([]domain.CategoryVideo, error) {
	c.calls = append(c.calls, "FetchTrendingByCategory")
	return nil, errBoom
}
func (c *countingYT) SearchVideos(string, string, string) ([]domain.SearchVideo, error) {
	c.calls = append(c.calls, "SearchVideos")
	return nil, nil
}
func (c *countingYT) FetchVideosByIDs(ids []string) ([]domain.SearchVideo, error) {
	c.calls = append(c.calls, "FetchVideosByIDs")
	return nil, nil
}

func TestWrapYouTube_QUOTA02_WorkerCallsCountedButNeverBlocked(t *testing.T) {
	h := newHarness(t)
	h.srv.Set("yt:quota:2026-11-06", "9500")
	inner := &countingYT{}
	yt := ytquota.WrapYouTube(inner, h.guard)

	songs, err := yt.FetchTrending("TH")
	require.NoError(t, err)
	assert.Len(t, songs, 1)
	v, _ := h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "9501", v)

	_, _ = yt.FetchVideoCategories("TH")
	_, err = yt.FetchTrendingByCategory("TH", "20")
	assert.ErrorIs(t, err, errBoom)
	v, _ = h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "9503", v, "failed requests still consume quota")

	ids := make([]string, 101)
	for i := range ids {
		ids[i] = strings.Repeat("x", 11)
	}
	_, _ = yt.FetchVideosByIDs(ids)
	v, _ = h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "9506", v)
	w, _ := h.srv.Value("yt:quota:2026-11-06:worker")
	assert.Equal(t, "6", w)

	_, _ = yt.SearchVideos("q", "TH", "10")
	v, _ = h.srv.Value("yt:quota:2026-11-06")
	assert.Equal(t, "9506", v, "search units are counted by Reserve, not the wrapper")
	_, hasSearch := h.srv.Value("yt:quota:2026-11-06:search")
	assert.False(t, hasSearch)

	assert.Equal(t, []string{"FetchTrending", "FetchVideoCategories", "FetchTrendingByCategory", "FetchVideosByIDs", "SearchVideos"}, inner.calls)

	h.srv.SetFail(true)
	_, err = yt.FetchTrending("TH")
	assert.NoError(t, err, "redis outage must not block worker calls")
}
