package socialRoutes

import (
	controller "backend/internal/controllers/socialController"
	"backend/internal/middleware/authMiddleware"

	"github.com/gofiber/fiber/v2"
)

func getUsers(c *fiber.Ctx) error {
	return controller.GetUsers(c)
}

func getUser(c *fiber.Ctx) error {
	return controller.GetUser(c)
}

func getFriends(c *fiber.Ctx) error {
	return controller.GetFriends(c)
}

func getIncomingFriendRequests(c *fiber.Ctx) error {
	return controller.GetIncomingFriendRequests(c)
}

func getOutgoingFriendRequests(c *fiber.Ctx) error {
	return controller.GetOutgoingFriendRequests(c)
}

func getBlockedUsers(c *fiber.Ctx) error {
	return controller.GetBlockedUsers(c)
}

func addFriend(c *fiber.Ctx) error {
	return controller.AddFriend(c)
}

func removeFriend(c *fiber.Ctx) error {
	return controller.RemoveFriend(c)
}

func blockUser(c *fiber.Ctx) error {
	return controller.BlockUser(c)
}

func unblockUser(c *fiber.Ctx) error {
	return controller.UnblockUser(c)
}

func SetupRoutes(app fiber.Router) {
	social := app.Group("/social")

	// Users
	social.Post("/users", authMiddleware.AuthMiddleware, getUsers)
	social.Get("/user/:id", authMiddleware.AuthMiddleware, getUser)

	// Friends
	social.Post("/friends", authMiddleware.AuthMiddleware, getFriends)

	// Friend requests
	social.Post("/friend-requests/incoming", authMiddleware.AuthMiddleware, getIncomingFriendRequests)
	social.Post("/friend-requests/outgoing", authMiddleware.AuthMiddleware, getOutgoingFriendRequests)

	// Block list
	social.Post("/block-list", authMiddleware.AuthMiddleware, getBlockedUsers)

	// Friend management
	social.Post("/users/:id/friend", authMiddleware.AuthMiddleware, addFriend)
	social.Delete("/users/:id/friend", authMiddleware.AuthMiddleware, removeFriend)

	// Block/Unblock
	social.Post("/users/:id/block", authMiddleware.AuthMiddleware, blockUser)
	social.Delete("/users/:id/block", authMiddleware.AuthMiddleware, unblockUser)

}
