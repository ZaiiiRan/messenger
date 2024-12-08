package socialController

import (
	appErr "backend/internal/errors/appError"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// parse target id from params
func parseTargetID(c *fiber.Ctx) (uint64, error) {
	targetID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, appErr.BadRequest("invalid request format")
	}
	return targetID, nil
}
