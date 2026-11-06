package httputil

import (
	"crypto/subtle"
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

const InternalKeyHeader = "X-Internal-Key"

var (
	defaultResolver     func(*fiber.Ctx) string
	defaultResolverOnce sync.Once
)

func ClientIP(c *fiber.Ctx) string {
	defaultResolverOnce.Do(func() {
		defaultResolver = NewClientIPResolver(os.Getenv("INTERNAL_PROXY_KEY"))
	})
	return defaultResolver(c)
}

// X-Forwarded-For is client-controlled except for the last hop, which the hosting
// proxy appends. Only our own frontend (proven by the shared key) may vouch for the
// leftmost entry, because it forwards the visitor's IP it received from Vercel.
func NewClientIPResolver(internalKey string) func(*fiber.Ctx) string {
	key := []byte(internalKey)
	return func(c *fiber.Ctx) string {
		fwd := c.Get("X-Forwarded-For")
		if fwd == "" {
			return c.IP()
		}
		hops := strings.Split(fwd, ",")
		if len(key) > 0 && subtle.ConstantTimeCompare([]byte(c.Get(InternalKeyHeader)), key) == 1 {
			return strings.TrimSpace(hops[0])
		}
		return strings.TrimSpace(hops[len(hops)-1])
	}
}
