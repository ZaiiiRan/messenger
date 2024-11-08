package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/token"

	"github.com/gofiber/fiber/v2"
)

// Logout
func Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return appErr.Unauthorized("unauthorized")
	}

	token.RemoveToken(refreshToken)

	clearTokenFromCookie(c)
	return c.JSON(fiber.Map{"message": "logout"})
}