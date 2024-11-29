package chatController

import (
	"strings"
)

// Chat request format
type ChatRequest struct {
	Name    string   `json:"name"`
	Members []uint64 `json:"members"`
	IsGroup bool     `json:"is_group"`
}

// trim spaces for chat request
func (r *ChatRequest) trimSpaces() {
	r.Name = strings.TrimSpace(r.Name)
}

// Change member role request format
type ChangeRoleRequest struct {
	Role string `json:"role"`
}

// trim spaces for change member role request
func (r *ChangeRoleRequest) trimSpaces() {
	r.Role = strings.TrimSpace(r.Role)
}
