package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"backend/internal/models/user/userActivation"
	"backend/internal/models/user/userDTO"

	"github.com/gofiber/fiber/v2"
)

// Resend Activation Code
func ResendActivationCode(c *fiber.Ctx) error {
	userDto, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || userDto == nil {
		return appErr.Unauthorized("unauthorized")
	}

	userObj, err := user.GetUserByID(userDto.ID)
	if err != nil {
		return err
	}

	activationCode, err := userActivation.GetActivationCode(userObj)
	if err != nil {
		return err
	}

	activationCode.Regenerate()
	err = activationCode.Save()
	if err != nil {
		return err
	}
	go activationCode.SendToEmail()

	return c.JSON(fiber.Map{
		"message": "new code has been sent",
	})
}
