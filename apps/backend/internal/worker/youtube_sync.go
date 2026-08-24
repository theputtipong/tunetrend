package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

const (
	jobNameYouTubeSync = "youtube_sync"
	syncInterval       = 3 * time.Hour

	syncLockKey = 7_301_994

	logRetention = 90 * 24 * time.Hour
)

func StartYouTubeSync(uc domain.SongUsecase, logRepo domain.WorkerLogRepository) {
	log.Println("🚀 [Worker] เริ่มต้นระบบ Background Sync แล้ว...")

	syncData(uc, logRepo)

	ticker := time.NewTicker(syncInterval)

	for range ticker.C {
		log.Println("⏰ [Worker] ถึงเวลาอัปเดตข้อมูลรอบใหม่...")
		syncData(uc, logRepo)
	}
}

func syncData(uc domain.SongUsecase, logRepo domain.WorkerLogRepository) {
	acquired, err := logRepo.WithLock(syncLockKey, func() error {
		runSync(uc, logRepo)
		return nil
	})

	if err != nil {
		log.Printf("⚠️ [Worker] เกิดข้อผิดพลาดขณะขอ Lock สำหรับรอบ sync: %v\n", err)
		return
	}

	if !acquired {
		log.Println("⏭️ [Worker] มี instance อื่นกำลังรัน sync รอบนี้อยู่แล้ว ข้ามรอบนี้ไปก่อน")
	}
}

func runSync(uc domain.SongUsecase, logRepo domain.WorkerLogRepository) {
	countries := []string{"TH", "KR", "JP", "US", "GB"}

	for _, country := range countries {
		startedAt := time.Now()
		err := uc.SyncTrendingMusic(country)
		finishedAt := time.Now()

		logEntry := domain.WorkerLog{
			JobName:         jobNameYouTubeSync,
			CountryCode:     country,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
			IntervalMinutes: int(syncInterval.Minutes()),
		}

		if err != nil {
			log.Printf("❌ [Worker] ซิงก์ข้อมูลประเทศ %s ไม่สำเร็จ: %v\n", country, err)
			logEntry.Status = domain.WorkerLogStatusFailed
			logEntry.Message = err.Error()
		} else {
			logEntry.Status = domain.WorkerLogStatusSuccess
			logEntry.Message = "ซิงก์ข้อมูลสำเร็จ"
		}

		if logErr := logRepo.CreateLog(logEntry); logErr != nil {
			log.Printf("⚠️ [Worker] บันทึก Log ของประเทศ %s ไม่สำเร็จ: %v\n", country, logErr)
		}
	}

	if err := logRepo.DeleteOlderThan(time.Now().Add(-logRetention)); err != nil {
		log.Printf("⚠️ [Worker] ลบ Log เก่าไม่สำเร็จ: %v\n", err)
	}
}
