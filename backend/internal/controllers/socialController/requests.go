package socialController

import (
	appErr "backend/internal/errors/appError"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Fetch users request format
type FetchUsersRequest struct {
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// trim spaces in fetch users request
func (f *FetchUsersRequest) trimSpaces() {
	f.Search = strings.TrimSpace(f.Search)
}

// parse target id from params
func parseTargetID(c *fiber.Ctx) (uint64, error) {
	targetID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return 0, appErr.BadRequest("invalid request format")
	}
	return targetID, nil
}
