package socialRoutes

import (
	controller "backend/internal/controllers/socialController"
	"backend/internal/middleware/authMiddleware"

	"github.com/gofiber/fiber/v2"
)

func getUsers(c *fiber.Ctx) error {
	return controller.GetUsers(c)
}

func getFriends(c *fiber.Ctx) error {
	return controller.GetFriends(c)
}

func getIncomingFriendRequests(c *fiber.Ctx) error {
	return controller.GetIncomingFriendRequests(c)
}

func getOutgoingFriendRequests(c *fiber.Ctx) error {
	return controller.GetIncomingFriendRequests(c)
}

func getBlockedUsers(c *fiber.Ctx) error {
	return controller.GetBlockedUsers(c)
}

func SetupRoutes(app fiber.Router) {
	social := app.Group("/social")
	social.Post("/get-users", authMiddleware.AuthMiddleware, getUsers)
	social.Post("/get-friends", authMiddleware.AuthMiddleware, getFriends)
	social.Post("/get-incoming-friend-requests", authMiddleware.AuthMiddleware, getIncomingFriendRequests)
	social.Post("/get-outgoing-friend-requests", authMiddleware.AuthMiddleware, getOutgoingFriendRequests)
	social.Post("/get-blocked-users", authMiddleware.AuthMiddleware, getBlockedUsers)
}
