package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/socialUser"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// User manipulation request format
type UserManipulationRequest struct {
	UserID uint64 `json:"user_id"`
}

func checkSelfID(selfID, targetID uint64) error {
	if selfID == targetID {
		return appErr.BadRequest("invalid user id")
	}
	return nil
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
	if err := checkSelfID(user.ID, req.UserID); err != nil {
		return err
	}
	friend, err := socialUser.AddFriend(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "ok",
		"user": friend,
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
	if err := checkSelfID(user.ID, req.UserID); err != nil {
		return err
	}
	target, err := socialUser.RemoveFriend(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "friend deleted",
		"user": target,
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
	if err := checkSelfID(user.ID, req.UserID); err != nil {
		return err
	}
	target, err := socialUser.BlockUser(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "user blocked",
		"user": target,
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
	if err := checkSelfID(user.ID, req.UserID); err != nil {
		return err
	}
	target, err := socialUser.UnblockUser(user.ID, req.UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "user unblocked",
		"user": target,
	})
}

// Get User dto with friend status
func GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return appErr.BadRequest("invalid request format")
	}

	user, err := getUserDTOFromLocals(c)
	if err != nil {
		return err
	}
	if err := checkSelfID(user.ID, id); err != nil {
		return err
	}
	dto, err := socialUser.GetTargetByID(user.ID, id)
	if err != nil {
		return err
	}
	return c.JSON(dto)
}
