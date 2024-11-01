package authController

import (
	"backend/internal/models/user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Login request format
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Login
func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request format",
		})
	}
	req.trimSpaces()

	userObject, err := user.GetUserByUsername(req.Login)
	if err != nil && err.Error() == "user not found" {
		userObject, err = user.GetUserByEmail(req.Login)
		if err != nil && err.Error() == "user not found" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid login or password",
			})
		}
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if !userObject.CheckPassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid login or password",
		})
	}

	return createUserDTOAndTokensResponse(userObject, c)
}

func (l *LoginRequest) trimSpaces() {
	l.Login = strings.TrimSpace(l.Login)
	l.Password = strings.TrimSpace(l.Password)
}