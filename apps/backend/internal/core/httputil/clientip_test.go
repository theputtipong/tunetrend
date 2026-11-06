package httputil

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

const testInternalKey = "internal-key-0123456789abcdef0123"

func resolveIP(t *testing.T, key string, headers map[string]string) string {
	t.Helper()
	app := fiber.New()
	resolver := NewClientIPResolver(key)
	app.Get("/", func(c *fiber.Ctx) error { return c.SendString(resolver(c)) })
	req := httptest.NewRequest("GET", "/", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func TestClientIP_SpoofedLeftmostIgnoredWithoutInternalKey(t *testing.T) {
	got := resolveIP(t, testInternalKey, map[string]string{"X-Forwarded-For": "6.6.6.6, 203.0.113.9"})
	if got != "203.0.113.9" {
		t.Fatalf("want rightmost hop 203.0.113.9, got %q", got)
	}
}

func TestClientIP_WrongInternalKeyUsesRightmost(t *testing.T) {
	got := resolveIP(t, testInternalKey, map[string]string{
		"X-Forwarded-For": "6.6.6.6, 203.0.113.9",
		InternalKeyHeader: "wrong",
	})
	if got != "203.0.113.9" {
		t.Fatalf("want 203.0.113.9, got %q", got)
	}
}

func TestClientIP_TrustedFrontendUsesForwardedClient(t *testing.T) {
	got := resolveIP(t, testInternalKey, map[string]string{
		"X-Forwarded-For": "198.51.100.7, 76.76.21.1",
		InternalKeyHeader: testInternalKey,
	})
	if got != "198.51.100.7" {
		t.Fatalf("want 198.51.100.7, got %q", got)
	}
}

func TestClientIP_EmptyConfiguredKeyNeverTrusts(t *testing.T) {
	got := resolveIP(t, "", map[string]string{
		"X-Forwarded-For": "6.6.6.6, 203.0.113.9",
		InternalKeyHeader: "",
	})
	if got != "203.0.113.9" {
		t.Fatalf("want 203.0.113.9, got %q", got)
	}
}

func TestClientIP_NoForwardedHeaderFallsBackToPeer(t *testing.T) {
	got := resolveIP(t, testInternalKey, nil)
	if got != "0.0.0.0" {
		t.Fatalf("want peer ip 0.0.0.0, got %q", got)
	}
}
