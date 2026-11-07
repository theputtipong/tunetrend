package http

import (
	"errors"
	"log"
	"time"

	"tunetrend-backend/internal/core/httputil"
	"tunetrend-backend/internal/core/ratelimit"
	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	usecase domain.SearchUsecase
}

func NewSearchHandler(usecase domain.SearchUsecase) *SearchHandler {
	return &SearchHandler{usecase: usecase}
}

func RegisterSearchRoutes(r fiber.Router, h *SearchHandler, authUsecase domain.AuthUsecase, trendsLimiter fiber.Handler) {
	r.Get("/trends/searched", trendsLimiter, h.GetSearchedTrends)
	r.Get("/trends/searched-queries", trendsLimiter, h.GetSearchedQueries)
	r.Get("/search", ratelimit.New(30, 60*time.Second).Named("search").Middleware(), h.Search)
	r.Get("/search/youtube", ratelimit.New(30, 60*time.Second).Named("search_youtube").Middleware(),
		RequireAuth(authUsecase), h.SearchYouTube)
	r.Post("/search/hit", ratelimit.New(60, 60*time.Second).Named("search_hit").Middleware(), h.RecordHit)
}

func searchErrorJSON(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"success": false, "error": message})
}

func searchCommonError(c *fiber.Ctx, err error) error {
	var ve *domain.ValidationError
	if errors.As(err, &ve) {
		return searchErrorJSON(c, fiber.StatusBadRequest, ve.Message)
	}
	log.Printf("❌ [Search] %v", err)
	return searchErrorJSON(c, fiber.StatusInternalServerError, "internal server error")
}

// Search godoc
// @Summary ค้นหาวิดีโอในชาร์ต TuneTrend
// @Description ค้นจาก songs (หมวดเพลง) หรือตารางของหมวดนั้น รวมกับ search_videos ไม่เรียก YouTube
// @Tags Search
// @Produce json
// @Param q query string true "คำค้น 2-100 ตัวอักษร"
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Param category query string false "รหัสหมวดหมู่ (Default: 10)"
// @Param track query string false "0 = ไม่นับเป็นการค้นในแท็บคำค้นยอดนิยม (ใช้ตอนเปิดจากลิงก์ในแท็บนั้นเอง)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Router /search [get]
func (h *SearchHandler) Search(c *fiber.Ctx) error {
	visitorIP := httputil.ClientIP(c)
	if c.Query("track") == "0" {
		visitorIP = ""
	}
	res, err := h.usecase.SearchLocal(visitorIP, c.Query("q"), c.Query("country", "TH"), c.Query("category"))
	if err != nil {
		return searchCommonError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

// SearchYouTube godoc
// @Summary ค้นหาเพิ่มบน YouTube (ต้อง login)
// @Description ใช้แคช 24 ชม. ก่อน แล้วตรวจสิทธิ์ต่อคนและโควตาระบบ error: user_limit | quota_exhausted | search_unavailable พร้อม retryAt
// @Tags Search
// @Produce json
// @Security BearerAuth
// @Param q query string true "คำค้น 2-100 ตัวอักษร"
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Param category query string false "รหัสหมวดหมู่ (Default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /search/youtube [get]
func (h *SearchHandler) SearchYouTube(c *fiber.Ctx) error {
	userID, _ := c.Locals(domain.LocalsUserID).(string)
	if userID == "" {
		return searchErrorJSON(c, fiber.StatusUnauthorized, "unauthorized")
	}

	res, err := h.usecase.SearchYouTube(userID, c.Query("q"), c.Query("country", "TH"), c.Query("category"))
	if err != nil {
		var qe *domain.QuotaError
		if errors.As(err, &qe) {
			status := fiber.StatusServiceUnavailable
			if errors.Is(err, domain.ErrUserLimit) {
				status = fiber.StatusTooManyRequests
			}
			return c.Status(status).JSON(fiber.Map{
				"success":   false,
				"error":     qe.Err.Error(),
				"retryAt":   qe.RetryAt.UTC().Format(time.RFC3339),
				"remaining": qe.Remaining,
			})
		}
		var ve *domain.ValidationError
		if errors.As(err, &ve) {
			return searchErrorJSON(c, fiber.StatusBadRequest, ve.Message)
		}
		log.Printf("❌ [Search] YouTube search failed: %v", err)
		return searchErrorJSON(c, fiber.StatusBadGateway, "youtube_error")
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

type recordHitRequest struct {
	VideoID string `json:"videoId"`
	Country string `json:"country"`
	Q       string `json:"q"`
}

// RecordHit godoc
// @Summary นับการกดเล่นจากหน้าผลค้นหา
// @Description นับได้ 1 ครั้ง/IP/วิดีโอ/ชม.
// @Tags Search
// @Accept json
// @Produce json
// @Param body body recordHitRequest true "videoId, country, q"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /search/hit [post]
func (h *SearchHandler) RecordHit(c *fiber.Ctx) error {
	var req recordHitRequest
	if err := c.BodyParser(&req); err != nil {
		return searchErrorJSON(c, fiber.StatusBadRequest, "invalid request body")
	}

	counted, err := h.usecase.RecordHit(httputil.ClientIP(c), req.VideoID, req.Country, req.Q)
	if err != nil {
		if errors.Is(err, domain.ErrUnknownVideo) {
			return searchErrorJSON(c, fiber.StatusNotFound, err.Error())
		}
		return searchCommonError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"counted": counted}})
}

// GetSearchedTrends godoc
// @Summary มาแรงจากการค้นหา
// @Description อันดับวิดีโอหมวดเพลงตามจำนวนการกดเล่นจากหน้าค้นหาใน 7 วัน (ไม่เกิน 30) พร้อมการเปลี่ยนอันดับเทียบ 7 วันก่อนหน้า
// @Tags Songs
// @Produce json
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /trends/searched [get]
func (h *SearchHandler) GetSearchedTrends(c *fiber.Ctx) error {
	trends, err := h.usecase.GetSearchedTrends(c.Query("country", "TH"))
	if err != nil {
		return searchCommonError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": trends})
}

// GetSearchedQueries godoc
// @Summary คำค้นยอดนิยม
// @Description กลุ่มคำค้นหมวดเพลงใน 7 วันล่าสุด เรียงตามจำนวนการค้น (นับ 1 ครั้ง/IP/คำ/ชม.) แสดงเฉพาะคำที่มีผู้ค้นต่างกันอย่างน้อย 2 คน ไม่เกิน 30 คำ คำละ 5 ผลจากการค้นล่าสุด
// @Tags Songs
// @Produce json
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /trends/searched-queries [get]
func (h *SearchHandler) GetSearchedQueries(c *fiber.Ctx) error {
	groups, err := h.usecase.GetSearchedQueries(c.Query("country", "TH"))
	if err != nil {
		return searchCommonError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": groups})
}
