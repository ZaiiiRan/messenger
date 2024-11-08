package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"backend/internal/models/token"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Activate Refister request format
type ActivateRequest struct {
	Code   string `json:"code"`
}

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

func (a *ActivateRequest) trimSpaces() {
	a.Code = strings.TrimSpace(a.Code)
}