package chatController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chatMember"
	"backend/internal/models/user"

	"github.com/gofiber/fiber/v2"
)

func CreateChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()

	chat, members, err := chatModel.CreateChat(req.Name, req.Members, req.IsGroup, user)
	if err != nil {
		return err
	}

	_, err = chat.SaveWithMembers(members)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "chat created",
		"chat":    chat,
	})
}

func AddMembers(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	if len(req.Members) == 0 {
		return appErr.BadRequest("empty user list")
	}

	chat, requestSendingMember, err := getChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	newMembers, err := chat.AddMembers(req.Members, requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":     "members added",
		"chat":        chat,
		"new_members": newMembers,
	})
}

func Leave(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, requestSendingMember, err := getChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	_, err = chat.LeaveFromChat(requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "you have left the chat",
		"chat":    chat,
	})
}

func ReturnToChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, requestSendingMember, err := chatModel.GetChatAndMember(chatID, user.ID)
	if err != nil {
		return err
	}

	_, err = chat.ReturnToChat(requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "you have returned to chat",
		"chat":    chat,
	})
}

func RemoveMembers(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}

	if len(req.Members) == 0 {
		return appErr.BadRequest("empty user list")
	}

	chat, requestSendingMember, err := getChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	removed, err := chat.RemoveMembers(req.Members, requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":         "members removed",
		"chat":            chat,
		"removed_members": removed,
	})
}

func RenameChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()

	chat, requestSendingMember, err := getChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	err = chat.Rename(req.Name, requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "chat renamed",
		"chat":    chat,
	})
}

func ChatMemberRoleChange(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}
	memberID, err := parseMemberID(c)
	if err != nil {
		return err
	}

	var req ChangeRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()

	chat, requestSendingMember, err := getChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	member, err := chat.ChatMemberRoleChange(memberID, req.Role, requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "role changed",
		"chat":    chat,
		"member":  member,
	})
}

func DeleteChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*user.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, requestSendingMember, err := getChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	err = chat.Delete(requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "chat has been deleted",
		"deleted_chat": chat,
	})
}

// get chat object and verify user access
func getChatAndVerifyAccess(chatID, userID uint64) (*chatModel.Chat, *chatMember.ChatMember, error) {
	chat, member, err := chatModel.GetChatAndMember(chatID, userID)
	if err != nil {
		return nil, nil, err
	}
	if member.Removed() {
		return nil, nil, appErr.Forbidden("you cannot access this chat")
	}
	return chat, member, nil
}
