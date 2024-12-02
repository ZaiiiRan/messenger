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

	clearTokenFromCookie(c)
	refreshTokenObj, err := token.FindRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	var userDTO *user.UserDTO
	userDTO, err = refreshTokenObj.ValidateRefreshToken()
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

	err = refreshTokenObj.RegenerateRefreshToken(newUserDTO)
	if err != nil {
		return err
	}
	err = refreshTokenObj.SaveRefreshToken()
	if err != nil {
		return err
	}
	newAccessToken, err := token.GenerateAccessToken(newUserDTO)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(newUserDTO, newAccessToken, refreshTokenObj.RefreshToken, c)
}
