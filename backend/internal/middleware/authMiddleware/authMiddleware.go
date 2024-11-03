package authMiddleware

import (
	"backend/internal/errors/appError"
	"backend/internal/models/token"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")
	
	if len(authorizationHeader) < len("Bearer ") {
		return appError.Unauthorized("unauthorized")
	}

	accessToken := authorizationHeader[len("Bearer "):]
	if accessToken == "" {
		return appError.Unauthorized("unauthorized")
	}

	user, err := token.ValidateAccessToken(accessToken)
	if err != nil || user == nil {
		return err
	}

	c.Locals("userDTO", user)
	return c.Next()
}