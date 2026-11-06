package ytquota

import "tunetrend-backend/internal/domain"

type Recorder interface {
	Record(units int64, source string)
}

type countingYouTube struct {
	inner    domain.YouTubeRepository
	recorder Recorder
}

func WrapYouTube(inner domain.YouTubeRepository, recorder Recorder) domain.YouTubeRepository {
	return &countingYouTube{inner: inner, recorder: recorder}
}

func (c *countingYouTube) FetchTrending(countryCode string) ([]domain.Song, error) {
	c.recorder.Record(VideosListUnits, SourceWorker)
	return c.inner.FetchTrending(countryCode)
}

func (c *countingYouTube) FetchVideoCategories(countryCode string) ([]domain.VideoCategory, error) {
	c.recorder.Record(1, SourceWorker)
	return c.inner.FetchVideoCategories(countryCode)
}

func (c *countingYouTube) FetchTrendingByCategory(countryCode, categoryID string) ([]domain.CategoryVideo, error) {
	c.recorder.Record(VideosListUnits, SourceWorker)
	return c.inner.FetchTrendingByCategory(countryCode, categoryID)
}

func (c *countingYouTube) SearchVideos(query, countryCode, categoryID string) ([]domain.SearchVideo, error) {
	return c.inner.SearchVideos(query, countryCode, categoryID)
}

func (c *countingYouTube) FetchVideosByIDs(ids []string) ([]domain.SearchVideo, error) {
	batches := (len(ids) + VideosBatchSize - 1) / VideosBatchSize
	c.recorder.Record(int64(batches*VideosListUnits), SourceWorker)
	return c.inner.FetchVideosByIDs(ids)
}
