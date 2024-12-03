package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user/userDTO"

	"github.com/gofiber/fiber/v2"
)

// get user dto from locals (from authMiddleware)
func getUserDTOFromLocals(c *fiber.Ctx) (*userDTO.UserDTO, error) {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return nil, appErr.Unauthorized("unauthorized")
	}
	return user, nil
}
