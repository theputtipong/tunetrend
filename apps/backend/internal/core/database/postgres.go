package database

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/repository/postgres"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var validTableName = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)

var reservedTableNames = map[string]bool{
	"songs":                  true,
	"worker_logs":            true,
	"contact_messages":       true,
	"api_logs":               true,
	"video_categories":       true,
	"category_video_configs": true,
	"worker_settings":        true,
}

func Connect() (*gorm.DB, []domain.CategoryVideoConfig, domain.WorkerSettings) {
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

	db, err := gorm.Open(gormpostgres.New(gormpostgres.Config{
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

	err = db.AutoMigrate(
		&domain.Song{}, &domain.WorkerLog{}, &domain.ContactMessage{}, &domain.ApiLog{}, &domain.VideoCategory{},
		&domain.CategoryVideoConfig{}, &domain.WorkerSettings{},
	)
	if err != nil {
		log.Fatalf("❌ ล้มเหลวในการ Auto Migrate: %v", err)
	}

	if err := seedCategoryVideoConfigs(db); err != nil {
		log.Fatalf("❌ ล้มเหลวในการ seed category_video_configs: %v", err)
	}
	if err := seedWorkerSettings(db); err != nil {
		log.Fatalf("❌ ล้มเหลวในการ seed worker_settings: %v", err)
	}

	configs, err := postgres.NewCategoryVideoConfigRepository(db).GetAll()
	if err != nil {
		log.Fatalf("❌ ไม่สามารถอ่าน category_video_configs ได้: %v", err)
	}

	for _, cfg := range configs {
		if !validTableName.MatchString(cfg.TableName) {
			log.Fatalf("❌ table_name \"%s\" ของหมวดหมู่ %s ไม่ถูกต้องตามรูปแบบที่อนุญาต", cfg.TableName, cfg.CategoryID)
		}
		if reservedTableNames[cfg.TableName] {
			log.Fatalf("❌ table_name \"%s\" ของหมวดหมู่ %s ชนกับตารางระบบ ห้ามใช้ชื่อนี้", cfg.TableName, cfg.CategoryID)
		}
		if err := db.Table(cfg.TableName).AutoMigrate(&domain.CategoryVideo{}); err != nil {
			log.Fatalf("❌ ล้มเหลวในการ Auto Migrate ตาราง %s: %v", cfg.TableName, err)
		}
	}

	settings, err := postgres.NewWorkerSettingsRepository(db).Get()
	if err != nil {
		log.Fatalf("❌ ไม่สามารถอ่าน worker_settings ได้: %v", err)
	}

	log.Println("✅ เชื่อมต่อ Database และ Auto Migrate สำเร็จ!")
	return db, configs, settings
}
