package http

import (
	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type DiscoverHandler struct {
	usecase domain.DiscoverUsecase
}

func NewDiscoverHandler(usecase domain.DiscoverUsecase) *DiscoverHandler {
	return &DiscoverHandler{usecase: usecase}
}

// GetDiscoverItems godoc
// @Summary ดึงรายการแนะนำหมวดหมู่อื่นๆ สำหรับ carousel
// @Description คืนค่าวิดีโอฮิตสุด 1-2 รายการต่อหมวดหมู่ รวมทุกประเทศ ไม่ซ้ำกัน
// @Tags Discover
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "คืนค่ารายการแนะนำ"
// @Failure 500 {object} map[string]interface{} "เกิดข้อผิดพลาดที่เซิร์ฟเวอร์"
// @Router /discover [get]
func (h *DiscoverHandler) GetDiscoverItems(c *fiber.Ctx) error {
	items, err := h.usecase.GetDiscoverItems()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    items,
	})
}
