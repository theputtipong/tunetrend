package http

import (
	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type SongHandler struct {
	usecase domain.SongUsecase
}

func NewSongHandler(usecase domain.SongUsecase) *SongHandler {
	return &SongHandler{usecase: usecase}
}

// GetTrends godoc
// @Summary ดึงข้อมูลเพลงฮิต
// @Description ดึงข้อมูลเพลงเทรนด์ดิ้งตามรหัสประเทศ (เช่น TH, US) จากฐานข้อมูล
// @Tags Songs
// @Accept json
// @Produce json
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Success 200 {object} map[string]interface{} "คืนค่ารายการเพลงฮิต"
// @Failure 500 {object} map[string]interface{} "เกิดข้อผิดพลาดที่เซิร์ฟเวอร์"
// @Router /trends [get]
func (h *SongHandler) GetTrends(c *fiber.Ctx) error {
	countryCode := c.Query("country", "TH")

	songs, err := h.usecase.GetTrends(countryCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    songs,
	})
}

// GetNewReleases godoc
// @Summary ดึงเพลงฮิตมาใหม่ (ไม่เกิน 7 วัน)
// @Description คัดกรองเฉพาะเพลงเทรนด์ดิ้งที่เพิ่งปล่อยออกมาภายใน 7 วันที่ผ่านมา
// @Tags Songs
// @Accept json
// @Produce json
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Success 200 {object} map[string]interface{} "คืนค่ารายการเพลงใหม่"
// @Router /trends/new [get]
func (h *SongHandler) GetNewReleases(c *fiber.Ctx) error {
	countryCode := c.Query("country", "TH")
	songs, err := h.usecase.GetNewReleases(countryCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": songs})
}

// GetMVs godoc
// @Summary ดึงเฉพาะคลิป Official MV
// @Description คัดกรองเฉพาะเทรนด์เพลงที่ระบบ AI วิเคราะห์ว่าเป็น Official Music Video
// @Tags Songs
// @Accept json
// @Produce json
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Success 200 {object} map[string]interface{} "คืนค่ารายการ MV"
// @Router /trends/mv [get]
func (h *SongHandler) GetMVs(c *fiber.Ctx) error {
	countryCode := c.Query("country", "TH")
	songs, err := h.usecase.GetMVs(countryCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": songs})
}
