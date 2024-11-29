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
	social.Post("/users/search", authMiddleware.AuthMiddleware, getUsers)
	social.Get("/users/:id", authMiddleware.AuthMiddleware, getUser)

	// Friends
	social.Post("/friends/friend-list", authMiddleware.AuthMiddleware, getFriends)

	// Friend requests
	social.Post("/friends/friend-requests/incoming", authMiddleware.AuthMiddleware, getIncomingFriendRequests)
	social.Post("/friends/friend-requests/outgoing", authMiddleware.AuthMiddleware, getOutgoingFriendRequests)

	// Block list
	social.Post("/block/block-list", authMiddleware.AuthMiddleware, getBlockedUsers)

	// Friend management
	social.Post("/friends/management/:id", authMiddleware.AuthMiddleware, addFriend)
	social.Delete("/friends/management/:id", authMiddleware.AuthMiddleware, removeFriend)

	// Block/Unblock
	social.Post("/block/management/:id", authMiddleware.AuthMiddleware, blockUser)
	social.Delete("/block/management/:id", authMiddleware.AuthMiddleware, unblockUser)

}
