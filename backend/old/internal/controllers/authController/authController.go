package authController

import (
	"backend/internal/models/token"
	"backend/internal/models/user"
	"backend/internal/models/user/userDTO"
	"time"

	"github.com/gofiber/fiber/v2"
)

// creating user dto and tokens for response
func createUserDTOAndTokensResponse(userObject *user.User, c *fiber.Ctx) error {
	userDTO := userDTO.CreateUserDTOFromUserObj(userObject)

	refreshToken, err := token.GenerateRefreshToken(userDTO)
	if err != nil {
		return err
	}
	err = refreshToken.SaveRefreshToken()
	if err != nil {
		return err
	}
	accessToken, err := token.GenerateAccessToken(userDTO)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(userDTO, accessToken, refreshToken.RefreshToken, c)
}

// sending tokens and user dto to client
func sendTokenAndJSON(userDTO *userDTO.UserDTO, accessToken, refreshToken string, c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{
		"user":        userDTO,
		"accessToken": accessToken,
	})
}

// clear refresh token from cookie
func clearTokenFromCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now(),
		HTTPOnly: true,
		Path:     "/",
		SameSite: "Lax",
	})
}
