package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"tunetrend-backend/internal/core/auth"
	deliveryhttp "tunetrend-backend/internal/delivery/http"
	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/repository/postgres"
	"tunetrend-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testJWTSecret = "test-secret-0123456789abcdef-0123456789"

type fakeVerifier map[string]domain.GoogleIdentity

func (f fakeVerifier) Verify(_ context.Context, token string) (*domain.GoogleIdentity, error) {
	id, ok := f[token]
	if !ok {
		return nil, domain.ErrUnauthorized
	}
	return &id, nil
}

var verifier = fakeVerifier{
	"google-a":         {Subject: "sub-a", Email: "a@example.com", Name: "Alice", Picture: "https://lh3.googleusercontent.com/a"},
	"google-a-renamed": {Subject: "sub-a", Email: "a@example.com", Name: "Alice Renamed"},
	"google-b":         {Subject: "sub-b", Email: "b@example.com", Name: "Bob"},
}

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{DSN: dsn, PreferSimpleProtocol: true}),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.User{}, &domain.Playlist{}, &domain.PlaylistItem{}))
	require.NoError(t, db.Exec("TRUNCATE users, playlists, playlist_items").Error)
	return db
}

func newApp(db *gorm.DB, secret string) *fiber.App {
	tokens := auth.NewSessionTokens(secret, auth.SessionTTL)
	authUC := usecase.NewAuthUsecase(verifier, tokens, postgres.NewUserRepository(db))
	playlistUC := usecase.NewPlaylistUsecase(postgres.NewPlaylistRepository(db))
	app := fiber.New()
	deliveryhttp.RegisterAuthRoutes(app, authUC, playlistUC, func(c *fiber.Ctx) error { return c.Next() })
	return app
}

type apiResp struct {
	Status  int
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
}

func call(t *testing.T, app *fiber.App, method, path, token string, body any) apiResp {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := app.Test(req, 10000)
	require.NoError(t, err)
	defer res.Body.Close()
	out := apiResp{Status: res.StatusCode}
	raw, _ := io.ReadAll(res.Body)
	require.NoError(t, json.Unmarshal(raw, &out), string(raw))
	return out
}

func decode[T any](t *testing.T, r apiResp) T {
	t.Helper()
	var v T
	require.NoError(t, json.Unmarshal(r.Data, &v), string(r.Data))
	return v
}

type sessionDTO struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func login(t *testing.T, app *fiber.App, googleToken string) sessionDTO {
	t.Helper()
	r := call(t, app, "POST", "/auth/google", "", map[string]string{"idToken": googleToken})
	require.Equal(t, 200, r.Status, r.Error)
	return decode[sessionDTO](t, r)
}

func createPlaylist(t *testing.T, app *fiber.App, token, name string) domain.Playlist {
	t.Helper()
	r := call(t, app, "POST", "/me/playlists", token, map[string]string{"name": name})
	require.Equal(t, 201, r.Status, r.Error)
	return decode[domain.Playlist](t, r)
}

func videoID(i int) string { return fmt.Sprintf("vid%08d", i) }

func addItem(t *testing.T, app *fiber.App, token, playlistID, vid string) apiResp {
	t.Helper()
	return call(t, app, "POST", "/me/playlists/"+playlistID+"/items", token, map[string]string{
		"videoId": vid, "title": "Title " + vid, "channelTitle": "Channel",
		"thumbnailUrl": "https://i.ytimg.com/vi/" + vid + "/hqdefault.jpg",
	})
}

func itemOrder(p domain.Playlist) []string {
	ids := make([]string, len(p.Items))
	for i, it := range p.Items {
		ids[i] = it.VideoID
	}
	return ids
}

func assertUnauthorized(t *testing.T, r apiResp) {
	t.Helper()
	assert.Equal(t, 401, r.Status)
	assert.False(t, r.Success)
	assert.Equal(t, "unauthorized", r.Error)
}

func TestAUTH01_InvalidTokensGet401(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, testJWTSecret)
	s := login(t, app, "google-a")

	t.Run("AUTH-01 forged google token", func(t *testing.T) {
		assertUnauthorized(t, call(t, app, "POST", "/auth/google", "", map[string]string{"idToken": "forged"}))
	})
	t.Run("AUTH-01 empty google token", func(t *testing.T) {
		assertUnauthorized(t, call(t, app, "POST", "/auth/google", "", map[string]string{"idToken": ""}))
	})
	t.Run("AUTH-01 missing bearer", func(t *testing.T) {
		assertUnauthorized(t, call(t, app, "GET", "/me", "", nil))
	})
	t.Run("AUTH-01 garbage bearer", func(t *testing.T) {
		assertUnauthorized(t, call(t, app, "GET", "/me/playlists", "garbage", nil))
	})
	t.Run("AUTH-01 expired session", func(t *testing.T) {
		expired, err := auth.NewSessionTokens(testJWTSecret, -time.Minute).Issue(s.User.ID)
		require.NoError(t, err)
		assertUnauthorized(t, call(t, app, "GET", "/me", expired, nil))
	})
	t.Run("AUTH-01 alg none", func(t *testing.T) {
		tok := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
			Issuer: "tunetrend", Subject: s.User.ID, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		})
		raw, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
		require.NoError(t, err)
		assertUnauthorized(t, call(t, app, "GET", "/me", raw, nil))
	})
	t.Run("AUTH-01 token signed with other secret", func(t *testing.T) {
		raw, err := auth.NewSessionTokens(strings.Repeat("x", 40), auth.SessionTTL).Issue(s.User.ID)
		require.NoError(t, err)
		assertUnauthorized(t, call(t, app, "GET", "/me", raw, nil))
	})
	t.Run("AUTH-01 non-bearer scheme", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Basic "+s.Token)
		res, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 401, res.StatusCode)
	})
	t.Run("valid session works", func(t *testing.T) {
		r := call(t, app, "GET", "/me", s.Token, nil)
		assert.Equal(t, 200, r.Status)
	})
}

func TestAUTH02_LoginCreatesThenReusesUser(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, testJWTSecret)

	first := login(t, app, "google-a")
	assert.NotEmpty(t, first.Token)
	assert.Equal(t, "Alice", first.User.Name)
	assert.Equal(t, "a@example.com", first.User.Email)

	second := login(t, app, "google-a-renamed")
	assert.Equal(t, first.User.ID, second.User.ID)
	assert.Equal(t, "Alice Renamed", second.User.Name)

	other := login(t, app, "google-b")
	assert.NotEqual(t, first.User.ID, other.User.ID)

	var count int64
	db.Model(&domain.User{}).Where("google_sub = ?", "sub-a").Count(&count)
	assert.Equal(t, int64(1), count)

	me := call(t, app, "GET", "/me", first.Token, nil)
	require.Equal(t, 200, me.Status)
	u := decode[map[string]any](t, me)
	assert.Equal(t, first.User.ID, u["id"])
	assert.NotContains(t, u, "googleSub")
	assert.NotContains(t, string(me.Data), "sub-a")
}

func TestAUTH03_DeleteMeRemovesUserPlaylistsItems(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, testJWTSecret)
	a := login(t, app, "google-a")
	b := login(t, app, "google-b")

	for i := 0; i < 2; i++ {
		p := createPlaylist(t, app, a.Token, fmt.Sprintf("A%d", i))
		require.Equal(t, 201, addItem(t, app, a.Token, p.ID, videoID(i)).Status)
	}
	pb := createPlaylist(t, app, b.Token, "B")
	require.Equal(t, 201, addItem(t, app, b.Token, pb.ID, videoID(9)).Status)

	r := call(t, app, "DELETE", "/me", a.Token, nil)
	require.Equal(t, 200, r.Status, r.Error)

	var users, playlists, items int64
	db.Model(&domain.User{}).Where("id = ?", a.User.ID).Count(&users)
	db.Model(&domain.Playlist{}).Where("user_id = ?", a.User.ID).Count(&playlists)
	db.Model(&domain.PlaylistItem{}).Where("playlist_id <> ?", pb.ID).Count(&items)
	assert.Zero(t, users)
	assert.Zero(t, playlists)
	assert.Zero(t, items)

	assertUnauthorized(t, call(t, app, "GET", "/me", a.Token, nil))
	assertUnauthorized(t, call(t, app, "GET", "/me/playlists", a.Token, nil))

	got := decode[domain.Playlist](t, call(t, app, "GET", "/me/playlists/"+pb.ID, b.Token, nil))
	assert.Equal(t, []string{videoID(9)}, itemOrder(got))

	again := login(t, app, "google-a")
	assert.NotEqual(t, a.User.ID, again.User.ID)
}

func TestPL01_OtherUsersPlaylistIs404(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, testJWTSecret)
	a := login(t, app, "google-a")
	b := login(t, app, "google-b")

	pa := createPlaylist(t, app, a.Token, "Mine")
	require.Equal(t, 201, addItem(t, app, a.Token, pa.ID, videoID(1)).Status)
	base := "/me/playlists/" + pa.ID

	cases := []struct {
		method, path string
		body         any
	}{
		{"GET", base, nil},
		{"PATCH", base, map[string]string{"name": "hacked"}},
		{"DELETE", base, nil},
		{"POST", base + "/items", map[string]string{"videoId": videoID(2), "title": "x"}},
		{"DELETE", base + "/items/" + videoID(1), nil},
		{"PUT", base + "/order", map[string][]string{"videoIds": {videoID(1)}}},
	}
	for _, tc := range cases {
		t.Run("PL-01 "+tc.method+" "+tc.path, func(t *testing.T) {
			r := call(t, app, tc.method, tc.path, b.Token, tc.body)
			assert.Equal(t, 404, r.Status)
			assert.False(t, r.Success)
		})
	}

	t.Run("PL-01 malformed id", func(t *testing.T) {
		assert.Equal(t, 404, call(t, app, "GET", "/me/playlists/not-a-uuid", a.Token, nil).Status)
	})

	list := decode[[]domain.Playlist](t, call(t, app, "GET", "/me/playlists", b.Token, nil))
	assert.Empty(t, list)

	got := decode[domain.Playlist](t, call(t, app, "GET", base, a.Token, nil))
	assert.Equal(t, "Mine", got.Name)
	assert.Equal(t, []string{videoID(1)}, itemOrder(got))
}

func TestPL02_ItemsOrderAndLimits(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, testJWTSecret)
	a := login(t, app, "google-a")
	p := createPlaylist(t, app, a.Token, "  Road trip  ")
	assert.Equal(t, "Road trip", p.Name)
	base := "/me/playlists/" + p.ID

	t.Run("PL-02 add keeps insertion order", func(t *testing.T) {
		for i := 1; i <= 4; i++ {
			r := addItem(t, app, a.Token, p.ID, videoID(i))
			require.Equal(t, 201, r.Status, r.Error)
			assert.Equal(t, i-1, decode[domain.PlaylistItem](t, r).Position)
		}
		got := decode[domain.Playlist](t, call(t, app, "GET", base, a.Token, nil))
		assert.Equal(t, []string{videoID(1), videoID(2), videoID(3), videoID(4)}, itemOrder(got))
		assert.Equal(t, int64(4), got.ItemCount)
		assert.Equal(t, "Title "+videoID(1), got.Items[0].Title)
		assert.Equal(t, "https://i.ytimg.com/vi/"+videoID(1)+"/hqdefault.jpg", got.Items[0].ThumbnailURL)
	})

	t.Run("PL-02 duplicate add is 409", func(t *testing.T) {
		assert.Equal(t, 409, addItem(t, app, a.Token, p.ID, videoID(1)).Status)
	})

	t.Run("PL-02 remove compacts positions", func(t *testing.T) {
		r := call(t, app, "DELETE", base+"/items/"+videoID(2), a.Token, nil)
		require.Equal(t, 200, r.Status, r.Error)
		got := decode[domain.Playlist](t, call(t, app, "GET", base, a.Token, nil))
		assert.Equal(t, []string{videoID(1), videoID(3), videoID(4)}, itemOrder(got))
		for i, it := range got.Items {
			assert.Equal(t, i, it.Position)
		}
		assert.Equal(t, 404, call(t, app, "DELETE", base+"/items/"+videoID(2), a.Token, nil).Status)
	})

	t.Run("PL-02 reorder", func(t *testing.T) {
		order := []string{videoID(4), videoID(1), videoID(3)}
		r := call(t, app, "PUT", base+"/order", a.Token, map[string][]string{"videoIds": order})
		require.Equal(t, 200, r.Status, r.Error)
		assert.Equal(t, order, itemOrder(decode[domain.Playlist](t, r)))
		got := decode[domain.Playlist](t, call(t, app, "GET", base, a.Token, nil))
		assert.Equal(t, order, itemOrder(got))
		for i, it := range got.Items {
			assert.Equal(t, i, it.Position)
		}
	})

	t.Run("PL-02 reorder must be permutation", func(t *testing.T) {
		bad := [][]string{
			{videoID(4), videoID(1)},
			{videoID(4), videoID(1), videoID(3), videoID(9)},
			{videoID(4), videoID(1), videoID(9)},
			{videoID(4), videoID(4), videoID(1)},
			{},
		}
		for _, ids := range bad {
			r := call(t, app, "PUT", base+"/order", a.Token, map[string][]string{"videoIds": ids})
			assert.Equal(t, 422, r.Status, "%v", ids)
		}
		got := decode[domain.Playlist](t, call(t, app, "GET", base, a.Token, nil))
		assert.Equal(t, []string{videoID(4), videoID(1), videoID(3)}, itemOrder(got))
	})

	t.Run("PL-02 invalid item input is 422", func(t *testing.T) {
		r := call(t, app, "POST", base+"/items", a.Token, map[string]string{"videoId": "bad", "title": "x"})
		assert.Equal(t, 422, r.Status)
		r = call(t, app, "POST", base+"/items", a.Token, map[string]string{
			"videoId": videoID(50), "title": "x", "thumbnailUrl": "https://evil.example.com/x.jpg",
		})
		assert.Equal(t, 422, r.Status)
	})

	t.Run("PL-02 rename and name limits", func(t *testing.T) {
		r := call(t, app, "PATCH", base, a.Token, map[string]string{"name": " ฟังตอนทำงาน "})
		require.Equal(t, 200, r.Status, r.Error)
		got := decode[domain.Playlist](t, r)
		assert.Equal(t, "ฟังตอนทำงาน", got.Name)
		assert.Equal(t, int64(3), got.ItemCount)

		assert.Equal(t, 422, call(t, app, "PATCH", base, a.Token, map[string]string{"name": "   "}).Status)
		assert.Equal(t, 422, call(t, app, "PATCH", base, a.Token, map[string]string{"name": strings.Repeat("ก", 61)}).Status)
		assert.Equal(t, 422, call(t, app, "POST", "/me/playlists", a.Token, map[string]string{"name": ""}).Status)
		assert.Equal(t, 201, call(t, app, "POST", "/me/playlists", a.Token, map[string]string{"name": strings.Repeat("ก", 60)}).Status)
	})

	t.Run("PL-02 200 items per playlist", func(t *testing.T) {
		full := createPlaylist(t, app, a.Token, "Full")
		items := make([]domain.PlaylistItem, 0, domain.MaxItemsPerPlaylist)
		for i := 0; i < domain.MaxItemsPerPlaylist; i++ {
			items = append(items, domain.PlaylistItem{
				PlaylistID: full.ID, VideoID: videoID(1000 + i), Position: i, Title: "t", AddedAt: time.Now(),
			})
		}
		require.NoError(t, db.Create(&items).Error)
		r := addItem(t, app, a.Token, full.ID, videoID(5000))
		assert.Equal(t, 422, r.Status)
		assert.False(t, r.Success)
	})

	t.Run("PL-02 50 playlists per user", func(t *testing.T) {
		var have int64
		db.Model(&domain.Playlist{}).Where("user_id = ?", a.User.ID).Count(&have)
		for i := int(have); i < domain.MaxPlaylistsPerUser; i++ {
			createPlaylist(t, app, a.Token, fmt.Sprintf("P%d", i))
		}
		r := call(t, app, "POST", "/me/playlists", a.Token, map[string]string{"name": "one too many"})
		assert.Equal(t, 422, r.Status)
		list := decode[[]domain.Playlist](t, call(t, app, "GET", "/me/playlists", a.Token, nil))
		assert.Len(t, list, domain.MaxPlaylistsPerUser)
	})

	t.Run("PL-02 delete playlist removes items", func(t *testing.T) {
		r := call(t, app, "DELETE", base, a.Token, nil)
		require.Equal(t, 200, r.Status, r.Error)
		var items int64
		db.Model(&domain.PlaylistItem{}).Where("playlist_id = ?", p.ID).Count(&items)
		assert.Zero(t, items)
		assert.Equal(t, 404, call(t, app, "GET", base, a.Token, nil).Status)
	})
}

func TestAuthNotConfiguredReturns500(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, "too-short")

	r := call(t, app, "POST", "/auth/google", "", map[string]string{"idToken": "google-a"})
	assert.Equal(t, 500, r.Status)
	r = call(t, app, "GET", "/me", "anything", nil)
	assert.Equal(t, 500, r.Status)
}

func TestMalformedBodyIs400(t *testing.T) {
	db := setupDB(t)
	app := newApp(db, testJWTSecret)
	a := login(t, app, "google-a")

	req := httptest.NewRequest("POST", "/me/playlists", strings.NewReader("{not json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.Token)
	res, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, nethttp.StatusBadRequest, res.StatusCode)
}

type unavailableVerifier struct{}

func (unavailableVerifier) Verify(context.Context, string) (*domain.GoogleIdentity, error) {
	return nil, domain.ErrAuthUnavailable
}

func TestGoogleLoginVerifierUnavailableIs503(t *testing.T) {
	tokens := auth.NewSessionTokens(testJWTSecret, auth.SessionTTL)
	authUC := usecase.NewAuthUsecase(unavailableVerifier{}, tokens, nil)
	app := fiber.New()
	deliveryhttp.RegisterAuthRoutes(app, authUC, nil, func(c *fiber.Ctx) error { return c.Next() })

	r := call(t, app, "POST", "/auth/google", "", map[string]string{"idToken": "google-a"})
	assert.Equal(t, 503, r.Status)
	assert.False(t, r.Success)
	assert.Equal(t, "auth_unavailable", r.Error)

	r = call(t, app, "POST", "/auth/google", "", map[string]string{"idToken": ""})
	assert.Equal(t, 401, r.Status)
}
