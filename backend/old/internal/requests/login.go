package requests

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Login request format
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// trim spaces in login request
func (r *LoginRequest) TrimSpaces() {
	r.Login = strings.TrimSpace(r.Login)
	r.Password = strings.TrimSpace(r.Password)
}

// parse login request
func ParseLoginRequest(c *fiber.Ctx) (*LoginRequest, error) {
	return ParseRequest[LoginRequest](c)
}