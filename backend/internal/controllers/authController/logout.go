package authController

import (
	"backend/internal/models/token"

	"github.com/gofiber/fiber/v2"
)

// Logout
func Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	token.RemoveToken(refreshToken)

	clearTokenFromCookie(c)
	return c.JSON(fiber.Map{"message": "logout"})
}