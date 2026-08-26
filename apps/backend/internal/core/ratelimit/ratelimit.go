package ratelimit

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"tunetrend-backend/internal/core/httputil"

	"github.com/gofiber/fiber/v2"
)

type Limiter struct {
	baseURL    string
	token      string
	client     *http.Client
	limit      int64
	window     time.Duration
	failClosed bool
}

func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		baseURL: strings.TrimRight(os.Getenv("UPSTASH_REDIS_REST_URL"), "/"),
		token:   os.Getenv("UPSTASH_REDIS_REST_TOKEN"),
		client:  &http.Client{Timeout: 3 * time.Second},
		limit:   int64(limit),
		window:  window,
	}
}

func (l *Limiter) FailClosed() *Limiter {
	l.failClosed = true
	return l
}

func (l *Limiter) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if l.baseURL == "" || l.token == "" {
			log.Printf("⚠️  [RateLimit] UPSTASH_REDIS_REST_URL/TOKEN not configured — failClosed=%v", l.failClosed)
			return l.onCheckFailed(c)
		}

		key := "tunetrend-backend-ratelimit:" + httputil.ClientIP(c)

		count, err := l.incrementAndExpire(key)
		if err != nil {
			log.Printf("⚠️  [RateLimit] Upstash check failed: %v — failClosed=%v", err, l.failClosed)
			return l.onCheckFailed(c)
		}

		if count > l.limit {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Too many requests. Please slow down and try again shortly.",
			})
		}

		return c.Next()
	}
}

func (l *Limiter) onCheckFailed(c *fiber.Ctx) error {
	if l.failClosed {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"error":   "Rate limiting is temporarily unavailable. Please try again shortly.",
		})
	}
	return c.Next()
}

type pipelineStep struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error,omitempty"`
}

func (l *Limiter) incrementAndExpire(key string) (int64, error) {
	body, err := json.Marshal([][]string{
		{"INCR", key},
		{"EXPIRE", key, strconv.Itoa(int(l.window.Seconds())), "NX"},
	})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, l.baseURL+"/pipeline", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+l.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := l.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var steps []pipelineStep
	if err := json.NewDecoder(res.Body).Decode(&steps); err != nil {
		return 0, err
	}
	if len(steps) == 0 {
		return 0, errors.New("ratelimit: empty pipeline response")
	}
	if steps[0].Error != "" {
		return 0, errors.New(steps[0].Error)
	}

	var count int64
	if err := json.Unmarshal(steps[0].Result, &count); err != nil {
		return 0, err
	}
	return count, nil
}
