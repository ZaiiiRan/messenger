package requests

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ChatRequest struct {
	Name    string   `json:"name"`
	Members []uint64 `json:"members"`
	IsGroup bool     `json:"is_group"`
}

// trim spaces for chat request
func (r *ChatRequest) TrimSpaces() {
	r.Name = strings.TrimSpace(r.Name)
}

// Parse chat request
func ParseChatRequest(c *fiber.Ctx) (*ChatRequest, error) {
	return ParseRequest[ChatRequest](c)
}
