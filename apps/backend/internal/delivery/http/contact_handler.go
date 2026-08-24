package http

import (
	"errors"
	"log"

	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type ContactRequest struct {
	Name         string `json:"name"`
	Message      string `json:"message"`
	ContactEmail string `json:"contactEmail"`
	ContactPhone string `json:"contactPhone"`
	Website      string `json:"website"`
}

type ContactHandler struct {
	usecase domain.ContactUsecase
}

func NewContactHandler(usecase domain.ContactUsecase) *ContactHandler {
	return &ContactHandler{usecase: usecase}
}

// SubmitContact godoc
// @Summary ส่งข้อความติดต่อถึงผู้พัฒนา
// @Description รับข้อความจากผู้ใช้พร้อมช่องทางติดต่อกลับ (อีเมลหรือเบอร์โทรไทย) บันทึกลงฐานข้อมูลและส่งอีเมลแจ้งเตือนผู้พัฒนา
// @Tags Contact
// @Accept json
// @Produce json
// @Param request body ContactRequest true "ข้อมูลฟอร์มติดต่อ"
// @Success 200 {object} map[string]interface{} "ส่งข้อความสำเร็จ"
// @Failure 400 {object} map[string]interface{} "ข้อมูลไม่ถูกต้อง เช่น เบอร์โทรไม่ใช่รูปแบบไทย"
// @Failure 429 {object} map[string]interface{} "ส่งคำขอถี่เกินไป"
// @Failure 500 {object} map[string]interface{} "เกิดข้อผิดพลาดที่เซิร์ฟเวอร์"
// @Router /contact [post]
func (h *ContactHandler) SubmitContact(c *fiber.Ctx) error {
	var req ContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}

	err := h.usecase.SubmitContactMessage(domain.ContactSubmission{
		Name:         req.Name,
		Message:      req.Message,
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
		Honeypot:     req.Website,
	})

	if err != nil {
		var verr *domain.ValidationError
		if errors.As(err, &verr) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": verr.Message})
		}

		log.Printf("❌ [Contact] internal error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "internal server error"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"message": "ส่งข้อความเรียบร้อยแล้ว ขอบคุณที่ติดต่อเรา"},
	})
}
