package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

const (
	jobNameVideoCategorySync = "video_category_sync"
	categorySyncInterval     = 24 * time.Hour
)

var categorySyncCountries = []string{"TH", "KR", "JP", "US", "GB"}

func StartVideoCategorySync(uc domain.VideoCategoryUsecase, logRepo domain.WorkerLogRepository) {
	log.Println("🚀 [Worker] เริ่มต้นระบบ Sync หมวดหมู่ YouTube แล้ว...")

	syncCategories(uc, logRepo)

	ticker := time.NewTicker(categorySyncInterval)
	for range ticker.C {
		syncCategories(uc, logRepo)
	}
}

func syncCategories(uc domain.VideoCategoryUsecase, logRepo domain.WorkerLogRepository) {
	for _, country := range categorySyncCountries {
		startedAt := time.Now()
		err := uc.SyncCategories(country)
		finishedAt := time.Now()

		logEntry := domain.WorkerLog{
			JobName:         jobNameVideoCategorySync,
			CountryCode:     country,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
			IntervalMinutes: int(categorySyncInterval.Minutes()),
		}

		if err != nil {
			log.Printf("❌ [Worker] ซิงก์หมวดหมู่ประเทศ %s ไม่สำเร็จ: %v\n", country, err)
			logEntry.Status = domain.WorkerLogStatusFailed
			logEntry.Message = err.Error()
		} else {
			logEntry.Status = domain.WorkerLogStatusSuccess
			logEntry.Message = "ซิงก์หมวดหมู่สำเร็จ"
		}

		if logErr := logRepo.CreateLog(logEntry); logErr != nil {
			log.Printf("⚠️ [Worker] บันทึก Log หมวดหมู่ของประเทศ %s ไม่สำเร็จ: %v\n", country, logErr)
		}
	}
}
