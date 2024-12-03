package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/token"
	"backend/internal/models/user"
	"backend/internal/models/user/userActivation"
	"backend/internal/models/user/userDTO"

	"github.com/gofiber/fiber/v2"
)

// Account Activation
func ActivateAccount(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	userDto, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || userDto == nil {
		return appErr.Unauthorized("unauthorized")
	}
	var req ActivateRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()

	userObj, err := user.GetUserByID(userDto.ID)
	if err != nil {
		return err
	}

	err = userActivation.ActivateAccount(userObj, req.Code)
	if err != nil {
		return err
	}

	userDto = userDTO.CreateUserDTOFromUserObj(userObj)
	refreshTokenObj, _ := token.FindRefreshToken(refreshToken)
	if refreshTokenObj == nil {
		refreshTokenObj.RemoveOtherTokens()
		refreshTokenObj, err = token.GenerateRefreshToken(userDto)
		if err != nil {
			return err
		}
	} else {
		err = refreshTokenObj.RegenerateRefreshToken(userDto)
		if err != nil {
			return err
		}
	}
	err = refreshTokenObj.SaveRefreshToken()
	if err != nil {
		return err
	}
	newAccessToken, err := token.GenerateAccessToken(userDto)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(userDto, newAccessToken, refreshTokenObj.RefreshToken, c)
}
