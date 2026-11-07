package worker

import (
	"fmt"
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

const (
	searchVideoRefreshLockKey  = 7_301_995
	searchVideoRefreshInterval = 6 * time.Hour
)

func StartSearchVideoRefresh(uc domain.SearchUsecase, logRepo domain.WorkerLogRepository) {
	log.Println("🚀 [Worker] เริ่มต้นระบบอัปเดตยอดวิว search_videos แล้ว...")

	runSearchVideoRefresh(uc, logRepo)

	ticker := time.NewTicker(searchVideoRefreshInterval)
	for range ticker.C {
		runSearchVideoRefresh(uc, logRepo)
	}
}

func runSearchVideoRefresh(uc domain.SearchUsecase, logRepo domain.WorkerLogRepository) {
	acquired, err := logRepo.WithLock(searchVideoRefreshLockKey, func() error {
		startedAt := time.Now()
		n, runErr := uc.RefreshSearchVideoStats()
		finishedAt := time.Now()

		entry := domain.WorkerLog{
			JobName:         domain.SearchVideoRefreshJobName,
			StartedAt:       startedAt,
			FinishedAt:      finishedAt,
			DurationMs:      finishedAt.Sub(startedAt).Milliseconds(),
			IntervalMinutes: int(searchVideoRefreshInterval.Minutes()),
			Status:          domain.WorkerLogStatusSuccess,
			Message:         fmt.Sprintf("อัปเดตยอดวิว search_videos สำเร็จ %d รายการ", n),
		}
		if runErr != nil {
			log.Printf("❌ [Worker] อัปเดต search_videos ไม่สำเร็จ: %v\n", runErr)
			entry.Status = domain.WorkerLogStatusFailed
			entry.Message = runErr.Error()
		}
		if pruned, pruneErr := uc.PruneSearchQueryHits(); pruneErr != nil {
			log.Printf("⚠️ [Worker] ลบ search_query_hits เก่าไม่สำเร็จ: %v\n", pruneErr)
		} else if pruned > 0 {
			log.Printf("🧹 [Worker] ลบ search_query_hits ที่เก่ากว่า 7 วัน %d แถว\n", pruned)
		}
		if logErr := logRepo.CreateLog(entry); logErr != nil {
			log.Printf("⚠️ [Worker] บันทึก Log ของ search_video_refresh ไม่สำเร็จ: %v\n", logErr)
		}
		return nil
	})
	if err != nil {
		log.Printf("⚠️ [Worker] ขอ Lock สำหรับ search_video_refresh ไม่สำเร็จ: %v\n", err)
		return
	}
	if !acquired {
		log.Println("⏭️ [Worker] มี instance อื่นกำลังอัปเดต search_videos อยู่ ข้ามรอบนี้")
	}
}
