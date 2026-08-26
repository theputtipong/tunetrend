package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

func StartCategoryVideoSync(
	categoryID, jobName string,
	uc domain.CategoryVideoUsecase,
	categoryRepo domain.VideoCategoryRepository,
	logRepo domain.WorkerLogRepository,
	intervalMinutes int,
	countries []string,
) {
	log.Printf("🚀 [Worker] เริ่มต้นระบบ Sync %s แล้ว...\n", jobName)

	interval := time.Duration(intervalMinutes) * time.Minute

	runCategoryVideoSync(categoryID, jobName, uc, categoryRepo, logRepo, interval, countries)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		runCategoryVideoSync(categoryID, jobName, uc, categoryRepo, logRepo, interval, countries)
	}
}

func runCategoryVideoSync(
	categoryID, jobName string,
	uc domain.CategoryVideoUsecase,
	categoryRepo domain.VideoCategoryRepository,
	logRepo domain.WorkerLogRepository,
	interval time.Duration,
	countries []string,
) {
	for _, country := range countries {
		reasons, err := categoryRepo.GetDeactivatedReasons(country, []string{categoryID})
		if err == nil && reasons[categoryID] == domain.DeactivatedReasonAutoFetchFailure {
			log.Printf("⏭️ [Worker] ข้าม %s ประเทศ %s (ปิดอัตโนมัติ รอ resume worker รายสัปดาห์)\n", jobName, country)
			continue
		}

		startedAt := time.Now()
		err = uc.SyncVideos(country)
		finishedAt := time.Now()

		logEntry := domain.WorkerLog{
			JobName:         jobName,
			CountryCode:     country,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
			IntervalMinutes: int(interval.Minutes()),
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
