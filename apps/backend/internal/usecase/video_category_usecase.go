package usecase

import (
	"log"

	"tunetrend-backend/internal/domain"
)

type videoCategoryUsecase struct {
	categoryRepo domain.VideoCategoryRepository
	ytRepo       domain.YouTubeRepository
}

func NewVideoCategoryUsecase(categoryRepo domain.VideoCategoryRepository, ytRepo domain.YouTubeRepository) domain.VideoCategoryUsecase {
	return &videoCategoryUsecase{categoryRepo: categoryRepo, ytRepo: ytRepo}
}

func (u *videoCategoryUsecase) SyncCategories(countryCode string) error {
	log.Printf("🔄 [Usecase] กำลังดึงข้อมูลหมวดหมู่ YouTube สำหรับประเทศ %s...\n", countryCode)

	categories, err := u.ytRepo.FetchVideoCategories(countryCode)
	if err != nil {
		return err
	}

	if err := u.categoryRepo.UpsertCategories(categories); err != nil {
		return err
	}

	log.Printf("✅ [Usecase] ซิงก์หมวดหมู่ %s สำเร็จ จำนวน %d รายการ\n", countryCode, len(categories))
	return nil
}

func (u *videoCategoryUsecase) GetCategories(countryCode string) ([]domain.VideoCategory, error) {
	if countryCode == "" {
		countryCode = "TH"
	}

	categories, err := u.categoryRepo.GetActiveCategories(countryCode)
	if err != nil {
		return nil, err
	}

	// กันเหนียวอีกชั้น เผื่อแอดมินเผลอตั้ง is_active=true ให้หมวดที่ยังไม่มีตารางวิดีโอ sync จริง
	// (domain.CategoryVideoConfigs คือ source of truth ว่าหมวดไหนมีตาราง/worker sync ให้แล้ว)
	trackable := make(map[string]bool, len(domain.CategoryVideoConfigs))
	for _, cfg := range domain.CategoryVideoConfigs {
		trackable[cfg.CategoryID] = true
	}

	filtered := make([]domain.VideoCategory, 0, len(categories))
	for _, c := range categories {
		if trackable[c.ID] {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}
