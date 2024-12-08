package requests

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Activate Register request format
type ActivateRequest struct {
	Code string `json:"code"`
}

// trim spaces in activation request
func (r *ActivateRequest) TrimSpaces() {
	r.Code = strings.TrimSpace(r.Code)
}

// parse activate request
func ParseActivateRequest(c *fiber.Ctx) (*ActivateRequest, error) {
	return ParseRequest[ActivateRequest](c)
}
