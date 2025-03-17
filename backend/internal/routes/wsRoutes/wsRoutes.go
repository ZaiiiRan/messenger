package wsRoutes

import (
	"backend/internal/controllers/wsController"
	"backend/internal/middleware/wsUpgradeMiddleware"

	"github.com/gofiber/fiber/v2"
)

func initConnection(c *fiber.Ctx) error {
	return wsController.InitConnection(c)
}

func SetupRoutes(app fiber.Router) {
	ws := app.Group("/ws")
	ws.Use(wsUpgradeMiddleware.WSUpgradeMiddleware)

	ws.Get("/", initConnection)
}
