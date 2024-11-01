package authController

import (
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}

	activationCode, err := user.GetActivationCode(req.UserID)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "inernal server error" {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	activationCode.Regenerate()
	err = activationCode.Save()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	activationCode.SendToEmail()

	return c.JSON(fiber.Map{
		"message": "new code has been sent",
	})
}
