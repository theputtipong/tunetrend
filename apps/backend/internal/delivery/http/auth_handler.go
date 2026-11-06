package http

import (
	"errors"
	"log"

	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type GoogleLoginRequest struct {
	IDToken string `json:"idToken"`
}

type AuthHandler struct {
	usecase domain.AuthUsecase
}

func NewAuthHandler(usecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func RegisterAuthRoutes(r fiber.Router, authUsecase domain.AuthUsecase, playlistUsecase domain.PlaylistUsecase, loginLimiter fiber.Handler) {
	authHandler := NewAuthHandler(authUsecase)
	playlistHandler := NewPlaylistHandler(playlistUsecase)
	requireAuth := RequireAuth(authUsecase)

	r.Post("/auth/google", loginLimiter, authHandler.GoogleLogin)

	me := r.Group("/me", requireAuth)
	me.Get("/", authHandler.GetMe)
	me.Delete("/", authHandler.DeleteMe)
	me.Get("/playlists", playlistHandler.List)
	me.Post("/playlists", playlistHandler.Create)
	me.Get("/playlists/:id", playlistHandler.Get)
	me.Patch("/playlists/:id", playlistHandler.Rename)
	me.Delete("/playlists/:id", playlistHandler.Delete)
	me.Post("/playlists/:id/items", playlistHandler.AddItem)
	me.Delete("/playlists/:id/items/:videoId", playlistHandler.RemoveItem)
	me.Put("/playlists/:id/order", playlistHandler.Reorder)
}

// GoogleLogin godoc
// @Summary เข้าสู่ระบบด้วย Google
// @Description ตรวจ Google ID token แล้วออก session JWT (HS256 อายุ 30 วัน) สร้าง user ใหม่ถ้ายังไม่เคย login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body GoogleLoginRequest true "Google ID token"
// @Success 200 {object} map[string]interface{} "token และข้อมูล user"
// @Failure 401 {object} map[string]interface{} "token ไม่ถูกต้องหรือหมดอายุ"
// @Failure 429 {object} map[string]interface{} "ส่งคำขอถี่เกินไป"
// @Failure 503 {object} map[string]interface{} "ตรวจ token กับ Google ไม่ได้ชั่วคราว (auth_unavailable)"
// @Router /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	var req GoogleLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
	}
	session, err := h.usecase.LoginWithGoogle(c.UserContext(), req.IDToken)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return unauthorized(c)
		}
		if errors.Is(err, domain.ErrAuthUnavailable) {
			log.Printf("⚠️ [Auth] Google token verification unavailable: %v", err)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"success": false, "error": "auth_unavailable"})
		}
		return authInternalError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": session})
}

// GetMe godoc
// @Summary ข้อมูลผู้ใช้ที่ login อยู่
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "ข้อมูล user"
// @Failure 401 {object} map[string]interface{} "ยังไม่ login"
// @Router /me [get]
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	user, err := h.usecase.GetUser(currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return unauthorized(c)
		}
		return authInternalError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": user})
}

// DeleteMe godoc
// @Summary ลบบัญชีพร้อม playlist ทั้งหมด
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "ลบบัญชีแล้ว"
// @Failure 401 {object} map[string]interface{} "ยังไม่ login"
// @Router /me [delete]
func (h *AuthHandler) DeleteMe(c *fiber.Ctx) error {
	if err := h.usecase.DeleteUser(currentUserID(c)); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return unauthorized(c)
		}
		return authInternalError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"deleted": true}})
}
