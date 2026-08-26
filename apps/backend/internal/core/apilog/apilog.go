package apilog

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"tunetrend-backend/internal/core/httputil"
	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type responseEnvelope struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// New คืนค่า middleware สำหรับบันทึก log ของทุก request
func New(repo domain.ApiLogRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if shouldSkip(c.Path()) {
			return c.Next()
		}

		start := time.Now()
		handlerErr := c.Next()

		entry := &domain.ApiLog{
			RequestID:  requestIDFrom(c),
			Method:     c.Method(),
			Path:       c.Path(),
			Query:      string(c.Request().URI().QueryString()),
			StatusCode: c.Response().StatusCode(),
			ClientIP:   httputil.ClientIP(c),
			UserAgent:  c.Get(fiber.HeaderUserAgent),
			DurationMs: time.Since(start).Milliseconds(),
			CreatedAt:  time.Now(),
		}

		var body responseEnvelope
		if json.Unmarshal(c.Response().Body(), &body) == nil {
			entry.Success = body.Success
			entry.ErrorMessage = body.Error
		} else {
			entry.Success = entry.StatusCode < 400
		}
		if handlerErr != nil {
			entry.Success = false
			if entry.ErrorMessage == "" {
				entry.ErrorMessage = handlerErr.Error()
			}
		}

		printConsoleLog(entry)

		go func() {
			if err := repo.Create(entry); err != nil {
				log.Printf("⚠️  [ApiLog] เขียน log ล้มเหลว (request_id=%s): %v", entry.RequestID, err)
			}
		}()

		return handlerErr
	}
}

func printConsoleLog(entry *domain.ApiLog) {
	icon := "✅"
	if !entry.Success {
		icon = "❌"
	}

	line := fmt.Sprintf("%s [API] %s %s %d %dms request_id=%s ip=%s",
		icon, entry.Method, entry.Path, entry.StatusCode, entry.DurationMs, entry.RequestID, entry.ClientIP)

	if entry.ErrorMessage != "" {
		line += fmt.Sprintf(" error=%q", entry.ErrorMessage)
	}

	log.Println(line)
}

func requestIDFrom(c *fiber.Ctx) string {
	rid, _ := c.Locals(requestid.ConfigDefault.ContextKey).(string)
	return rid
}

func shouldSkip(path string) bool {
	return path == "/health" || strings.HasPrefix(path, "/docs")
}
