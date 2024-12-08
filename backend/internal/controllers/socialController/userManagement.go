package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/socialUser"
	"backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func checkSelfID(selfID, targetID uint64) error {
	if selfID == targetID {
		return appErr.BadRequest("invalid user id")
	}
	return nil
}

// Add friend
func AddFriend(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	targetID, err := parseTargetID(c)
	if err != nil {
		return err
	}

	if err := checkSelfID(user.ID, targetID); err != nil {
		return err
	}
	friend, err := socialUser.AddFriend(user.ID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "ok",
		"user":    friend,
	})
}

// Remove friend
func RemoveFriend(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	targetID, err := parseTargetID(c)
	if err != nil {
		return err
	}

	if err := checkSelfID(user.ID, targetID); err != nil {
		return err
	}
	target, err := socialUser.RemoveFriend(user.ID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "friend deleted",
		"user":    target,
	})
}

// Block user
func BlockUser(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	targetID, err := parseTargetID(c)
	if err != nil {
		return err
	}

	if err := checkSelfID(user.ID, targetID); err != nil {
		return err
	}
	target, err := socialUser.BlockUser(user.ID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "user blocked",
		"user":    target,
	})
}

// Unblock user
func UnblockUser(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	targetID, err := parseTargetID(c)
	if err != nil {
		return err
	}

	if err := checkSelfID(user.ID, targetID); err != nil {
		return err
	}
	target, err := socialUser.UnblockUser(user.ID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "user unblocked",
		"user":    target,
	})
}

// Get User dto with friend status
func GetUser(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	targetID, err := parseTargetID(c)
	if err != nil {
		return err
	}

	if err := checkSelfID(user.ID, targetID); err != nil {
		return err
	}
	dto, err := socialUser.GetTargetByID(user.ID, targetID)
	if err != nil {
		return err
	}

	return c.JSON(dto)
}
