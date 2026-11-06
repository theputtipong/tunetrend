package http

import (
	"errors"
	"log"
	"strings"

	"tunetrend-backend/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func RequireAuth(authUsecase domain.AuthUsecase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return unauthorized(c)
		}
		userID, err := authUsecase.Authenticate(token)
		if err != nil {
			if errors.Is(err, domain.ErrUnauthorized) {
				return unauthorized(c)
			}
			return authInternalError(c, err)
		}
		c.Locals(domain.LocalsUserID, userID)
		return c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	scheme, token, found := strings.Cut(strings.TrimSpace(header), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	return token, token != ""
}

func currentUserID(c *fiber.Ctx) string {
	id, _ := c.Locals(domain.LocalsUserID).(string)
	return id
}

func unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "unauthorized"})
}

func authInternalError(c *fiber.Ctx, err error) error {
	if errors.Is(err, domain.ErrAuthNotConfigured) {
		log.Printf("❌ [Auth] JWT_SECRET (>= 32 bytes) or GOOGLE_CLIENT_IDS not configured")
	} else {
		log.Printf("❌ [Auth] internal error: %v", err)
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": "internal server error"})
}
