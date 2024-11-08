package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

// Resend Activation Code
func ResendActivationCode(c *fiber.Ctx) error {
	userDTO, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || userDTO == nil {
		return appErr.Unauthorized("unauthorized")
	}

	activationCode, err := user.GetActivationCode(userDTO.ID)
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
