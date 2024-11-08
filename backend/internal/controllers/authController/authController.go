package authController

import (
	"backend/internal/models/token"
	"backend/internal/models/user"
	"time"

	"github.com/gofiber/fiber/v2"
)

// creating user dto and tokens for response
func createUserDTOAndTokensResponse(userObject *user.User, c *fiber.Ctx) error {
	userDTO := user.CreateUserDTOFromUserObj(userObject)
	accessToken, refreshToken, err := token.GenerateTokens(userDTO)
	if err != nil {
		return err
	}
	_, err = token.InsertToken(userDTO.ID, refreshToken)
	if err != nil {
		return err
	}

	return sendTokenAndJSON(userDTO, accessToken, refreshToken, c)
}

// sending tokens and user dto to client
func sendTokenAndJSON(userDTO *user.UserDTO, accessToken, refreshToken string, c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HTTPOnly: true,
		Path:     "/",
		SameSite: "None",
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
	})
}
