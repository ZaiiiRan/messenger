package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Activate Refister request format
type ActivateRequest struct {
	UserID uint64 `json:"user_id"`
	Code   string `json:"code"`
}

// Account Activation
func ActivateAccount(c *fiber.Ctx) error {
	var req ActivateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 	"error": "invalid request format",
		// })
	}
	req.trimSpaces()

	err := user.ActivateAccount(req.UserID, req.Code)
	if err != nil {
		return err
		// status := fiber.StatusBadRequest
		// if err.Error() == "inernal server error" {
		// 	status = fiber.StatusInternalServerError
		// }
		// return c.Status(status).JSON(fiber.Map{
		// 	"error": err.Error(),
		// })
	}

	return c.JSON(fiber.Map{
		"message": "account activated successfully",
	})
}

func (a *ActivateRequest) trimSpaces() {
	a.Code = strings.TrimSpace(a.Code)
}