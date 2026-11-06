package postgres_test

import (
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/repository/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const userTestSchema = "tt_user_repo_test"

func openUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}
	admin, err := gorm.Open(gormpostgres.New(gormpostgres.Config{DSN: dsn, PreferSimpleProtocol: true}), cfg)
	require.NoError(t, err)
	require.NoError(t, admin.Exec("CREATE SCHEMA IF NOT EXISTS "+userTestSchema).Error)
	if sqlDB, err := admin.DB(); err == nil {
		sqlDB.Close()
	}

	u, err := url.Parse(dsn)
	require.NoError(t, err)
	q := u.Query()
	q.Set("search_path", userTestSchema)
	q.Set("application_name", userTestSchema)
	u.RawQuery = q.Encode()
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{DSN: u.String(), PreferSimpleProtocol: true}), cfg)
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.User{}, &domain.Playlist{}, &domain.PlaylistItem{}))
	require.NoError(t, db.Exec("TRUNCATE users, playlists, playlist_items").Error)
	t.Cleanup(func() {
		db.Exec("TRUNCATE users, playlists, playlist_items")
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})
	return db
}

func createTestUser(t *testing.T, db *gorm.DB, sub string) *domain.User {
	t.Helper()
	u, err := postgres.NewUserRepository(db).UpsertByGoogleSub(domain.GoogleIdentity{Subject: sub, Email: sub + "@example.com"}, time.Now())
	require.NoError(t, err)
	return u
}

func waitForLockWaiter(t *testing.T, db *gorm.DB) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var n int64
		require.NoError(t, db.Raw(`SELECT COUNT(*) FROM pg_stat_activity WHERE datname = current_database() AND application_name = ? AND wait_event_type = 'Lock'`, userTestSchema).Scan(&n).Error)
		if n > 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("DeleteWithData never blocked on a row lock")
}

func countOrphans(t *testing.T, db *gorm.DB) (playlists, items int64) {
	t.Helper()
	require.NoError(t, db.Raw(`SELECT COUNT(*) FROM playlists p WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = p.user_id)`).Scan(&playlists).Error)
	require.NoError(t, db.Raw(`SELECT COUNT(*) FROM playlist_items i WHERE NOT EXISTS (SELECT 1 FROM playlists p WHERE p.id = i.playlist_id)`).Scan(&items).Error)
	return
}

func TestDeleteWithData_ConcurrentPlaylistCreateLeavesNoOrphans(t *testing.T) {
	db := openUserTestDB(t)
	u := createTestUser(t, db, "sub-del-create")
	users := postgres.NewUserRepository(db)

	tx := db.Begin()
	require.NoError(t, tx.Error)
	require.NoError(t, tx.Exec(`SELECT id FROM users WHERE id = ? FOR UPDATE`, u.ID).Error)

	var wg sync.WaitGroup
	var delErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		delErr = users.DeleteWithData(u.ID)
	}()
	waitForLockWaiter(t, db)

	require.NoError(t, tx.Create(&domain.Playlist{UserID: u.ID, Name: "late"}).Error)
	require.NoError(t, tx.Commit().Error)
	wg.Wait()
	require.NoError(t, delErr)

	playlists, items := countOrphans(t, db)
	assert.Zero(t, playlists)
	assert.Zero(t, items)
}

func TestDeleteWithData_ConcurrentItemAddLeavesNoOrphans(t *testing.T) {
	db := openUserTestDB(t)
	u := createTestUser(t, db, "sub-del-item")
	p := &domain.Playlist{UserID: u.ID, Name: "mix"}
	require.NoError(t, db.Create(p).Error)
	users := postgres.NewUserRepository(db)

	tx := db.Begin()
	require.NoError(t, tx.Error)
	require.NoError(t, tx.Exec(`SELECT id FROM playlists WHERE id = ? FOR UPDATE`, p.ID).Error)

	var wg sync.WaitGroup
	var delErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		delErr = users.DeleteWithData(u.ID)
	}()
	waitForLockWaiter(t, db)

	require.NoError(t, tx.Create(&domain.PlaylistItem{PlaylistID: p.ID, VideoID: "dQw4w9WgXcQ", AddedAt: time.Now()}).Error)
	require.NoError(t, tx.Commit().Error)
	wg.Wait()
	require.NoError(t, delErr)

	playlists, items := countOrphans(t, db)
	assert.Zero(t, playlists)
	assert.Zero(t, items)
}

func TestDeleteWithData_WritesAfterDeletionFail(t *testing.T) {
	db := openUserTestDB(t)
	u := createTestUser(t, db, "sub-del-after")
	playlists := postgres.NewPlaylistRepository(db)
	p, err := playlists.CreateWithLimit(u.ID, "mix", domain.MaxPlaylistsPerUser)
	require.NoError(t, err)

	require.NoError(t, postgres.NewUserRepository(db).DeleteWithData(u.ID))

	_, err = playlists.CreateWithLimit(u.ID, "again", domain.MaxPlaylistsPerUser)
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
	_, err = playlists.AddItemWithLimit(u.ID, p.ID, domain.PlaylistItem{VideoID: "dQw4w9WgXcQ"}, domain.MaxItemsPerPlaylist)
	assert.ErrorIs(t, err, domain.ErrPlaylistNotFound)
	assert.ErrorIs(t, postgres.NewUserRepository(db).DeleteWithData(u.ID), domain.ErrUserNotFound)
}

func TestDeleteWithData_RacingWritersLeaveNoOrphans(t *testing.T) {
	db := openUserTestDB(t)
	users := postgres.NewUserRepository(db)
	playlists := postgres.NewPlaylistRepository(db)

	for round := 0; round < 20; round++ {
		u := createTestUser(t, db, "sub-race-"+time.Now().Format("150405.000000000"))
		p, err := playlists.CreateWithLimit(u.ID, "base", domain.MaxPlaylistsPerUser)
		require.NoError(t, err)

		var wg sync.WaitGroup
		start := make(chan struct{})
		for i := 0; i < 4; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				<-start
				_, _ = playlists.CreateWithLimit(u.ID, "x", domain.MaxPlaylistsPerUser)
			}()
			go func(i int) {
				defer wg.Done()
				<-start
				_, _ = playlists.AddItemWithLimit(u.ID, p.ID, domain.PlaylistItem{VideoID: "vid0000000" + string(rune('a'+i))}, domain.MaxItemsPerPlaylist)
			}(i)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_ = users.DeleteWithData(u.ID)
		}()
		close(start)
		wg.Wait()
	}

	orphanPlaylists, orphanItems := countOrphans(t, db)
	assert.Zero(t, orphanPlaylists)
	assert.Zero(t, orphanItems)
}
