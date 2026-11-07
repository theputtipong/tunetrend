package postgres_test

import (
	"fmt"
	"os"
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

func openSearchTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.Song{}, &domain.SearchVideo{}, &domain.SearchHit{}, &domain.SearchQuery{}, &domain.SearchQueryHit{}))
	require.NoError(t, db.Table("gaming_videos").AutoMigrate(&domain.CategoryVideo{}))
	require.NoError(t, db.Exec("TRUNCATE songs, search_videos, search_hits, gaming_videos, search_queries, search_query_hits").Error)
	t.Cleanup(func() {
		db.Exec("TRUNCATE songs, search_videos, search_hits, gaming_videos, search_queries, search_query_hits")
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	})
	return db
}

func newSearchRepo(db *gorm.DB) domain.SearchRepository {
	return postgres.NewSearchRepository(db, []domain.CategoryVideoConfig{{CategoryID: "20", TableName: "gaming_videos", Label: "Gaming"}})
}

func song(id, title, channel string, views int, publishedDaysAgo int) domain.Song {
	return domain.Song{
		ID: id, Title: title, ChannelTitle: channel, ViewCount: fmt.Sprint(views), CountryCode: "TH",
		CategoryID: "10", PublishedAt: time.Now().AddDate(0, 0, -publishedDaysAgo), VideoType: "MV",
	}
}

func TestSearchLocal_MusicSearchesSongsAndSearchVideosWithChartRank(t *testing.T) {
	db := openSearchTestDB(t)
	songRepo := postgres.NewSongRepository(db)
	repo := newSearchRepo(db)

	songs := []domain.Song{
		song("s1", "Top Hit", "Band", 9000, 30),
		song("s2", "ฝนตกที่หน้าต่าง (Official MV)", "วงสมมติ", 5000, 30),
		song("s3", "Other", "ฝนตก Channel", 100, 1),
	}
	require.NoError(t, songRepo.UpsertSongs(songs))
	require.NoError(t, songRepo.UpsertSongs([]domain.Song{{ID: "s2", Title: "ฝนตกที่หน้าต่าง (Official MV)", CountryCode: "US", ViewCount: "1", CategoryID: "10"}}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{
		{ID: "s2", Title: "ฝนตกที่หน้าต่าง dup", CategoryID: "10", ViewCount: "1"},
		{ID: "y1", Title: "ฟังเพลงตอนฝนตก", ChannelTitle: "ช่องรวมเพลง", CategoryID: "10", ViewCount: "3100000"},
		{ID: "y2", Title: "ฝนตก gameplay", CategoryID: "20", ViewCount: "5"},
	}))

	items, err := repo.SearchLocal("TH", "10", "ฝนตก", 30)
	require.NoError(t, err)
	require.Len(t, items, 3)

	assert.Equal(t, "s2", items[0].ID)
	assert.Equal(t, domain.SearchSourceChart, items[0].Source)
	assert.Equal(t, domain.ChartTabTrending, items[0].ChartTab)
	assert.Equal(t, 2, items[0].ChartRank)
	assert.Equal(t, "MV", items[0].VideoType)

	assert.Equal(t, "s3", items[1].ID, "channel_title match")
	assert.Equal(t, 3, items[1].ChartRank)

	assert.Equal(t, "y1", items[2].ID)
	assert.Equal(t, domain.SearchSourceSearched, items[2].Source)
	assert.Empty(t, items[2].ChartTab)
}

func TestSearchLocal_EscapesLikeWildcards(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	require.NoError(t, postgres.NewSongRepository(db).UpsertSongs([]domain.Song{
		song("p1", "100% Love", "A", 10, 1),
		song("p2", "1000 Love", "A", 10, 1),
		song("u1", "a_b", "A", 10, 1),
		song("u2", "axb", "A", 10, 1),
		song("b1", `back\slash`, "A", 10, 1),
	}))

	items, err := repo.SearchLocal("TH", "10", "0%", 30)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "p1", items[0].ID)

	items, err = repo.SearchLocal("TH", "10", "a_b", 30)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "u1", items[0].ID)

	items, err = repo.SearchLocal("TH", "10", `k\s`, 30)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "b1", items[0].ID)
}

func TestSearchLocal_SRCH08_CategorySearchesItsTable(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	require.NoError(t, postgres.NewCategoryVideoRepository(db, "gaming_videos").UpsertVideos([]domain.CategoryVideo{
		{ID: "g1", Title: "Minecraft Speedrun", ChannelTitle: "Gamer", ViewCount: "500", CountryCode: "TH", PublishedAt: time.Now()},
		{ID: "g2", Title: "Minecraft US", ViewCount: "500", CountryCode: "US", PublishedAt: time.Now()},
	}))
	require.NoError(t, postgres.NewSongRepository(db).UpsertSongs([]domain.Song{song("s1", "Minecraft song", "A", 10, 1)}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{
		{ID: "y20", Title: "minecraft tips", CategoryID: "20", ViewCount: "1"},
		{ID: "y10", Title: "minecraft music", CategoryID: "10", ViewCount: "1"},
	}))

	items, err := repo.SearchLocal("TH", "20", "minecraft", 30)
	require.NoError(t, err)
	ids := []string{}
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	assert.Equal(t, []string{"g1", "y20"}, ids)
	assert.Equal(t, "20", items[0].CategoryID)
	assert.Equal(t, 1, items[0].ChartRank)

	_, err = repo.SearchLocal("TH", "999", "minecraft", 30)
	assert.ErrorIs(t, err, domain.ErrUnknownCategory)
}

func TestUpsertSearchVideos_PreservesFirstSeenAt(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	first := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	later := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{{ID: "v", Title: "old", ViewCount: "1", FirstSeenAt: first, StatsUpdatedAt: first}}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{{ID: "v", Title: "new", ViewCount: "2", FirstSeenAt: later, StatsUpdatedAt: later}}))

	var got domain.SearchVideo
	require.NoError(t, db.First(&got, "id = ?", "v").Error)
	assert.Equal(t, "new", got.Title)
	assert.Equal(t, "2", got.ViewCount)
	assert.True(t, got.FirstSeenAt.Equal(first))
	assert.True(t, got.StatsUpdatedAt.Equal(later))
}

func TestSearchVideos_TRND01_NotVisibleInChartEndpoints(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	songRepo := postgres.NewSongRepository(db)
	require.NoError(t, songRepo.UpsertSongs([]domain.Song{song("chart1", "Chart (Official MV)", "A", 10, 1)}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{
		{ID: "yt1", Title: "Huge (Official MV)", ViewCount: "9999999999", CategoryID: "10", VideoType: "MV", PublishedAt: time.Now()},
	}))

	for name, fn := range map[string]func(string) ([]domain.Song, error){
		"/trends": songRepo.GetTrends, "/trends/new": songRepo.GetNewReleases, "/trends/mv": songRepo.GetMVs,
	} {
		songs, err := fn("TH")
		require.NoError(t, err, name)
		require.Len(t, songs, 1, name)
		assert.Equal(t, "chart1", songs[0].ID, name)
	}
}

func TestVideoExists(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	require.NoError(t, postgres.NewSongRepository(db).UpsertSongs([]domain.Song{song("s1", "x", "y", 1, 1)}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{{ID: "y1", CategoryID: "10"}}))

	for id, want := range map[string]bool{"s1": true, "y1": true, "nope": false} {
		got, err := repo.VideoExists(id)
		require.NoError(t, err)
		assert.Equal(t, want, got, id)
	}
}

func hits(t *testing.T, repo domain.SearchRepository, videoID, country string, at time.Time, n int) {
	for i := 0; i < n; i++ {
		require.NoError(t, repo.CreateHit(domain.SearchHit{VideoID: videoID, CountryCode: country, QueryNorm: "q", HitAt: at}))
	}
}

func TestSearchedTrends_TRND02_OrderedByHitsInLast7Days(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)

	require.NoError(t, postgres.NewSongRepository(db).UpsertSongs([]domain.Song{song("songA", "Song A", "Band", 100, 1)}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{
		{ID: "ytB", Title: "Video B", CategoryID: "10", ViewCount: "5", VideoType: "Lyric"},
		{ID: "ytC", Title: "Video C", CategoryID: "10"},
		{ID: "ytOld", Title: "Old", CategoryID: "10"},
		{ID: "ytGame", Title: "Game", CategoryID: "20"},
	}))

	hits(t, repo, "songA", "TH", now.Add(-1*time.Hour), 2)
	hits(t, repo, "ytB", "TH", now.Add(-6*24*time.Hour), 3)
	hits(t, repo, "ytC", "TH", now.Add(-2*time.Hour), 1)
	hits(t, repo, "ytOld", "TH", now.Add(-8*24*time.Hour), 10)
	hits(t, repo, "ytGame", "TH", now.Add(-time.Hour), 10)
	hits(t, repo, "ytC", "US", now.Add(-time.Hour), 10)

	hits(t, repo, "songA", "TH", now.Add(-9*24*time.Hour), 20)
	hits(t, repo, "ytB", "TH", now.Add(-10*24*time.Hour), 5)

	trends, err := repo.SearchedTrends("TH", now, 30)
	require.NoError(t, err)
	require.Len(t, trends, 3)

	assert.Equal(t, "ytB", trends[0].ID)
	assert.Equal(t, int64(3), trends[0].Hits)
	assert.Equal(t, 1, trends[0].Rank)
	require.NotNil(t, trends[0].PrevRank)
	assert.Equal(t, 3, *trends[0].PrevRank)
	assert.Equal(t, "Video B", trends[0].Title)
	assert.Equal(t, "Lyric", trends[0].VideoType)

	assert.Equal(t, "songA", trends[1].ID)
	assert.Equal(t, int64(2), trends[1].Hits)
	require.NotNil(t, trends[1].PrevRank)
	assert.Equal(t, 1, *trends[1].PrevRank)
	assert.Equal(t, "Song A", trends[1].Title)
	assert.Equal(t, "MV", trends[1].VideoType)

	assert.Equal(t, "ytC", trends[2].ID)
	assert.Nil(t, trends[2].PrevRank)

	limited, err := repo.SearchedTrends("TH", now, 2)
	require.NoError(t, err)
	assert.Len(t, limited, 2)
}

func TestSearchVideoIDsWithHitsSince(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Now()
	require.NoError(t, postgres.NewSongRepository(db).UpsertSongs([]domain.Song{song("songA", "A", "B", 1, 1)}))
	require.NoError(t, repo.UpsertSearchVideos([]domain.SearchVideo{{ID: "y1", CategoryID: "10"}, {ID: "y2", CategoryID: "10"}, {ID: "y3", CategoryID: "10"}}))
	hits(t, repo, "y1", "TH", now.Add(-time.Hour), 2)
	hits(t, repo, "y1", "US", now.Add(-time.Hour), 1)
	hits(t, repo, "y2", "TH", now.Add(-8*24*time.Hour), 1)
	hits(t, repo, "songA", "TH", now.Add(-time.Hour), 1)

	ids, err := repo.SearchVideoIDsWithHitsSince(now.Add(-7 * 24 * time.Hour))
	require.NoError(t, err)
	assert.Equal(t, []string{"y1"}, ids)
}

func queryHit(t *testing.T, repo domain.SearchRepository, country, q, visitor string, at time.Time) {
	t.Helper()
	require.NoError(t, repo.CreateQueryHit(domain.SearchQueryHit{CountryCode: country, QueryNorm: q, VisitorHash: visitor, HitAt: at}))
}

func TestSearchedQueries_TRND11_GroupsByQueryWithThresholdAndWindow(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	since := now.Add(-7 * 24 * time.Hour)

	require.NoError(t, repo.SaveQueryResults("TH", "rain song", []domain.SearchItem{{ID: "a", Title: "A"}, {ID: "b"}}, now))
	require.NoError(t, repo.SaveQueryResults("TH", "solo search", []domain.SearchItem{{ID: "c"}}, now))
	require.NoError(t, repo.SaveQueryResults("TH", "old song", []domain.SearchItem{{ID: "d"}}, now))
	require.NoError(t, repo.SaveQueryResults("TH", "busy", []domain.SearchItem{{ID: "e"}}, now))
	require.NoError(t, repo.SaveQueryResults("US", "rain song", []domain.SearchItem{{ID: "f"}}, now))

	queryHit(t, repo, "TH", "rain song", "v1", now.Add(-time.Hour))
	queryHit(t, repo, "TH", "rain song", "v2", now.Add(-2*time.Hour))
	queryHit(t, repo, "TH", "solo search", "v1", now.Add(-time.Hour))
	queryHit(t, repo, "TH", "solo search", "v1", now.Add(-3*time.Hour))
	queryHit(t, repo, "TH", "solo search", "v1", now.Add(-5*time.Hour))
	queryHit(t, repo, "TH", "old song", "v1", now.Add(-8*24*time.Hour))
	queryHit(t, repo, "TH", "old song", "v2", now.Add(-8*24*time.Hour))
	for i, v := range []string{"v1", "v2", "v3"} {
		queryHit(t, repo, "TH", "busy", v, now.Add(-time.Duration(i+1)*time.Minute))
	}
	queryHit(t, repo, "US", "rain song", "v9", now.Add(-time.Hour))
	queryHit(t, repo, "US", "rain song", "v8", now.Add(-time.Hour))

	got, err := repo.SearchedQueries("TH", since, 2, 30)
	require.NoError(t, err)
	require.Len(t, got, 2, "solo search has one visitor, old song is outside the window")
	assert.Equal(t, "busy", got[0].Query)
	assert.Equal(t, 1, got[0].Rank)
	assert.Equal(t, int64(3), got[0].Searches)
	assert.Equal(t, "rain song", got[1].Query)
	assert.Equal(t, 2, got[1].Rank)
	assert.Equal(t, int64(2), got[1].Searches)
	require.Len(t, got[1].Items, 2)
	assert.Equal(t, "A", got[1].Items[0].Title)

	limited, err := repo.SearchedQueries("TH", since, 2, 1)
	require.NoError(t, err)
	assert.Len(t, limited, 1)
}

func TestSaveQueryResults_LatestSearchReplacesGroupResults(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)

	require.NoError(t, repo.SaveQueryResults("TH", "rain", []domain.SearchItem{{ID: "old"}}, now.Add(-time.Hour)))
	require.NoError(t, repo.SaveQueryResults("TH", "rain", []domain.SearchItem{{ID: "new1"}, {ID: "new2"}}, now))
	queryHit(t, repo, "TH", "rain", "v1", now)
	queryHit(t, repo, "TH", "rain", "v2", now)

	got, err := repo.SearchedQueries("TH", now.Add(-24*time.Hour), 2, 30)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, []string{"new1", "new2"}, []string{got[0].Items[0].ID, got[0].Items[1].ID})

	var count int64
	db.Model(&domain.SearchQuery{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestSearchedQueries_GroupWithoutResultsHidden(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)

	require.NoError(t, repo.SaveQueryResults("TH", "empty", []domain.SearchItem{}, now))
	queryHit(t, repo, "TH", "empty", "v1", now)
	queryHit(t, repo, "TH", "empty", "v2", now)
	queryHit(t, repo, "TH", "no snapshot", "v1", now)
	queryHit(t, repo, "TH", "no snapshot", "v2", now)

	got, err := repo.SearchedQueries("TH", now.Add(-time.Hour), 2, 30)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestDeleteQueryHitsBefore_TRND12(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)

	queryHit(t, repo, "TH", "q", "v1", now.Add(-8*24*time.Hour))
	queryHit(t, repo, "TH", "q", "v2", now.Add(-6*24*time.Hour))

	require.NoError(t, repo.SaveQueryResults("TH", "q", []domain.SearchItem{{ID: "a"}}, now))
	queryHit(t, repo, "TH", "private old", "v3", now.Add(-9*24*time.Hour))
	require.NoError(t, repo.SaveQueryResults("TH", "private old", []domain.SearchItem{{ID: "b"}}, now.Add(-9*24*time.Hour)))
	require.NoError(t, repo.SaveQueryResults("TH", "never counted", []domain.SearchItem{{ID: "c"}}, now))

	n, err := repo.DeleteQueryHitsBefore(now.Add(-7 * 24 * time.Hour))
	require.NoError(t, err)
	assert.Equal(t, int64(2), n)
	var left int64
	db.Model(&domain.SearchQueryHit{}).Count(&left)
	assert.Equal(t, int64(1), left)

	var kept []string
	db.Model(&domain.SearchQuery{}).Order("query_norm").Pluck("query_norm", &kept)
	assert.Equal(t, []string{"q"}, kept, "groups without hits in the window are deleted with their snapshot")
}

func TestSearchedQueries_TieBreakByLatestSearchThenQuery(t *testing.T) {
	db := openSearchTestDB(t)
	repo := newSearchRepo(db)
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	for _, q := range []string{"alpha", "bravo", "charlie"} {
		require.NoError(t, repo.SaveQueryResults("TH", q, []domain.SearchItem{{ID: q}}, now))
	}
	queryHit(t, repo, "TH", "alpha", "v1", now.Add(-3*time.Hour))
	queryHit(t, repo, "TH", "alpha", "v2", now.Add(-3*time.Hour))
	queryHit(t, repo, "TH", "charlie", "v1", now.Add(-time.Hour))
	queryHit(t, repo, "TH", "charlie", "v2", now.Add(-3*time.Hour))
	queryHit(t, repo, "TH", "bravo", "v1", now.Add(-3*time.Hour))
	queryHit(t, repo, "TH", "bravo", "v2", now.Add(-3*time.Hour))

	got, err := repo.SearchedQueries("TH", now.Add(-24*time.Hour), 2, 30)
	require.NoError(t, err)
	require.Len(t, got, 3)
	assert.Equal(t, []string{"charlie", "alpha", "bravo"}, []string{got[0].Query, got[1].Query, got[2].Query})
}
