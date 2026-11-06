package ratelimit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamedLimiterUsesSeparateKey(t *testing.T) {
	var keys []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cmds [][]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&cmds))
		keys = append(keys, cmds[0][1])
		_, _ = w.Write([]byte(`[{"result":1},{"result":1}]`))
	}))
	defer srv.Close()
	t.Setenv("UPSTASH_REDIS_REST_URL", srv.URL)
	t.Setenv("UPSTASH_REDIS_REST_TOKEN", "t")

	app := fiber.New()
	ok := func(c *fiber.Ctx) error { return c.SendStatus(200) }
	app.Get("/default", New(10, time.Minute).Middleware(), ok)
	app.Get("/named", New(10, time.Minute).Named("auth").Middleware(), ok)

	for _, path := range []string{"/default", "/named"} {
		req := httptest.NewRequest("GET", path, nil)
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		res, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	}
	assert.Equal(t, []string{"tunetrend-backend-ratelimit:1.2.3.4", "tunetrend-backend-ratelimit:auth:1.2.3.4"}, keys)
}
