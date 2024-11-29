package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/token"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

// Account Activation
func ActivateAccount(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	userDTO, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || userDTO == nil {
		return appErr.Unauthorized("unauthorized")
	}
	var req ActivateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()

	err := user.ActivateAccount(userDTO.ID, req.Code)
	if err != nil {
		return err
	}

	userDTO.IsActivated = true

	newAccessToken, newRefreshToken, err := token.GenerateTokens(userDTO)
	if err != nil {
		return err
	}
	_, err = token.UpdateToken(refreshToken, newRefreshToken, userDTO.ID)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(userDTO, newAccessToken, newRefreshToken, c)
}
