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

	members, err = chat.SaveWithMembers(members)
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

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	if len(req.Members) == 0 {
		return appErr.BadRequest("empty user list")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	newMembers, err := chat.AddMembers(req.Members, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":     "members added",
		"chat":        chat,
		"new_members": newMembers,
	})
}

type LeaveRequest struct {
	ChatID uint64 `json:"chat_id"`
}

func Leave(c *fiber.Ctx) error {
	var req LeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	_, err = chat.RemoveMembers([]uint64{user.ID}, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "you have left the chat",
		"chat":    chat,
	})
}

func ReturnToChat(c *fiber.Ctx) error {
	var req LeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	_, err = chat.ReturnToChat(user.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "you have returned to chat",
		"chat":    chat,
	})
}

func RemoveMembers(c *fiber.Ctx) error {
	var req MembersManipulationRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	if len(req.Members) == 0 {
		return appErr.BadRequest("empty user list")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	removed, err := chat.RemoveMembers(req.Members, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":         "members removed",
		"chat":            chat,
		"removed_members": removed,
	})
}

type RenameChatRequest struct {
	ChatID uint64 `json:"chat_id"`
	Name   string `json:"name"`
}

func RenameChat(c *fiber.Ctx) error {
	var req RenameChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	err = chat.Rename(req.Name, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "chat renamed",
		"chat":    chat,
	})
}

type ChangeRoleRequest struct {
	ChatID   uint64 `json:"chat_id"`
	MemberID uint64 `json:"member_id"`
	Role     string `json:"role"`
}

func ChatMemberRoleChange(c *fiber.Ctx) error {
	var req ChangeRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chat, err := chatModel.GetChatByID(req.ChatID)
	if err != nil {
		return err
	}

	member, err := chat.ChatMemberRoleChange(req.MemberID, req.Role, user.ID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "chat renamed",
		"chat":    chat,
		"member":  member,
	})
}
