package ytquota

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	_ "time/tzdata"

	"tunetrend-backend/internal/core/upstash"
	"tunetrend-backend/internal/domain"
)

const (
	SourceWorker = "worker"
	SourceSearch = "search"

	SearchListUnits = 100
	VideosListUnits = 1
	SearchCallUnits = SearchListUnits + VideosListUnits
	VideosBatchSize = 50

	DefaultGuardUnits   = 9000
	DefaultPerUserDaily = 10

	counterTTL = 48 * time.Hour
)

var pacific = mustLoadPacific()

func mustLoadPacific() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	return loc
}

type Options struct {
	GuardUnits   int64
	PerUserDaily int64
	Now          func() time.Time
}

func OptionsFromEnv() Options {
	return Options{
		GuardUnits:   envInt("YOUTUBE_QUOTA_GUARD_UNITS", DefaultGuardUnits),
		PerUserDaily: envInt("YOUTUBE_SEARCH_PER_USER_DAILY", DefaultPerUserDaily),
	}
}

func envInt(key string, def int64) int64 {
	v, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

type Guard struct {
	store   *upstash.Client
	mailer  domain.MailSender
	logRepo domain.WorkerLogRepository
	opts    Options
}

func New(store *upstash.Client, mailer domain.MailSender, logRepo domain.WorkerLogRepository, opts Options) *Guard {
	if opts.GuardUnits <= 0 {
		opts.GuardUnits = DefaultGuardUnits
	}
	if opts.PerUserDaily <= 0 {
		opts.PerUserDaily = DefaultPerUserDaily
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	return &Guard{store: store, mailer: mailer, logRepo: logRepo, opts: opts}
}

func (g *Guard) day() string {
	return g.opts.Now().In(pacific).Format("2006-01-02")
}

func (g *Guard) quotaKey(day string) string { return "yt:quota:" + day }

func (g *Guard) userKey(userID, day string) string { return "yt:user:" + userID + ":" + day }

func (g *Guard) RetryAt() time.Time {
	now := g.opts.Now().In(pacific)
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, pacific)
}

func (g *Guard) Record(units int64, source string) {
	if units <= 0 {
		return
	}
	day := g.day()
	total, err := g.incr(day, units, source)
	if err != nil {
		log.Printf("⚠️ [YTQuota] บันทึก %d units (%s) ไม่สำเร็จ: %v", units, source, err)
		return
	}
	if total >= g.opts.GuardUnits {
		g.alert(day, total)
	}
}

func (g *Guard) incr(day string, units int64, source string) (int64, error) {
	key := g.quotaKey(day)
	sourceKey := key + ":" + source
	n := strconv.FormatInt(units, 10)
	ttl := strconv.FormatInt(int64(counterTTL.Seconds()), 10)
	res, err := g.store.Pipeline(
		[]string{"INCRBY", key, n},
		[]string{"EXPIRE", key, ttl, "NX"},
		[]string{"INCRBY", sourceKey, n},
		[]string{"EXPIRE", sourceKey, ttl, "NX"},
	)
	if err != nil {
		return 0, err
	}
	return res[0].Int()
}

func (g *Guard) decr(day string, units int64, source string) {
	key := g.quotaKey(day)
	n := strconv.FormatInt(units, 10)
	if _, err := g.store.Pipeline(
		[]string{"DECRBY", key, n},
		[]string{"DECRBY", key + ":" + source, n},
	); err != nil {
		log.Printf("⚠️ [YTQuota] คืน %d units (%s) ไม่สำเร็จ: %v", units, source, err)
	}
}

func (g *Guard) Reserve(userID string) (int64, error) {
	day := g.day()
	userKey := g.userKey(userID, day)

	used, err := g.store.IncrByWithTTL(userKey, 1, counterTTL)
	if err != nil {
		log.Printf("⚠️ [YTQuota] ตรวจสิทธิ์ผู้ใช้ไม่สำเร็จ (fail-closed): %v", err)
		return 0, domain.ErrQuotaUnavailable
	}
	if used > g.opts.PerUserDaily {
		g.refund(userKey)
		return 0, domain.ErrUserLimit
	}

	total, err := g.incr(day, SearchCallUnits, SourceSearch)
	if err != nil {
		g.refund(userKey)
		log.Printf("⚠️ [YTQuota] จองโควตาระบบไม่สำเร็จ (fail-closed): %v", err)
		return 0, domain.ErrQuotaUnavailable
	}
	if total > g.opts.GuardUnits {
		g.decr(day, SearchCallUnits, SourceSearch)
		g.refund(userKey)
		g.alert(day, total-SearchCallUnits)
		return 0, domain.ErrQuotaExhausted
	}

	return g.opts.PerUserDaily - used, nil
}

func (g *Guard) Release(userID string) {
	g.refund(g.userKey(userID, g.day()))
}

func (g *Guard) refund(userKey string) {
	if err := g.store.DecrBy(userKey, 1); err != nil {
		log.Printf("⚠️ [YTQuota] คืนสิทธิ์ผู้ใช้ไม่สำเร็จ: %v", err)
	}
}

func (g *Guard) Remaining(userID string) (int64, error) {
	used, err := g.store.GetInt(g.userKey(userID, g.day()))
	if err != nil {
		return 0, domain.ErrQuotaUnavailable
	}
	remaining := g.opts.PerUserDaily - used
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

func (g *Guard) alert(day string, total int64) {
	alertedKey := "yt:quota:alerted:" + day
	first, err := g.store.SetNX(alertedKey, "1", counterTTL)
	if err != nil || !first {
		return
	}

	var workerUnits, searchUnits int64
	if res, err := g.store.Pipeline(
		[]string{"GET", g.quotaKey(day) + ":" + SourceWorker},
		[]string{"GET", g.quotaKey(day) + ":" + SourceSearch},
	); err == nil {
		workerUnits, _ = res[0].Int()
		searchUnits, _ = res[1].Int()
	}

	limit := formatThousands(g.opts.GuardUnits)
	message := fmt.Sprintf(
		"ค้นบน YouTube ถูกปิดอัตโนมัติจนถึง reset (00:00 PT) วันที่ %s\nตัวนับรวม: %s units\nใช้ไปแล้ว: worker %s · search %s",
		day, formatThousands(total), formatThousands(workerUnits), formatThousands(searchUnits),
	)

	if g.mailer != nil {
		if err := g.mailer.Send(domain.MailMessage{
			Subject: "[TuneTrend] YouTube quota ถึง " + limit + " units",
			Body:    message,
		}); err != nil {
			log.Printf("⚠️ [YTQuota] ส่งอีเมล alert ไม่สำเร็จ: %v", err)
			if err := g.store.Del(alertedKey); err != nil {
				log.Printf("⚠️ [YTQuota] ล้าง flag alert ไม่สำเร็จ: %v", err)
			}
			return
		}
	}

	if g.logRepo != nil {
		now := g.opts.Now()
		if err := g.logRepo.CreateLog(domain.WorkerLog{
			JobName:    domain.YouTubeQuotaAlertJobName,
			Status:     domain.WorkerLogStatusFailed,
			Message:    message,
			StartedAt:  now,
			FinishedAt: now,
		}); err != nil {
			log.Printf("⚠️ [YTQuota] บันทึก worker log ของ alert ไม่สำเร็จ: %v", err)
		}
	}
}

func formatThousands(n int64) string {
	s := strconv.FormatInt(n, 10)
	neg := ""
	if n < 0 {
		neg, s = "-", s[1:]
	}
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return neg + s
}
