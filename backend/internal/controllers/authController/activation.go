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

	refreshTokenObj, _ := token.FindRefreshToken(refreshToken)
	if refreshTokenObj == nil {
		refreshTokenObj.RemoveOtherTokens()
		refreshTokenObj, err = token.GenerateRefreshToken(userDTO)
		if err != nil {
			return err
		}
	} else {
		err = refreshTokenObj.RegenerateRefreshToken(userDTO)
		if err != nil {
			return err
		}
	}
	err = refreshTokenObj.SaveRefreshToken()
	if err != nil {
		return err
	}

	newAccessToken, err := token.GenerateAccessToken(userDTO)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(userDTO, newAccessToken, refreshTokenObj.RefreshToken, c)
}
