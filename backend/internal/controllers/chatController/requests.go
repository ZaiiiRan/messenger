package chatController

import (
	appErr "backend/internal/errors/appError"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
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

// Get chat members request format
type GetChatMembersRequest struct {
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// trim spaces for get chat members request
func (r *GetChatMembersRequest) trimSpaces() {
	r.Search = strings.TrimSpace(r.Search)
}

// parse chat id from params
func parseChatID(c *fiber.Ctx) (uint64, error) {
	chatIDParam := c.Params("chat_id")
	chatID, err := strconv.ParseUint(chatIDParam, 0, 64)
	if err != nil {
		return 0, appErr.BadRequest("invalid request format")
	}
	return chatID, nil
}

// parse member id from params
func parseMemberID(c *fiber.Ctx) (uint64, error) {
	memberIDParam := c.Params("member_id")
	memberID, err := strconv.ParseUint(memberIDParam, 0, 64)
	if err != nil {
		return 0, appErr.BadRequest("invalid request format")
	}
	return memberID, nil
}

// Get chat members or chat request format
type GetListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
