package usecase

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"tunetrend-backend/internal/domain"
)

const (
	searchQueryMinRunes   = 2
	searchQueryMaxRunes   = 100
	localSearchLimit      = 30
	suggestYouTubeBelow   = 5
	searchedTrendsLimit   = 30
	searchedTrendsWindow  = 7 * 24 * time.Hour
	searchCacheKeyPattern = "yt:search:%s:%s:%s"

	searchedQueriesMinVisitors = 2
	searchedQueriesLimit       = 30
	searchedQueryItemsLimit    = 5
)

var videoIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{11}$`)

type SearchDeps struct {
	Repo            domain.SearchRepository
	YouTube         domain.YouTubeRepository
	Guard           domain.YouTubeQuotaGuard
	Cache           domain.SearchCache
	Deduper         domain.HitDeduper
	CategoryRepo    domain.VideoCategoryRepository
	CategoryConfigs []domain.CategoryVideoConfig
	Countries       []string
	VisitorKey      []byte
	Now             func() time.Time
}

type searchUsecase struct {
	deps        SearchDeps
	countries   map[string]bool
	categoryIDs map[string]bool
}

func NewSearchUsecase(deps SearchDeps) domain.SearchUsecase {
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if len(deps.VisitorKey) == 0 {
		log.Println("⚠️ [Search] ไม่มี VisitorKey (JWT_SECRET) ใช้ key สุ่มต่อ process แทน — รีสตาร์ตแล้วคนเดิมจะถูกนับเป็นคนใหม่")
		deps.VisitorKey = make([]byte, 32)
		if _, err := rand.Read(deps.VisitorKey); err != nil {
			panic(fmt.Sprintf("generate visitor key: %v", err))
		}
	}
	countries := make(map[string]bool, len(deps.Countries))
	for _, c := range deps.Countries {
		countries[strings.ToUpper(strings.TrimSpace(c))] = true
	}
	categoryIDs := make(map[string]bool, len(deps.CategoryConfigs))
	for _, cfg := range deps.CategoryConfigs {
		categoryIDs[cfg.CategoryID] = true
	}
	return &searchUsecase{deps: deps, countries: countries, categoryIDs: categoryIDs}
}

type searchInput struct {
	query      string
	queryNorm  string
	country    string
	categoryID string
}

func normalizeQuery(q string) string {
	return strings.ToLower(strings.Join(strings.Fields(q), " "))
}

func (u *searchUsecase) validateCountry(country string) (string, error) {
	country = strings.ToUpper(strings.TrimSpace(country))
	if !u.countries[country] {
		return "", &domain.ValidationError{Message: "unsupported country"}
	}
	return country, nil
}

func (u *searchUsecase) validate(query, country, categoryID string) (searchInput, error) {
	q := strings.TrimSpace(query)
	if n := utf8.RuneCountInString(q); n < searchQueryMinRunes || n > searchQueryMaxRunes {
		return searchInput{}, &domain.ValidationError{Message: fmt.Sprintf("q must be %d-%d characters", searchQueryMinRunes, searchQueryMaxRunes)}
	}

	country, err := u.validateCountry(country)
	if err != nil {
		return searchInput{}, err
	}

	categoryID = strings.TrimSpace(categoryID)
	if categoryID == "" {
		categoryID = domain.MusicCategoryID
	}
	if categoryID != domain.MusicCategoryID {
		if err := u.validateCategory(country, categoryID); err != nil {
			return searchInput{}, err
		}
	}

	return searchInput{query: q, queryNorm: normalizeQuery(q), country: country, categoryID: categoryID}, nil
}

func (u *searchUsecase) validateCategory(country, categoryID string) error {
	invalid := &domain.ValidationError{Message: "unknown or inactive category"}
	if !u.categoryIDs[categoryID] {
		return invalid
	}
	active, err := u.deps.CategoryRepo.GetActiveCategories(country)
	if err != nil {
		return err
	}
	for _, c := range active {
		if c.ID == categoryID {
			return nil
		}
	}
	return invalid
}

func (u *searchUsecase) SearchLocal(ip, query, countryCode, categoryID string) (domain.LocalSearchResult, error) {
	in, err := u.validate(query, countryCode, categoryID)
	if err != nil {
		return domain.LocalSearchResult{}, err
	}

	items, err := u.deps.Repo.SearchLocal(in.country, in.categoryID, in.query, localSearchLimit)
	if err != nil {
		return domain.LocalSearchResult{}, err
	}
	if items == nil {
		items = []domain.SearchItem{}
	}

	u.countQuerySearch(ip, in, items)

	return domain.LocalSearchResult{Items: items, SuggestYouTube: len(items) < suggestYouTubeBelow}, nil
}

func (u *searchUsecase) visitorHash(ip string) string {
	mac := hmac.New(sha256.New, u.deps.VisitorKey)
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))[:32]
}

func firstItems(items []domain.SearchItem, n int) []domain.SearchItem {
	if len(items) > n {
		items = items[:n]
	}
	return append([]domain.SearchItem{}, items...)
}

func (u *searchUsecase) countQuerySearch(ip string, in searchInput, items []domain.SearchItem) {
	if ip == "" || in.categoryID != domain.MusicCategoryID || len(items) == 0 {
		return
	}
	visitor := u.visitorHash(ip)
	first, err := u.deps.Deduper.FirstSearch(visitor, in.country, in.queryNorm)
	if err != nil {
		log.Printf("⚠️ [Search] ตรวจการค้นซ้ำไม่สำเร็จ ข้ามการนับคำค้น: %v", err)
		return
	}
	if !first {
		return
	}

	now := u.deps.Now()
	if err := u.deps.Repo.SaveQueryResults(in.country, in.queryNorm, firstItems(items, searchedQueryItemsLimit), now); err != nil {
		log.Printf("⚠️ [Search] บันทึกผลของคำค้นไม่สำเร็จ: %v", err)
		return
	}
	hit := domain.SearchQueryHit{CountryCode: in.country, QueryNorm: in.queryNorm, VisitorHash: visitor, HitAt: now}
	if err := u.deps.Repo.CreateQueryHit(hit); err != nil {
		log.Printf("⚠️ [Search] บันทึกการค้นไม่สำเร็จ: %v", err)
	}
}

func (u *searchUsecase) refreshQueryResults(in searchInput) {
	if in.categoryID != domain.MusicCategoryID {
		return
	}
	items, err := u.deps.Repo.SearchLocal(in.country, in.categoryID, in.query, searchedQueryItemsLimit)
	if err != nil {
		log.Printf("⚠️ [Search] โหลดผลค้นใหม่ของคำค้นไม่สำเร็จ: %v", err)
		return
	}
	if len(items) == 0 {
		return
	}
	if err := u.deps.Repo.SaveQueryResults(in.country, in.queryNorm, firstItems(items, searchedQueryItemsLimit), u.deps.Now()); err != nil {
		log.Printf("⚠️ [Search] บันทึกผลของคำค้นไม่สำเร็จ: %v", err)
	}
}

func (u *searchUsecase) quotaError(err error) error {
	return &domain.QuotaError{Err: err, RetryAt: u.deps.Guard.RetryAt()}
}

func (u *searchUsecase) quotaErrorWithRemaining(userID string, err error) error {
	qe := &domain.QuotaError{Err: err, RetryAt: u.deps.Guard.RetryAt()}
	if remaining, rerr := u.deps.Guard.Remaining(userID); rerr == nil {
		qe.Remaining = remaining
	}
	return qe
}

func (u *searchUsecase) SearchYouTube(userID, query, countryCode, categoryID string) (domain.YouTubeSearchResult, error) {
	in, err := u.validate(query, countryCode, categoryID)
	if err != nil {
		return domain.YouTubeSearchResult{}, err
	}

	cacheKey := fmt.Sprintf(searchCacheKeyPattern, in.country, in.categoryID, in.queryNorm)
	cached, hit, err := u.deps.Cache.Get(cacheKey)
	if err != nil {
		log.Printf("⚠️ [Search] อ่านแคชไม่สำเร็จ (fail-closed): %v", err)
		return domain.YouTubeSearchResult{}, u.quotaErrorWithRemaining(userID, domain.ErrQuotaUnavailable)
	}
	if hit {
		remaining, err := u.deps.Guard.Remaining(userID)
		if err != nil {
			return domain.YouTubeSearchResult{}, u.quotaError(err)
		}
		return domain.YouTubeSearchResult{Items: cached, Remaining: remaining, Cached: true}, nil
	}

	remaining, err := u.deps.Guard.Reserve(userID)
	if err != nil {
		return domain.YouTubeSearchResult{}, u.quotaErrorWithRemaining(userID, err)
	}

	videos, err := u.deps.YouTube.SearchVideos(in.query, in.country, in.categoryID)
	if err != nil {
		u.deps.Guard.Release(userID)
		return domain.YouTubeSearchResult{}, fmt.Errorf("youtube search: %w", err)
	}

	now := u.deps.Now()
	items := make([]domain.SearchItem, len(videos))
	for i := range videos {
		if videos[i].CategoryID == "" {
			videos[i].CategoryID = in.categoryID
		}
		videos[i].FirstSeenAt = now
		videos[i].StatsUpdatedAt = now
		items[i] = searchVideoToItem(videos[i], domain.SearchSourceYouTube)
	}

	if err := u.deps.Repo.UpsertSearchVideos(videos); err != nil {
		log.Printf("⚠️ [Search] บันทึก search_videos ไม่สำเร็จ: %v", err)
	}
	if err := u.deps.Cache.Set(cacheKey, items); err != nil {
		log.Printf("⚠️ [Search] เขียนแคชไม่สำเร็จ: %v", err)
	}
	u.refreshQueryResults(in)

	return domain.YouTubeSearchResult{Items: items, Remaining: remaining}, nil
}

func searchVideoToItem(v domain.SearchVideo, source string) domain.SearchItem {
	return domain.SearchItem{
		ID:           v.ID,
		Title:        v.Title,
		ChannelTitle: v.ChannelTitle,
		ThumbnailURL: v.ThumbnailURL,
		ViewCount:    v.ViewCount,
		PublishedAt:  v.PublishedAt,
		VideoType:    v.VideoType,
		CategoryID:   v.CategoryID,
		Source:       source,
	}
}

func truncateRunes(s string, max int) string {
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	return string([]rune(s)[:max])
}

func (u *searchUsecase) RecordHit(ip, videoID, countryCode, query string) (bool, error) {
	if !videoIDPattern.MatchString(videoID) {
		return false, &domain.ValidationError{Message: "invalid videoId"}
	}
	country, err := u.validateCountry(countryCode)
	if err != nil {
		return false, err
	}

	exists, err := u.deps.Repo.VideoExists(videoID)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, domain.ErrUnknownVideo
	}

	first, err := u.deps.Deduper.FirstHit(ip, videoID)
	if err != nil {
		log.Printf("⚠️ [Search] ตรวจ hit ซ้ำไม่สำเร็จ ข้ามการนับ: %v", err)
		return false, nil
	}
	if !first {
		return false, nil
	}

	err = u.deps.Repo.CreateHit(domain.SearchHit{
		VideoID:     videoID,
		CountryCode: country,
		QueryNorm:   truncateRunes(normalizeQuery(query), searchQueryMaxRunes),
		HitAt:       u.deps.Now(),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *searchUsecase) GetSearchedTrends(countryCode string) ([]domain.SearchedTrend, error) {
	country, err := u.validateCountry(countryCode)
	if err != nil {
		return nil, err
	}

	trends, err := u.deps.Repo.SearchedTrends(country, u.deps.Now(), searchedTrendsLimit)
	if err != nil {
		return nil, err
	}
	if trends == nil {
		trends = []domain.SearchedTrend{}
	}

	for i := range trends {
		if trends[i].PrevRank == nil {
			trends[i].IsNew = true
			trends[i].RankChange = 0
			continue
		}
		trends[i].RankChange = *trends[i].PrevRank - trends[i].Rank
	}
	return trends, nil
}

func (u *searchUsecase) GetSearchedQueries(countryCode string) ([]domain.SearchedQuery, error) {
	country, err := u.validateCountry(countryCode)
	if err != nil {
		return nil, err
	}
	groups, err := u.deps.Repo.SearchedQueries(country, u.deps.Now().Add(-searchedTrendsWindow), searchedQueriesMinVisitors, searchedQueriesLimit)
	if err != nil {
		return nil, err
	}
	if groups == nil {
		groups = []domain.SearchedQuery{}
	}
	return groups, nil
}

func (u *searchUsecase) PruneSearchQueryHits() (int64, error) {
	return u.deps.Repo.DeleteQueryHitsBefore(u.deps.Now().Add(-searchedTrendsWindow))
}

func (u *searchUsecase) RefreshSearchVideoStats() (int, error) {
	now := u.deps.Now()
	ids, err := u.deps.Repo.SearchVideoIDsWithHitsSince(now.Add(-searchedTrendsWindow))
	if err != nil {
		return 0, err
	}
	if len(ids) == 0 {
		return 0, nil
	}

	videos, err := u.deps.YouTube.FetchVideosByIDs(ids)
	if err != nil {
		return 0, err
	}

	for i := range videos {
		videos[i].FirstSeenAt = now
		videos[i].StatsUpdatedAt = now
	}
	if err := u.deps.Repo.UpsertSearchVideos(videos); err != nil {
		return 0, err
	}
	return len(videos), nil
}
