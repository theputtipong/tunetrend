package usecase

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

type videoCategoryUsecase struct {
	categoryRepo    domain.VideoCategoryRepository
	ytRepo          domain.YouTubeRepository
	workerLogRepo   domain.WorkerLogRepository
	categoryConfigs []domain.CategoryVideoConfig
	gracePeriod     time.Duration
}

func NewVideoCategoryUsecase(
	categoryRepo domain.VideoCategoryRepository,
	ytRepo domain.YouTubeRepository,
	workerLogRepo domain.WorkerLogRepository,
	categoryConfigs []domain.CategoryVideoConfig,
	graceMinutes int,
) domain.VideoCategoryUsecase {
	return &videoCategoryUsecase{
		categoryRepo:    categoryRepo,
		ytRepo:          ytRepo,
		workerLogRepo:   workerLogRepo,
		categoryConfigs: categoryConfigs,
		gracePeriod:     time.Duration(graceMinutes) * time.Minute,
	}
}

func (u *videoCategoryUsecase) SyncCategories(countryCode string) error {
	log.Printf("🔄 [Usecase] กำลังดึงข้อมูลหมวดหมู่ YouTube สำหรับประเทศ %s...\n", countryCode)

	categories, err := u.ytRepo.FetchVideoCategories(countryCode)
	if err != nil {
		return err
	}

	tableNameByCategoryID := make(map[string]string, len(u.categoryConfigs)+1)
	jobNames := make([]string, 0, len(u.categoryConfigs)+1)
	for _, cfg := range u.categoryConfigs {
		tableNameByCategoryID[cfg.CategoryID] = cfg.TableName
		jobNames = append(jobNames, cfg.TableName)
	}
	tableNameByCategoryID[domain.MusicCategoryID] = domain.MusicSyncJobName
	jobNames = append(jobNames, domain.MusicSyncJobName)

	ids := make([]string, 0, len(categories))
	for _, c := range categories {
		ids = append(ids, c.ID)
	}

	latestStatuses, err := u.workerLogRepo.LatestStatuses(jobNames, countryCode)
	if err != nil {
		return err
	}

	since := time.Now().Add(-u.gracePeriod)
	recentSuccess, err := u.workerLogRepo.JobsWithSuccessSince(jobNames, countryCode, since)
	if err != nil {
		return err
	}

	existingReasons, err := u.categoryRepo.GetDeactivatedReasons(countryCode, ids)
	if err != nil {
		return err
	}

	setActiveBatch := make([]domain.VideoCategory, 0, len(categories))
	preserveActiveBatch := make([]domain.VideoCategory, 0, len(categories))

	for _, cat := range categories {
		reason := existingReasons[cat.ID]
		if reason != "" && reason != domain.DeactivatedReasonAutoFetchFailure {
			preserveActiveBatch = append(preserveActiveBatch, cat)
			continue
		}

		tableName, trackable := tableNameByCategoryID[cat.ID]

		switch {
		case !cat.Assignable || !trackable:
			cat.IsActive = false
			cat.DeactivatedReason = ""
			setActiveBatch = append(setActiveBatch, cat)

		case latestStatuses[tableName] != domain.WorkerLogStatusFailed:
			cat.IsActive = true
			cat.DeactivatedReason = ""
			setActiveBatch = append(setActiveBatch, cat)

		case recentSuccess[tableName]:
			preserveActiveBatch = append(preserveActiveBatch, cat)

		default:
			cat.IsActive = false
			cat.DeactivatedReason = domain.DeactivatedReasonAutoFetchFailure
			setActiveBatch = append(setActiveBatch, cat)
		}
	}

	if err := u.categoryRepo.UpsertCategoriesSetActive(setActiveBatch); err != nil {
		return err
	}
	if err := u.categoryRepo.UpsertCategoriesPreserveActive(preserveActiveBatch); err != nil {
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

	trackable := make(map[string]bool, len(u.categoryConfigs))
	for _, cfg := range u.categoryConfigs {
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
