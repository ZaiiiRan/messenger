package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/token"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

// Refresh
func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	if refreshToken == "" {
		return appErr.Unauthorized("unauthorized")
	}

	var userDTO *user.UserDTO
	userDTO, err := token.ValidateRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	if userDTO == nil {
		return appErr.Unauthorized("unauthorized")
	}

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
