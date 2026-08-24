package worker

import (
	"log"
	"time"

	"tunetrend-backend/internal/domain"
)

const (
	apiLogCleanupInterval = 24 * time.Hour
	apiLogRetention       = 14 * 24 * time.Hour
)

func StartApiLogCleanup(repo domain.ApiLogRepository) {
	log.Println("🚀 [ApiLog] เริ่มต้นระบบลบ log เก่าอัตโนมัติแล้ว...")

	cleanupApiLogs(repo)

	ticker := time.NewTicker(apiLogCleanupInterval)

	for range ticker.C {
		cleanupApiLogs(repo)
	}
}

func cleanupApiLogs(repo domain.ApiLogRepository) {
	if err := repo.DeleteOlderThan(time.Now().Add(-apiLogRetention)); err != nil {
		log.Printf("⚠️ [ApiLog] ลบ log เก่าไม่สำเร็จ: %v\n", err)
	}
}
