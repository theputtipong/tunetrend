package main

import (
	"log"
	"os"
	"strings"

	"time"

	"tunetrend-backend/internal/core/apilog"
	"tunetrend-backend/internal/core/auth"
	"tunetrend-backend/internal/core/database"
	"tunetrend-backend/internal/core/mail"
	"tunetrend-backend/internal/core/ratelimit"
	"tunetrend-backend/internal/core/searchcache"
	"tunetrend-backend/internal/core/upstash"
	"tunetrend-backend/internal/core/ytquota"
	"tunetrend-backend/internal/delivery/http"
	"tunetrend-backend/internal/domain"
	"tunetrend-backend/internal/repository/postgres"
	"tunetrend-backend/internal/repository/youtube"
	"tunetrend-backend/internal/usecase"

	"tunetrend-backend/internal/worker"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
)

// @title TuneTrend API
// @version 1.0
// @description Backend Service สำหรับดึงข้อมูลเทรนด์เพลงฮิตจาก YouTube
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()

	db, categoryConfigs, workerSettings := database.Connect()
	countries := strings.Split(workerSettings.Countries, ",")

	songRepo := postgres.NewSongRepository(db)
	workerLogRepo := postgres.NewWorkerLogRepository(db)
	upstashClient := upstash.NewFromEnv()
	mailer := mail.NewResendMailer()
	quotaGuard := ytquota.New(upstashClient, mailer, workerLogRepo, ytquota.OptionsFromEnv())
	ytRepo := ytquota.WrapYouTube(youtube.NewYouTubeRepository(), quotaGuard)
	categoryRepo := postgres.NewVideoCategoryRepository(db)

	categoryRepos := make(map[string]domain.CategoryVideoRepository, len(categoryConfigs))
	categoryVideoUsecases := make(map[string]domain.CategoryVideoUsecase, len(categoryConfigs))
	for _, cfg := range categoryConfigs {
		catVideoRepo := postgres.NewCategoryVideoRepository(db, cfg.TableName)
		categoryRepos[cfg.CategoryID] = catVideoRepo
		catVideoUsecase := usecase.NewCategoryVideoUsecase(catVideoRepo, ytRepo, cfg.CategoryID, cfg.Label)
		categoryVideoUsecases[cfg.CategoryID] = catVideoUsecase
		go worker.StartCategoryVideoSync(
			cfg.CategoryID, cfg.TableName, catVideoUsecase, categoryRepo, workerLogRepo,
			workerSettings.CategoryVideoSyncIntervalMinutes, countries,
		)
	}

	songUsecase := usecase.NewSongUsecase(songRepo, ytRepo, categoryRepos)
	songHandler := http.NewSongHandler(songUsecase)

	contactRepo := postgres.NewContactRepository(db)
	contactUsecase := usecase.NewContactUsecase(contactRepo, mailer)
	contactHandler := http.NewContactHandler(contactUsecase)

	sessionTokens := auth.NewSessionTokens(os.Getenv("JWT_SECRET"), auth.SessionTTL)
	googleVerifier := auth.NewGoogleVerifier(os.Getenv("GOOGLE_CLIENT_IDS"))
	if !sessionTokens.Configured() || !googleVerifier.Configured() {
		log.Println("⚠️  JWT_SECRET (>=32 bytes) หรือ GOOGLE_CLIENT_IDS ยังไม่ได้ตั้ง — endpoint auth จะตอบ 500")
	}
	authUsecase := usecase.NewAuthUsecase(googleVerifier, sessionTokens, postgres.NewUserRepository(db))
	playlistUsecase := usecase.NewPlaylistUsecase(postgres.NewPlaylistRepository(db))

	searchUsecase := usecase.NewSearchUsecase(usecase.SearchDeps{
		Repo: postgres.NewSearchRepository(db, categoryConfigs), YouTube: ytRepo, Guard: quotaGuard,
		Cache: searchcache.New(upstashClient), Deduper: searchcache.NewHitDeduper(upstashClient),
		CategoryRepo: categoryRepo, CategoryConfigs: categoryConfigs, Countries: countries,
		VisitorKey: []byte(os.Getenv("JWT_SECRET")),
	})
	searchHandler := http.NewSearchHandler(searchUsecase)

	apiLogRepo := postgres.NewApiLogRepository(db)

	categoryUsecase := usecase.NewVideoCategoryUsecase(
		categoryRepo, ytRepo, workerLogRepo, categoryConfigs, workerSettings.CategoryFetchFailureGraceMinutes,
	)
	categoryHandler := http.NewVideoCategoryHandler(categoryUsecase)

	discoverUsecase := usecase.NewDiscoverUsecase(categoryRepos, categoryConfigs)
	discoverHandler := http.NewDiscoverHandler(discoverUsecase)

	go worker.StartYouTubeSync(songUsecase, workerLogRepo, workerSettings.MusicSyncIntervalMinutes, countries)
	go worker.StartApiLogCleanup(apiLogRepo)
	go worker.StartSearchVideoRefresh(searchUsecase, workerLogRepo)
	go worker.StartVideoCategorySync(categoryUsecase, workerLogRepo, workerSettings.CategorySyncIntervalMinutes, countries)
	go worker.StartCategoryVideoResumeWorker(
		categoryRepo, categoryVideoUsecases, workerLogRepo, categoryConfigs,
		workerSettings.CategoryResumeIntervalMinutes, countries,
	)

	app := fiber.New(fiber.Config{BodyLimit: 1 * 1024 * 1024})

	corsAllowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsAllowOrigins == "" {
		corsAllowOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsAllowOrigins,
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
		AllowHeaders: "Authorization,Content-Type",
	}))
	app.Use(requestid.New())
	app.Use(apilog.New(apiLogRepo))

	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"status":  "ERROR",
				"message": "Failed to get database instance",
				"error":   err.Error(),
			})
		}

		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success": false,
				"status":  "DOWN",
				"message": "Database is unreachable 💀",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"status":  "OK",
			"message": "API and Database are up and running! 🚀",
		})
	})

	trendsRateLimit := ratelimit.New(60, 60*time.Second).Middleware()
	app.Get("/trends", trendsRateLimit, songHandler.GetTrends)
	app.Get("/trends/new", trendsRateLimit, songHandler.GetNewReleases)
	app.Get("/trends/mv", trendsRateLimit, songHandler.GetMVs)
	app.Get("/categories", trendsRateLimit, categoryHandler.GetCategories)
	app.Get("/discover", trendsRateLimit, discoverHandler.GetDiscoverItems)

	contactRateLimit := ratelimit.New(5, 10*time.Minute).FailClosed().Middleware()
	app.Post("/contact", contactRateLimit, contactHandler.SubmitContact)

	http.RegisterSearchRoutes(app, searchHandler, authUsecase, trendsRateLimit)

	http.RegisterAuthRoutes(app, authUsecase, playlistUsecase, ratelimit.New(10, time.Minute).Named("auth").Middleware())

	SetupSwaggerRoutes(app)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Fatal(app.Listen(":" + appPort))
}
