package authController

import (
	"backend/internal/models/token"
	"backend/internal/models/user"
	"backend/internal/models/user/userActivation"
	"backend/internal/models/user/userDTO"
	"backend/internal/requests"
	"backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Account Activation
func ActivateAccount(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refreshToken")
	userDto, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseActivateRequest(c)
	if err != nil {
		return err
	}

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
		token.RemoveAllRefreshTokens(userObj.ID)
		refreshTokenObj, err = token.GenerateRefreshToken(userDto)
		if err != nil {
			return err
		}
	} else {
		refreshTokenObj.RemoveOtherTokens()
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
