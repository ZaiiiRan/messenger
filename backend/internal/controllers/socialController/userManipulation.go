package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/socialUser"

	"github.com/gofiber/fiber/v2"
)

// User manipulation request format
type UserManipulationRequest struct {
	UserID uint64 `json:"user_id"`
}

// Add friend
func AddFriend(c *fiber.Ctx) error {
	var req UserManipulationRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, err := getUserDTOFromLocals(c)
	if err != nil {
		return err
	}
	err = socialUser.AddFriend(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "ok",
	})
}

// Remove friend
func RemoveFriend(c *fiber.Ctx) error {
	var req UserManipulationRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, err := getUserDTOFromLocals(c)
	if err != nil {
		return err
	}
	err = socialUser.RemoveFriend(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "friend deleted",
	})
}

// Block user
func BlockUser(c *fiber.Ctx) error {
	var req UserManipulationRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, err := getUserDTOFromLocals(c)
	if err != nil {
		return err
	}
	err = socialUser.BlockUser(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "user blocked",
	})
}

// Unblock user
func UnblockUser(c *fiber.Ctx) error {
	var req UserManipulationRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	user, err := getUserDTOFromLocals(c)
	if err != nil {
		return err
	}
	err = socialUser.UnblockUser(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "user unblocked",
	})
}
