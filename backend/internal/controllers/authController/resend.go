package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

// Resend activation code request format
type ResendActivationCodeRequest struct {
	UserID uint64 `json:"user_id"`
}

// Resend Activation Code
func ResendActivationCode(c *fiber.Ctx) error {
	var req ResendActivationCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	activationCode, err := user.GetActivationCode(req.UserID)
	if err != nil {
		return err
	}

	activationCode.Regenerate()
	err = activationCode.Save()
	if err != nil {
		return err
	}
	activationCode.SendToEmail()

	return c.JSON(fiber.Map{
		"message": "new code has been sent",
	})
}
