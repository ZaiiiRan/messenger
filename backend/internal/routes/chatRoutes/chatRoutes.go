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

func removeMembers(c *fiber.Ctx) error {
	return controller.RemoveMembers(c)
}

func returnToChat(c *fiber.Ctx) error {
	return controller.ReturnToChat(c)
}

func renameChat(c *fiber.Ctx) error {
	return controller.RenameChat(c)
}

func chatMemberRoleChange(c *fiber.Ctx) error {
	return controller.ChatMemberRoleChange(c)
}

func SetupRoutes(app *fiber.App) {
	chat := app.Group("/chat")
	chat.Post("/create-chat", authMiddleware.AuthMiddleware, createChat)
	chat.Post("/add-members", authMiddleware.AuthMiddleware, addMembers)
	chat.Post("/leave", authMiddleware.AuthMiddleware, leave)
	chat.Post("/remove-members", authMiddleware.AuthMiddleware, removeMembers)
	chat.Post("/return", authMiddleware.AuthMiddleware, returnToChat)
	chat.Post("/rename", authMiddleware.AuthMiddleware, renameChat)
	chat.Post("change-role", authMiddleware.AuthMiddleware, chatMemberRoleChange)
}
