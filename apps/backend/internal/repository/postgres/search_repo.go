package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	chartRankDepth   = maxListResults
	newReleaseWindow = 7 * 24 * time.Hour
	searchedWindow   = 7 * 24 * time.Hour
)

type searchRepository struct {
	db             *gorm.DB
	categoryTables map[string]string
}

func NewSearchRepository(db *gorm.DB, configs []domain.CategoryVideoConfig) domain.SearchRepository {
	tables := make(map[string]string, len(configs))
	for _, cfg := range configs {
		tables[cfg.CategoryID] = cfg.TableName
	}
	return &searchRepository{db: db, categoryTables: tables}
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func containsPattern(q string) string {
	return "%" + likeEscaper.Replace(q) + "%"
}

type chartRow struct {
	ID           string
	Title        string
	ChannelTitle string
	ThumbnailURL string
	ViewCount    string
	PublishedAt  time.Time
	VideoType    string
	TrendRank    sql.NullInt64
	NewRank      sql.NullInt64
}

// table is either "songs" or a category table name that database.Connect validated at boot, so concatenation is safe.
func (r *searchRepository) searchChartTable(table string, hasVideoType bool, country, pattern string, limit int) ([]chartRow, error) {
	videoType := "''"
	if hasVideoType {
		videoType = "t.video_type"
	}
	query := `
		SELECT t.id, t.title, t.channel_title, t.thumbnail_url, t.view_count, t.published_at,
			` + videoType + ` AS video_type, tr.rn AS trend_rank, nr.rn AS new_rank
		FROM ` + table + ` t
		LEFT JOIN (
			SELECT id, ROW_NUMBER() OVER (ORDER BY CAST(view_count AS BIGINT) DESC) AS rn
			FROM ` + table + ` WHERE country_code = @country
		) tr ON tr.id = t.id AND tr.rn <= @depth
		LEFT JOIN (
			SELECT id, ROW_NUMBER() OVER (ORDER BY CAST(view_count AS BIGINT) DESC) AS rn
			FROM ` + table + ` WHERE country_code = @country AND published_at >= @since
		) nr ON nr.id = t.id AND nr.rn <= @depth
		WHERE t.country_code = @country
			AND (t.title ILIKE @pattern ESCAPE '\' OR t.channel_title ILIKE @pattern ESCAPE '\')
		ORDER BY COALESCE(tr.rn, @depth + nr.rn, 2 * @depth + 1), CAST(t.view_count AS BIGINT) DESC, t.id
		LIMIT @limit
	`
	var rows []chartRow
	err := r.db.Raw(query, map[string]any{
		"country": country,
		"pattern": pattern,
		"depth":   chartRankDepth,
		"since":   time.Now().Add(-newReleaseWindow),
		"limit":   limit,
	}).Scan(&rows).Error
	return rows, err
}

func (r *searchRepository) SearchLocal(countryCode, categoryID, query string, limit int) ([]domain.SearchItem, error) {
	table, hasVideoType := "songs", true
	if categoryID != domain.MusicCategoryID {
		t, ok := r.categoryTables[categoryID]
		if !ok {
			return nil, fmt.Errorf("%w: %s", domain.ErrUnknownCategory, categoryID)
		}
		table, hasVideoType = t, false
	}

	pattern := containsPattern(query)
	chart, err := r.searchChartTable(table, hasVideoType, countryCode, pattern, limit)
	if err != nil {
		return nil, err
	}

	items := make([]domain.SearchItem, 0, limit)
	seen := make(map[string]bool, limit)
	for _, row := range chart {
		item := domain.SearchItem{
			ID:           row.ID,
			Title:        row.Title,
			ChannelTitle: row.ChannelTitle,
			ThumbnailURL: row.ThumbnailURL,
			ViewCount:    row.ViewCount,
			PublishedAt:  row.PublishedAt,
			VideoType:    row.VideoType,
			CategoryID:   categoryID,
			Source:       domain.SearchSourceChart,
		}
		switch {
		case row.TrendRank.Valid:
			item.ChartTab, item.ChartRank = domain.ChartTabTrending, int(row.TrendRank.Int64)
		case row.NewRank.Valid:
			item.ChartTab, item.ChartRank = domain.ChartTabNew, int(row.NewRank.Int64)
		}
		seen[row.ID] = true
		items = append(items, item)
	}

	if len(items) >= limit {
		return items, nil
	}

	var searched []domain.SearchVideo
	err = r.db.Where("category_id = ? AND (title ILIKE ? ESCAPE '\\' OR channel_title ILIKE ? ESCAPE '\\')", categoryID, pattern, pattern).
		Order("CAST(NULLIF(view_count, '') AS BIGINT) DESC NULLS LAST, id").
		Limit(limit + len(seen)).
		Find(&searched).Error
	if err != nil {
		return nil, err
	}

	for _, v := range searched {
		if len(items) >= limit {
			break
		}
		if seen[v.ID] {
			continue
		}
		seen[v.ID] = true
		items = append(items, domain.SearchItem{
			ID:           v.ID,
			Title:        v.Title,
			ChannelTitle: v.ChannelTitle,
			ThumbnailURL: v.ThumbnailURL,
			ViewCount:    v.ViewCount,
			PublishedAt:  v.PublishedAt,
			VideoType:    v.VideoType,
			CategoryID:   v.CategoryID,
			Source:       domain.SearchSourceSearched,
		})
	}
	return items, nil
}

func (r *searchRepository) UpsertSearchVideos(videos []domain.SearchVideo) error {
	if len(videos) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "channel_title", "thumbnail_url", "view_count", "published_at", "video_type", "category_id", "stats_updated_at",
		}),
	}).Create(&videos).Error
}

func (r *searchRepository) SearchVideoIDsWithHitsSince(since time.Time) ([]string, error) {
	var ids []string
	err := r.db.Raw(`
		SELECT v.id FROM search_videos v
		WHERE EXISTS (SELECT 1 FROM search_hits h WHERE h.video_id = v.id AND h.hit_at >= ?)
		ORDER BY v.id
	`, since).Scan(&ids).Error
	return ids, err
}

func (r *searchRepository) VideoExists(videoID string) (bool, error) {
	var exists bool
	err := r.db.Raw(`
		SELECT EXISTS (SELECT 1 FROM songs WHERE id = @id) OR EXISTS (SELECT 1 FROM search_videos WHERE id = @id)
	`, map[string]any{"id": videoID}).Scan(&exists).Error
	return exists, err
}

func (r *searchRepository) CreateHit(hit domain.SearchHit) error {
	return r.db.Create(&hit).Error
}

type searchedTrendRow struct {
	ID           string
	Hits         int64
	Rank         int
	PrevRank     sql.NullInt64
	Title        string
	ChannelTitle string
	ThumbnailURL string
	ViewCount    string
	PublishedAt  sql.NullTime
	VideoType    string
}

func (r *searchRepository) SearchedTrends(countryCode string, now time.Time, limit int) ([]domain.SearchedTrend, error) {
	query := `
		WITH music AS (
			SELECT id FROM songs WHERE country_code = @country
			UNION
			SELECT id FROM search_videos WHERE category_id = @music
		),
		cur AS (
			SELECT video_id, COUNT(*) AS hits, MAX(hit_at) AS last_hit FROM search_hits
			WHERE country_code = @country AND hit_at >= @curFrom AND hit_at <= @now AND video_id IN (SELECT id FROM music)
			GROUP BY video_id
		),
		prev AS (
			SELECT video_id, COUNT(*) AS hits, MAX(hit_at) AS last_hit FROM search_hits
			WHERE country_code = @country AND hit_at >= @prevFrom AND hit_at < @curFrom AND video_id IN (SELECT id FROM music)
			GROUP BY video_id
		),
		cur_ranked AS (
			SELECT video_id, hits, ROW_NUMBER() OVER (ORDER BY hits DESC, last_hit DESC, video_id) AS rn FROM cur
		),
		prev_ranked AS (
			SELECT video_id, ROW_NUMBER() OVER (ORDER BY hits DESC, last_hit DESC, video_id) AS rn FROM prev
		)
		SELECT c.video_id AS id, c.hits, c.rn AS rank, p.rn AS prev_rank,
			COALESCE(s.title, v.title, '') AS title,
			COALESCE(s.channel_title, v.channel_title, '') AS channel_title,
			COALESCE(s.thumbnail_url, v.thumbnail_url, '') AS thumbnail_url,
			COALESCE(s.view_count, v.view_count, '') AS view_count,
			COALESCE(s.published_at, v.published_at) AS published_at,
			COALESCE(s.video_type, v.video_type, '') AS video_type
		FROM cur_ranked c
		LEFT JOIN prev_ranked p ON p.video_id = c.video_id
		LEFT JOIN songs s ON s.id = c.video_id AND s.country_code = @country
		LEFT JOIN search_videos v ON v.id = c.video_id
		WHERE c.rn <= @limit
		ORDER BY c.rn
	`
	var rows []searchedTrendRow
	err := r.db.Raw(query, map[string]any{
		"country":  countryCode,
		"music":    domain.MusicCategoryID,
		"now":      now,
		"curFrom":  now.Add(-searchedWindow),
		"prevFrom": now.Add(-2 * searchedWindow),
		"limit":    limit,
	}).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	trends := make([]domain.SearchedTrend, len(rows))
	for i, row := range rows {
		trends[i] = domain.SearchedTrend{
			SearchItem: domain.SearchItem{
				ID:           row.ID,
				Title:        row.Title,
				ChannelTitle: row.ChannelTitle,
				ThumbnailURL: row.ThumbnailURL,
				ViewCount:    row.ViewCount,
				PublishedAt:  row.PublishedAt.Time,
				VideoType:    row.VideoType,
				CategoryID:   domain.MusicCategoryID,
				Source:       domain.SearchSourceSearched,
			},
			Rank: row.Rank,
			Hits: row.Hits,
		}
		if row.PrevRank.Valid {
			prev := int(row.PrevRank.Int64)
			trends[i].PrevRank = &prev
		}
	}
	return trends, nil
}
