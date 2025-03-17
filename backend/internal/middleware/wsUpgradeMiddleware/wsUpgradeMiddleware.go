package wsUpgradeMiddleware

import (
	"backend/internal/errors/appError"
	"backend/internal/models/token"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WSUpgradeMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	accessToken := c.Query("token")
	if accessToken == "" {
		return appError.Unauthorized("unauthorized")
	}

	user, err := token.ValidateAccessToken(accessToken)
	if err != nil || user == nil {
		return err
	}

	c.Locals("userDTO", user)
	c.Locals("allowed", true)
	return c.Next()
}
