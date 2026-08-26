package usecase

import (
	"log"

	"tunetrend-backend/internal/domain"
)

type categoryVideoUsecase struct {
	repo       domain.CategoryVideoRepository
	ytRepo     domain.YouTubeRepository
	categoryID string
	label      string
}

func NewCategoryVideoUsecase(repo domain.CategoryVideoRepository, ytRepo domain.YouTubeRepository, categoryID, label string) domain.CategoryVideoUsecase {
	return &categoryVideoUsecase{repo: repo, ytRepo: ytRepo, categoryID: categoryID, label: label}
}

func (u *categoryVideoUsecase) SyncVideos(countryCode string) error {
	log.Printf("🔄 [Usecase] กำลังดึงข้อมูล %s สำหรับประเทศ %s...\n", u.label, countryCode)

	videos, err := u.ytRepo.FetchTrendingByCategory(countryCode, u.categoryID)
	if err != nil {
		return err
	}

	if err := u.repo.UpsertVideos(videos); err != nil {
		return err
	}

	log.Printf("✅ [Usecase] ซิงก์ %s ประเทศ %s สำเร็จ จำนวน %d รายการ\n", u.label, countryCode, len(videos))
	return nil
}
