package authController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

// Login
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid login or password")
	}
	req.trimSpaces()

	if req.Login == "" {
		return appErr.BadRequest("login is empty")
	}
	if req.Password == "" {
		return appErr.BadRequest("password is empty")
	}

	userObject, err := user.GetUserByUsername(req.Login)
	if err != nil && err.Error() == "user not found" {
		userObject, err = user.GetUserByEmail(req.Login)
		if err != nil {
			return appErr.BadRequest("invalid login or password")
		}
	} else if err != nil {
		return err
	}

	if !userObject.CheckPassword(req.Password) {
		return appErr.BadRequest("invalid login or password")
	}

	return createUserDTOAndTokensResponse(userObject, c)
}
