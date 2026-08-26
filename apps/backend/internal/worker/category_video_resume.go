package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

func StartCategoryVideoResumeWorker(
	categoryRepo domain.VideoCategoryRepository,
	usecases map[string]domain.CategoryVideoUsecase,
	logRepo domain.WorkerLogRepository,
	configs []domain.CategoryVideoConfig,
	intervalMinutes int,
	countries []string,
) {
	log.Println("🚀 [Worker] เริ่มต้นระบบ Resume หมวดหมู่ที่ถูกปิดอัตโนมัติแล้ว...")

	interval := time.Duration(intervalMinutes) * time.Minute

	runCategoryVideoResume(categoryRepo, usecases, logRepo, configs, interval, countries)

	ticker := time.NewTicker(interval)
	for range ticker.C {
		runCategoryVideoResume(categoryRepo, usecases, logRepo, configs, interval, countries)
	}
}

func runCategoryVideoResume(
	categoryRepo domain.VideoCategoryRepository,
	usecases map[string]domain.CategoryVideoUsecase,
	logRepo domain.WorkerLogRepository,
	configs []domain.CategoryVideoConfig,
	interval time.Duration,
	countries []string,
) {
	for _, cfg := range configs {
		uc, ok := usecases[cfg.CategoryID]
		if !ok {
			continue
		}

		for _, country := range countries {
			resumeIfPaused(cfg, country, uc, categoryRepo, logRepo, interval)
		}
	}
}

func resumeIfPaused(
	cfg domain.CategoryVideoConfig,
	country string,
	uc domain.CategoryVideoUsecase,
	categoryRepo domain.VideoCategoryRepository,
	logRepo domain.WorkerLogRepository,
	interval time.Duration,
) {
	reasons, err := categoryRepo.GetDeactivatedReasons(country, []string{cfg.CategoryID})
	if err != nil || reasons[cfg.CategoryID] != domain.DeactivatedReasonAutoFetchFailure {
		return
	}

	log.Printf("🔁 [Worker] Resume %s ประเทศ %s (รายสัปดาห์)\n", cfg.TableName, country)

	startedAt := time.Now()
	err = uc.SyncVideos(country)
	finishedAt := time.Now()

	logEntry := domain.WorkerLog{
		JobName:         cfg.TableName,
		CountryCode:     country,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
		IntervalMinutes: int(interval.Minutes()),
	}

	if err != nil {
		log.Printf("❌ [Worker] Resume %s ประเทศ %s ไม่สำเร็จ: %v\n", cfg.TableName, country, err)
		logEntry.Status = domain.WorkerLogStatusFailed
		logEntry.Message = err.Error()
	} else {
		log.Printf("✅ [Worker] Resume %s ประเทศ %s สำเร็จ\n", cfg.TableName, country)
		logEntry.Status = domain.WorkerLogStatusSuccess
		logEntry.Message = "resume สำเร็จ"
	}

	if logErr := logRepo.CreateLog(logEntry); logErr != nil {
		log.Printf("⚠️ [Worker] บันทึก Log resume ของ %s ประเทศ %s ไม่สำเร็จ: %v\n", cfg.TableName, country, logErr)
	}
}
