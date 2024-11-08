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

	userObj, err := user.GetUserByID(userDTO.ID)
	if err != nil {
		return err
	}

	newUserDTO := user.CreateUserDTOFromUserObj(userObj)

	newAccessToken, newRefreshToken, err := token.GenerateTokens(newUserDTO)
	if err != nil {
		return err
	}
	_, err = token.UpdateToken(refreshToken, newRefreshToken, userDTO.ID)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(newUserDTO, newAccessToken, newRefreshToken, c)
}
