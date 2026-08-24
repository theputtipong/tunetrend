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
