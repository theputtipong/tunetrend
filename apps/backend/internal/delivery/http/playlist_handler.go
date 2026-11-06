package http

import (
	"errors"
	"log"

	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type PlaylistNameRequest struct {
	Name string `json:"name"`
}

type PlaylistItemRequest struct {
	VideoID      string `json:"videoId"`
	Title        string `json:"title"`
	ChannelTitle string `json:"channelTitle"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type PlaylistOrderRequest struct {
	VideoIDs []string `json:"videoIds"`
}

type PlaylistHandler struct {
	usecase domain.PlaylistUsecase
}

func NewPlaylistHandler(usecase domain.PlaylistUsecase) *PlaylistHandler {
	return &PlaylistHandler{usecase: usecase}
}

// List godoc
// @Summary รายการ playlist ของฉัน
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "รายการ playlist พร้อม itemCount"
// @Failure 401 {object} map[string]interface{} "ยังไม่ login"
// @Router /me/playlists [get]
func (h *PlaylistHandler) List(c *fiber.Ctx) error {
	playlists, err := h.usecase.List(currentUserID(c))
	if err != nil {
		return playlistError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": playlists})
}

// Create godoc
// @Summary สร้าง playlist
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PlaylistNameRequest true "ชื่อ playlist (1-60 ตัวอักษร)"
// @Success 201 {object} map[string]interface{} "playlist ที่สร้าง"
// @Failure 422 {object} map[string]interface{} "ชื่อไม่ถูกต้องหรือเกิน 50 playlist"
// @Router /me/playlists [post]
func (h *PlaylistHandler) Create(c *fiber.Ctx) error {
	var req PlaylistNameRequest
	if err := c.BodyParser(&req); err != nil {
		return badBody(c)
	}
	p, err := h.usecase.Create(currentUserID(c), req.Name)
	if err != nil {
		return playlistError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": p})
}

// Get godoc
// @Summary playlist พร้อมรายการเพลงเรียงตามลำดับ
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path string true "playlist id"
// @Success 200 {object} map[string]interface{} "playlist และ items"
// @Failure 404 {object} map[string]interface{} "ไม่พบ playlist"
// @Router /me/playlists/{id} [get]
func (h *PlaylistHandler) Get(c *fiber.Ctx) error {
	p, err := h.usecase.Get(currentUserID(c), c.Params("id"))
	if err != nil {
		return playlistError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": p})
}

// Rename godoc
// @Summary แก้ชื่อ playlist
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "playlist id"
// @Param request body PlaylistNameRequest true "ชื่อใหม่"
// @Success 200 {object} map[string]interface{} "playlist ที่แก้แล้ว"
// @Failure 404 {object} map[string]interface{} "ไม่พบ playlist"
// @Failure 422 {object} map[string]interface{} "ชื่อไม่ถูกต้อง"
// @Router /me/playlists/{id} [patch]
func (h *PlaylistHandler) Rename(c *fiber.Ctx) error {
	var req PlaylistNameRequest
	if err := c.BodyParser(&req); err != nil {
		return badBody(c)
	}
	p, err := h.usecase.Rename(currentUserID(c), c.Params("id"), req.Name)
	if err != nil {
		return playlistError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": p})
}

// Delete godoc
// @Summary ลบ playlist
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path string true "playlist id"
// @Success 200 {object} map[string]interface{} "ลบแล้ว"
// @Failure 404 {object} map[string]interface{} "ไม่พบ playlist"
// @Router /me/playlists/{id} [delete]
func (h *PlaylistHandler) Delete(c *fiber.Ctx) error {
	if err := h.usecase.Delete(currentUserID(c), c.Params("id")); err != nil {
		return playlistError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"deleted": true}})
}

// AddItem godoc
// @Summary เพิ่มเพลงลง playlist
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "playlist id"
// @Param request body PlaylistItemRequest true "ข้อมูลเพลง (snapshot)"
// @Success 201 {object} map[string]interface{} "item ที่เพิ่ม"
// @Failure 404 {object} map[string]interface{} "ไม่พบ playlist"
// @Failure 409 {object} map[string]interface{} "มีเพลงนี้ใน playlist แล้ว"
// @Failure 422 {object} map[string]interface{} "ข้อมูลไม่ถูกต้องหรือเกิน 200 เพลง"
// @Router /me/playlists/{id}/items [post]
func (h *PlaylistHandler) AddItem(c *fiber.Ctx) error {
	var req PlaylistItemRequest
	if err := c.BodyParser(&req); err != nil {
		return badBody(c)
	}
	item, err := h.usecase.AddItem(currentUserID(c), c.Params("id"), domain.PlaylistItemInput{
		VideoID:      req.VideoID,
		Title:        req.Title,
		ChannelTitle: req.ChannelTitle,
		ThumbnailURL: req.ThumbnailURL,
	})
	if err != nil {
		return playlistError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": item})
}

// RemoveItem godoc
// @Summary ลบเพลงออกจาก playlist
// @Tags Playlists
// @Produce json
// @Security BearerAuth
// @Param id path string true "playlist id"
// @Param videoId path string true "YouTube video id"
// @Success 200 {object} map[string]interface{} "ลบแล้ว"
// @Failure 404 {object} map[string]interface{} "ไม่พบ playlist หรือเพลง"
// @Router /me/playlists/{id}/items/{videoId} [delete]
func (h *PlaylistHandler) RemoveItem(c *fiber.Ctx) error {
	if err := h.usecase.RemoveItem(currentUserID(c), c.Params("id"), c.Params("videoId")); err != nil {
		return playlistError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": fiber.Map{"deleted": true}})
}

// Reorder godoc
// @Summary จัดลำดับเพลงใน playlist
// @Description videoIds ต้องมีเพลงทุกเพลงใน playlist ครบและไม่ซ้ำ
// @Tags Playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "playlist id"
// @Param request body PlaylistOrderRequest true "ลำดับใหม่"
// @Success 200 {object} map[string]interface{} "playlist พร้อมลำดับใหม่"
// @Failure 404 {object} map[string]interface{} "ไม่พบ playlist"
// @Failure 422 {object} map[string]interface{} "videoIds ไม่ตรงกับเพลงใน playlist"
// @Router /me/playlists/{id}/order [put]
func (h *PlaylistHandler) Reorder(c *fiber.Ctx) error {
	var req PlaylistOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return badBody(c)
	}
	p, err := h.usecase.Reorder(currentUserID(c), c.Params("id"), req.VideoIDs)
	if err != nil {
		return playlistError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": p})
}

func badBody(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "invalid request body"})
}

func playlistError(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "internal server error"
	var verr *domain.ValidationError
	switch {
	case errors.As(err, &verr):
		status, message = fiber.StatusUnprocessableEntity, verr.Message
	case errors.Is(err, domain.ErrPlaylistNotFound), errors.Is(err, domain.ErrPlaylistItemNotFound):
		status, message = fiber.StatusNotFound, err.Error()
	case errors.Is(err, domain.ErrPlaylistLimitReached),
		errors.Is(err, domain.ErrPlaylistItemLimitReached),
		errors.Is(err, domain.ErrPlaylistInvalidOrder):
		status, message = fiber.StatusUnprocessableEntity, err.Error()
	case errors.Is(err, domain.ErrVideoAlreadyInPlaylist):
		status, message = fiber.StatusConflict, err.Error()
	case errors.Is(err, domain.ErrUserNotFound):
		return unauthorized(c)
	default:
		log.Printf("❌ [Playlist] internal error: %v", err)
	}
	return c.Status(status).JSON(fiber.Map{"success": false, "error": message})
}
