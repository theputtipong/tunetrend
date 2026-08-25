package main

import (
	"log"
	"os"

	"time"

	"tunetrend-backend/internal/core/apilog"
	"tunetrend-backend/internal/core/database"
	"tunetrend-backend/internal/core/mail"
	"tunetrend-backend/internal/core/ratelimit"
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

	db := database.Connect()

	songRepo := postgres.NewSongRepository(db)
	ytRepo := youtube.NewYouTubeRepository()
	workerLogRepo := postgres.NewWorkerLogRepository(db)

	categoryRepos := make(map[string]domain.CategoryVideoRepository, len(domain.CategoryVideoConfigs))
	for _, cfg := range domain.CategoryVideoConfigs {
		catVideoRepo := postgres.NewCategoryVideoRepository(db, cfg.TableName)
		categoryRepos[cfg.CategoryID] = catVideoRepo
		catVideoUsecase := usecase.NewCategoryVideoUsecase(catVideoRepo, ytRepo, cfg.CategoryID, cfg.Label)
		go worker.StartCategoryVideoSync(cfg.TableName, catVideoUsecase, workerLogRepo)
	}

	songUsecase := usecase.NewSongUsecase(songRepo, ytRepo, categoryRepos)
	songHandler := http.NewSongHandler(songUsecase)

	contactRepo := postgres.NewContactRepository(db)
	mailer := mail.NewResendMailer()
	contactUsecase := usecase.NewContactUsecase(contactRepo, mailer)
	contactHandler := http.NewContactHandler(contactUsecase)

	apiLogRepo := postgres.NewApiLogRepository(db)

	categoryRepo := postgres.NewVideoCategoryRepository(db)
	categoryUsecase := usecase.NewVideoCategoryUsecase(categoryRepo, ytRepo)
	categoryHandler := http.NewVideoCategoryHandler(categoryUsecase)

	go worker.StartYouTubeSync(songUsecase, workerLogRepo)
	go worker.StartApiLogCleanup(apiLogRepo)
	go worker.StartVideoCategorySync(categoryUsecase, workerLogRepo)

	app := fiber.New(fiber.Config{BodyLimit: 1 * 1024 * 1024})

	corsAllowOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsAllowOrigins == "" {
		corsAllowOrigins = "*"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsAllowOrigins,
		AllowMethods: "GET,POST",
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

	contactRateLimit := ratelimit.New(5, 10*time.Minute).FailClosed().Middleware()
	app.Post("/contact", contactRateLimit, contactHandler.SubmitContact)

	SetupSwaggerRoutes(app)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	log.Fatal(app.Listen(":" + appPort))
}
