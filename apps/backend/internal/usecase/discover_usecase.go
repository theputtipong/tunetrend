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
	// seenGlobal กันวิดีโอเดียวกันโผล่ซ้ำข้ามหมวดหมู่ (เช่นวิดีโอที่ YouTube
	// จัดให้ติด top ทั้ง Film & Animation และ Entertainment พร้อมกัน) — ไม่ใช่
	// แค่กันซ้ำภายในหมวดเดียวกันเท่านั้น
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

// dedupeTopN คัดวิดีโอไม่ให้ซ้ำกับที่เคยเลือกไปแล้ว (ทั้งในหมวดนี้และหมวดอื่น
// ก่อนหน้า ผ่าน seenGlobal) แล้วเอาแค่ N รายการแรกที่ยอดวิวสูงสุด (videos
// ถูกเรียงมาจากยอดวิวมากไปน้อยแล้ว) — อาจได้น้อยกว่า N ถ้าตัวเลือกที่เหลือ
// ถูกหมวดอื่นเอาไปแล้วหมด ถือว่ายอมรับได้ตามที่ตกลงไว้ (1-2 รายการต่อหมวด)
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
