package http

import (
	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type VideoCategoryHandler struct {
	usecase domain.VideoCategoryUsecase
}

func NewVideoCategoryHandler(usecase domain.VideoCategoryUsecase) *VideoCategoryHandler {
	return &VideoCategoryHandler{usecase: usecase}
}

// GetCategories godoc
// @Summary ดึงรายการหมวดหมู่ที่ใช้กรอง /trends ได้
// @Description คืนค่าหมวดหมู่ที่ is_active=true (แอดมินเปิดใช้งานแล้ว) และมีตารางวิดีโอ sync ไว้จริง สำหรับใช้เป็นค่า category ใน /trends
// @Tags Categories
// @Accept json
// @Produce json
// @Param country query string false "รหัสประเทศ (Default: TH)"
// @Success 200 {object} map[string]interface{} "คืนค่ารายการหมวดหมู่"
// @Failure 500 {object} map[string]interface{} "เกิดข้อผิดพลาดที่เซิร์ฟเวอร์"
// @Router /categories [get]
func (h *VideoCategoryHandler) GetCategories(c *fiber.Ctx) error {
	countryCode := c.Query("country", "TH")

	categories, err := h.usecase.GetCategories(countryCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    categories,
	})
}
