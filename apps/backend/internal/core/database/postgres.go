package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"tunetrend-backend/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	var dsn string

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		dsn = dbURL
		log.Println("🔌 Connecting to database using DATABASE_URL...")
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_TIMEZONE"),
		)
		log.Println("🔌 Connecting to database using DB_HOST, DB_USER, etc...")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatalf("❌ ล้มเหลวในการเชื่อมต่อ Database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ ไม่สามารถดึง sql.DB ออกมาตั้งค่าได้: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&domain.Song{}, &domain.WorkerLog{}, &domain.ContactMessage{}, &domain.ApiLog{}, &domain.VideoCategory{})
	if err != nil {
		log.Fatalf("❌ ล้มเหลวในการ Auto Migrate: %v", err)
	}

	for _, cfg := range domain.CategoryVideoConfigs {
		if err := db.Table(cfg.TableName).AutoMigrate(&domain.CategoryVideo{}); err != nil {
			log.Fatalf("❌ ล้มเหลวในการ Auto Migrate ตาราง %s: %v", cfg.TableName, err)
		}
	}

	log.Println("✅ เชื่อมต่อ Database และ Auto Migrate สำเร็จ!")
	return db
}
