package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

// get user dto from locals (from authMiddleware)
func getUserDTOFromLocals(c *fiber.Ctx) (*user.UserDTO, error) {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return nil, appErr.Unauthorized("unauthorized")
	}
	return user, nil
}
