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

func deleteChat(c *fiber.Ctx) error {
	return controller.DeleteChat(c)
}

func SetupRoutes(app *fiber.App) {
	chat := app.Group("/chats")

	// Chat
	chat.Post("/", authMiddleware.AuthMiddleware, createChat)
	chat.Patch("/:chat_id", authMiddleware.AuthMiddleware, renameChat)
	chat.Delete("/:chat_id", authMiddleware.AuthMiddleware, deleteChat)

	// Members management
	chat.Post("/:chat_id/members", authMiddleware.AuthMiddleware, addMembers)
	chat.Delete("/:chat_id/members", authMiddleware.AuthMiddleware, removeMembers)
	chat.Patch("/:chat_id/members/:member_id/role", authMiddleware.AuthMiddleware, chatMemberRoleChange)

	// Leave/Return
	chat.Patch("/:chat_id/leave", authMiddleware.AuthMiddleware, leave)
	chat.Patch("/:chat_id/return", authMiddleware.AuthMiddleware, returnToChat)
}
