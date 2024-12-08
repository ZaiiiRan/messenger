package requests

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type MemberRoleRequest struct {
	Role string `json:"role"`
}

// trim spaces for member role request
func (r *MemberRoleRequest) TrimSpaces() {
	r.Role = strings.TrimSpace(r.Role)
}

// parse member role request
func ParseMemberRoleRequest(c *fiber.Ctx) (*MemberRoleRequest, error) {
	return ParseRequest[MemberRoleRequest](c)
}
