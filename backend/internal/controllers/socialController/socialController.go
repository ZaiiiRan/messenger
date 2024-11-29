package socialController

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"strconv"

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

// parse target id from params
func parseTargetID(c *fiber.Ctx) (uint64, error) {
	targetID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, appErr.BadRequest("invalid request format")
	}
	return targetID, nil
}
