package chatRoutes

import (
	controller "backend/internal/controllers/chatController"
	"backend/internal/middleware/authMiddleware"

	"github.com/gofiber/fiber/v2"
)

func createChat(c *fiber.Ctx) error {
	return controller.CreateChat(c)
} 

func addMembers(c *fiber.Ctx) error {
	return controller.AddMembers(c)
}

func leave(c *fiber.Ctx) error {
	return controller.Leave(c)
} 

func SetupRoutes(app *fiber.App) {
	chat := app.Group("/chat")
	chat.Post("/create-chat", authMiddleware.AuthMiddleware, createChat)
	chat.Post("/add-members", authMiddleware.AuthMiddleware, addMembers)
	chat.Post("/leave", authMiddleware.AuthMiddleware, leave)
}

