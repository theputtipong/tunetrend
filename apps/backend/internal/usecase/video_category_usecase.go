package usecase

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

// categoryFetchFailureGracePeriod คือระยะเวลาที่ยอมให้ content-fetch ของหมวดหมู่ fail ต่อเนื่องได้
// ก่อนที่ sync จะบังคับปิด is_active ให้เอง
const categoryFetchFailureGracePeriod = 24 * time.Hour

type videoCategoryUsecase struct {
	categoryRepo  domain.VideoCategoryRepository
	ytRepo        domain.YouTubeRepository
	workerLogRepo domain.WorkerLogRepository
}

func NewVideoCategoryUsecase(
	categoryRepo domain.VideoCategoryRepository,
	ytRepo domain.YouTubeRepository,
	workerLogRepo domain.WorkerLogRepository,
) domain.VideoCategoryUsecase {
	return &videoCategoryUsecase{
		categoryRepo:  categoryRepo,
		ytRepo:        ytRepo,
		workerLogRepo: workerLogRepo,
	}
}

func (u *videoCategoryUsecase) SyncCategories(countryCode string) error {
	log.Printf("🔄 [Usecase] กำลังดึงข้อมูลหมวดหมู่ YouTube สำหรับประเทศ %s...\n", countryCode)

	categories, err := u.ytRepo.FetchVideoCategories(countryCode)
	if err != nil {
		return err
	}

	// ตารางที่มี worker content-sync ให้จริง (source of truth: domain.CategoryVideoConfigs)
	tableNameByCategoryID := make(map[string]string, len(domain.CategoryVideoConfigs))
	jobNames := make([]string, 0, len(domain.CategoryVideoConfigs))
	for _, cfg := range domain.CategoryVideoConfigs {
		tableNameByCategoryID[cfg.CategoryID] = cfg.TableName
		jobNames = append(jobNames, cfg.TableName)
	}

	ids := make([]string, 0, len(categories))
	for _, c := range categories {
		ids = append(ids, c.ID)
	}

	latestStatuses, err := u.workerLogRepo.LatestStatuses(jobNames, countryCode)
	if err != nil {
		return err
	}

	since := time.Now().Add(-categoryFetchFailureGracePeriod)
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
		// เคยถูกปิดด้วยเหตุผลอื่นที่ไม่ใช่ auto (เช่นแอดมินตั้งเอง) — sync จะไม่แตะ is_active ของหมวดนี้อีกเลย
		reason := existingReasons[cat.ID]
		if reason != "" && reason != domain.DeactivatedReasonAutoFetchFailure {
			preserveActiveBatch = append(preserveActiveBatch, cat)
			continue
		}

		tableName, trackable := tableNameByCategoryID[cat.ID]

		switch {
		case !cat.Assignable || !trackable:
			// assignable=false, หรือ assignable=true แต่ยังไม่มีตาราง/worker content sync ให้หมวดนี้เลย
			cat.IsActive = false
			cat.DeactivatedReason = ""
			setActiveBatch = append(setActiveBatch, cat)

		case latestStatuses[tableName] != domain.WorkerLogStatusFailed:
			// ไม่เคย sync content เลย (key ไม่พบ -> "") หรือ sync ล่าสุดสำเร็จ -> เปิดใช้งาน (self-healing)
			cat.IsActive = true
			cat.DeactivatedReason = ""
			setActiveBatch = append(setActiveBatch, cat)

		case recentSuccess[tableName]:
			// กำลัง fail อยู่ตอนนี้ แต่ยังมี success ภายใน 24 ชม.ล่าสุด -> ยังไม่ครบ grace period, ค้างค่าเดิมไว้
			preserveActiveBatch = append(preserveActiveBatch, cat)

		default:
			// fail ต่อเนื่องเกิน 24 ชม. โดยไม่มี success เลยในช่วงนั้น -> บังคับปิด พร้อมระบุเหตุผลว่าระบบปิดเอง
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

	// กันเหนียวอีกชั้น เผื่อมีหมวดที่ is_active=true แต่ไม่มีตารางวิดีโอ sync จริง
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
