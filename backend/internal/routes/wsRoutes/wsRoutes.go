package wsRoutes

import (
	"backend/internal/controllers/wsController"
	"backend/internal/middleware/authMiddleware"

	"github.com/gofiber/fiber/v2"
)

func initConnection(c *fiber.Ctx) error {
	return wsController.InitConnection(c)
}

func SetupRoutes(app fiber.Router) {
	ws := app.Group("/ws", authMiddleware.AuthMiddleware)

	ws.Get("/", initConnection)
}
