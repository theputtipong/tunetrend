package usecase

import "tunetrend-backend/internal/domain"

const (
	discoverPerCountryLimit  = 2
	discoverPerCategoryLimit = 2
)

var discoverCountries = []string{"TH", "KR", "JP", "US", "GB"}

type discoverUsecase struct {
	categoryRepos   map[string]domain.CategoryVideoRepository
	categoryConfigs []domain.CategoryVideoConfig
}

func NewDiscoverUsecase(
	categoryRepos map[string]domain.CategoryVideoRepository,
	categoryConfigs []domain.CategoryVideoConfig,
) domain.DiscoverUsecase {
	return &discoverUsecase{categoryRepos: categoryRepos, categoryConfigs: categoryConfigs}
}

func (u *discoverUsecase) GetDiscoverItems() ([]domain.DiscoverItem, error) {
	items := make([]domain.DiscoverItem, 0)
	// A video can top several categories at once; each one appears only under the first.
	seenGlobal := make(map[string]bool)

	for _, cfg := range u.categoryConfigs {
		repo, ok := u.categoryRepos[cfg.CategoryID]
		if !ok {
			continue
		}

		videos, err := repo.GetTopAcrossCountries(discoverCountries, discoverPerCountryLimit)
		if err != nil {
			return nil, err
		}

		items = append(items, dedupeTopN(cfg, videos, discoverPerCategoryLimit, seenGlobal)...)
	}

	return items, nil
}

// Returning fewer than n items is accepted when earlier categories already took the rest.
func dedupeTopN(
	cfg domain.CategoryVideoConfig,
	videos []domain.CategoryVideo,
	n int,
	seenGlobal map[string]bool,
) []domain.DiscoverItem {
	items := make([]domain.DiscoverItem, 0, n)

	for _, v := range videos {
		if len(items) >= n {
			break
		}
		if seenGlobal[v.ID] {
			continue
		}
		seenGlobal[v.ID] = true

		items = append(items, domain.DiscoverItem{
			CategoryID:    cfg.CategoryID,
			CategoryLabel: cfg.Label,
			ID:            v.ID,
			Title:         v.Title,
			ChannelTitle:  v.ChannelTitle,
			ThumbnailURL:  v.ThumbnailURL,
			ViewCount:     v.ViewCount,
			CountryCode:   v.CountryCode,
		})
	}

	return items
}
