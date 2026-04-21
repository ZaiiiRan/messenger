package wsUpgradeMiddleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WSUpgradeMiddleware(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	c.Locals("allowed", true)
	return c.Next()
}
