package chatController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/user"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CreateChatRequest struct {
	Name    string   `json:"name"`
	Members []uint64 `json:"members"`
	IsGroup bool     `json:"is_group"`
}

func (r *CreateChatRequest) TrimSpaces() {
	r.Name = strings.TrimSpace(r.Name)
}

func CreateChat(c *fiber.Ctx) error {
	var req CreateChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.TrimSpaces()

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chat, members, err := chatModel.CreateChat(req.Name, req.Members, req.IsGroup, user)
	if err != nil {
		return err
	}

	err = chat.Save(members)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "chat created",
		"chat":    chat,
		"members": members,
	})
}

type MembersManipulationRequest struct {
	ChatID  uint64   `json:"chat_id"`
	Members []uint64 `json:"members"`
}

func AddMembers(c *fiber.Ctx) error {
	var req MembersManipulationRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	if len(req.Members) == 0 {
		return appErr.BadRequest("empty user list")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	for _, memberID := range req.Members {
		err = chat.AddMember(memberID)
		if err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"message": "members added",
		"chat":    chat,
	})
}
