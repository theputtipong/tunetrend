package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

const categoryVideoSyncInterval = 3 * time.Hour

var categoryVideoSyncCountries = []string{"TH", "KR", "JP", "US", "GB"}

func StartCategoryVideoSync(jobName string, uc domain.CategoryVideoUsecase, logRepo domain.WorkerLogRepository) {
	log.Printf("🚀 [Worker] เริ่มต้นระบบ Sync %s แล้ว...\n", jobName)

	runCategoryVideoSync(jobName, uc, logRepo)

	ticker := time.NewTicker(categoryVideoSyncInterval)
	for range ticker.C {
		runCategoryVideoSync(jobName, uc, logRepo)
	}
}

func runCategoryVideoSync(jobName string, uc domain.CategoryVideoUsecase, logRepo domain.WorkerLogRepository) {
	for _, country := range categoryVideoSyncCountries {
		startedAt := time.Now()
		err := uc.SyncVideos(country)
		finishedAt := time.Now()

		logEntry := domain.WorkerLog{
			JobName:         jobName,
			CountryCode:     country,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
			IntervalMinutes: int(categoryVideoSyncInterval.Minutes()),
		}

		if err != nil {
			log.Printf("❌ [Worker] ซิงก์ %s ประเทศ %s ไม่สำเร็จ: %v\n", jobName, country, err)
			logEntry.Status = domain.WorkerLogStatusFailed
			logEntry.Message = err.Error()
		} else {
			logEntry.Status = domain.WorkerLogStatusSuccess
			logEntry.Message = "ซิงก์สำเร็จ"
		}

		if logErr := logRepo.CreateLog(logEntry); logErr != nil {
			log.Printf("⚠️ [Worker] บันทึก Log ของ %s ประเทศ %s ไม่สำเร็จ: %v\n", jobName, country, logErr)
		}
	}
}
