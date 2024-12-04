package chatController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/shortUser"
	"backend/internal/models/user/userDTO"

	"github.com/gofiber/fiber/v2"
)

func CreateChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
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
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
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

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	newMembers, err := chat.AddMembers(req.Members, requestSendingMember)
	if err != nil {
		return err
	}

	newMembersDTOs := chatMemberDTO.GetChatMembersDTOs(newMembers)

	return c.JSON(fiber.Map{
		"message":     "members added",
		"chat":        chat,
		"new_members": newMembersDTOs,
	})
}

func Leave(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
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
		"you":     chatMemberDTO.CreateChatMemberDTO(requestSendingMember),
	})
}

func ReturnToChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
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
		"you":     chatMemberDTO.CreateChatMemberDTO(requestSendingMember),
	})
}

func RemoveMembers(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
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

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	removed, err := chat.RemoveMembers(req.Members, requestSendingMember)
	if err != nil {
		return err
	}

	removedDTOs := chatMemberDTO.GetChatMembersDTOs(removed)

	return c.JSON(fiber.Map{
		"message":         "members removed",
		"chat":            chat,
		"removed_members": removedDTOs,
	})
}

func RenameChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
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

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
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
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
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

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
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
		"member":  chatMemberDTO.CreateChatMemberDTO(member),
	})
}

func DeleteChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	err = chat.Delete(requestSendingMember)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message":      "chat has been deleted",
		"deleted_chat": chat,
	})
}

func GetChat(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, member, err := chatModel.GetChatAndMember(chatID, user.ID)
	if err != nil {
		return err
	}

	var members []chatMember.ChatMember
	if !member.IsRemoved() && !member.IsLeft() {
		members, err = chat.GetChatMembers(member)
		if err != nil {
			return err
		}
	}
	membersDTOs := chatMemberDTO.GetChatMembersDTOs(members)

	return c.JSON(fiber.Map{
		"chat":    chat,
		"you":     chatMemberDTO.CreateChatMemberDTO(member),
		"members": membersDTOs,
	})
}

func GetFriendsAreNotChatting(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}

	var req GetChatMembersRequest
	if err := c.BodyParser(&req); err != nil {
		return appErr.BadRequest("invalid request format")
	}
	req.trimSpaces()

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, _, err := chatModel.GetChatAndVerifyAccess(chatID, user.ID)
	if err != nil {
		return err
	}

	friends, err := shortUser.SearchFriendsAreNotChatting(user.ID, chat.ID, req.Search, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"friends": friends,
	})
}
