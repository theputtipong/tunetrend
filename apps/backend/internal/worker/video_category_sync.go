package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

const jobNameVideoCategorySync = "video_category_sync"

func StartVideoCategorySync(uc domain.VideoCategoryUsecase, logRepo domain.WorkerLogRepository, intervalMinutes int, countries []string) {
	log.Println("🚀 [Worker] เริ่มต้นระบบ Sync หมวดหมู่ YouTube แล้ว...")

	interval := time.Duration(intervalMinutes) * time.Minute

	syncCategories(uc, logRepo, interval, countries)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		syncCategories(uc, logRepo, interval, countries)
	}
}

func syncCategories(uc domain.VideoCategoryUsecase, logRepo domain.WorkerLogRepository, interval time.Duration, countries []string) {
	for _, country := range countries {
		startedAt := time.Now()
		err := uc.SyncCategories(country)
		finishedAt := time.Now()

		logEntry := domain.WorkerLog{
			JobName:         jobNameVideoCategorySync,
			CountryCode:     country,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
			IntervalMinutes: int(interval.Minutes()),
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
