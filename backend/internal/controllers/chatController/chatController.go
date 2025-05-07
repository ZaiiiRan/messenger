package chatController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatDTO"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
	"backend/internal/models/shortUser"
	"backend/internal/requests"
	"backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// Create chat
func CreateChat(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseChatRequest(c)
	if err != nil {
		return err
	}

	chat, members, err := chatModel.CreateChat(req.Name, req.Members, req.IsGroup, user)
	if err != nil {
		return err
	}

	members, err = chat.SaveWithMembers(members)
	if err != nil {
		return err
	}

	membersDTOs := chatMemberDTO.CreateChatMembersDTOs(members)

	return c.JSON(fiber.Map{
		"message": "chat created",
		"chat":    chatDTO.CreateChatDTO(chat),
		"members": membersDTOs[1:],
		"you":     membersDTOs[0],
	})
}

// Add members to chat
func AddMembers(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseChatRequest(c)
	if err != nil {
		return err
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

	newMembersDTOs := chatMemberDTO.CreateChatMembersDTOs(newMembers)

	return c.JSON(fiber.Map{
		"message":     "members added",
		"chat":        chatDTO.CreateChatDTO(chat),
		"new_members": newMembersDTOs,
	})
}

// Remove members from chat
func RemoveMembers(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseChatRequest(c)
	if err != nil {
		return err
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

	removedDTOs := chatMemberDTO.CreateChatMembersDTOs(removed)

	return c.JSON(fiber.Map{
		"message":         "members removed",
		"chat":            chatDTO.CreateChatDTO(chat),
		"removed_members": removedDTOs,
	})
}

// Leave from chat
func Leave(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
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
		"chat":    chatDTO.CreateChatDTO(chat),
		"you":     chatMemberDTO.CreateChatMemberDTO(requestSendingMember),
	})
}

// Change member role in chat
func ChatMemberRoleChange(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}
	memberID, err := parseMemberID(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseMemberRoleRequest(c)
	if err != nil {
		return err
	}

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
		"chat":    chatDTO.CreateChatDTO(chat),
		"member":  chatMemberDTO.CreateChatMemberDTO(member),
	})
}

// Return to chat
func ReturnToChat(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
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
		"chat":    chatDTO.CreateChatDTO(chat),
		"you":     chatMemberDTO.CreateChatMemberDTO(requestSendingMember),
	})
}

// Rename chat
func RenameChat(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseChatRequest(c)
	if err != nil {
		return err
	}

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
		"chat":    chatDTO.CreateChatDTO(chat),
	})
}

// Delete chat
func DeleteChat(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
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
		"deleted_chat": chatDTO.CreateChatDTO(chat),
	})
}

// Get chat
func GetChat(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, member, err := chatModel.GetChatAndMember(chatID, user.ID)
	if err != nil {
		return err
	}

	return createChatResponse(c, chat, member)
}

// Get private chat
func GetPrivateChat(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	memberID, err := parseMemberID(c)
	if err != nil {
		return err
	}

	if user.ID == memberID {
		return appErr.BadRequest("you can't have a chat with yourself")
	}

	chat, member, err := chatModel.GetPrivateChatAndMember(user.ID, memberID)
	if err != nil {
		return err
	}

	return createChatResponse(c, chat, member)
}

// Get friends are not chatting
func GetFriendsAreNotChatting(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	req, err := requests.ParseSearchRequest(c)
	if err != nil {
		return err
	}

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

// Get messages from chat (with limit and offset)
func GetMessages(c *fiber.Ctx) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	req, err := requests.ParsePaginationRequest(c)
	if err != nil {
		return err
	}

	chatID, err := parseChatID(c)
	if err != nil {
		return err
	}

	chat, requestSendingMember, err := chatModel.GetChatAndMember(chatID, user.ID)
	if err != nil {
		return err
	}

	messages, err := message.GetMessages(chat, requestSendingMember, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	messagesDTOs := messageDTO.CreateMessagesDTOs(messages)

	return c.JSON(fiber.Map{
		"messages": messagesDTOs,
	})
}

// Get Group chats
func GetGroupChats(c *fiber.Ctx) error {
	return getChatList(c, true)
}

// Get private chats
func GetPrivateChats(c *fiber.Ctx) error {
	return getChatList(c, false)
}

// get chat list (with limit and offset)
func getChatList(c *fiber.Ctx, isGroup bool) error {
	user, err := utils.GetUserDTOFromLocals(c)
	if err != nil {
		return err
	}

	req, err := requests.ParsePaginationRequest(c)
	if err != nil {
		return err
	}

	chats, you, messageIDs, err := chatModel.GetChatList(user.ID, isGroup, req.Limit, req.Offset)
	if err != nil {
		return err
	}

	messages, err := getMessagesByID(messageIDs)
	if err != nil {
		return err
	}

	response, err := createReponseForChatList(chats, you, messages)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"chats": response,
	})
}

// create chat response
func createChatResponse(c *fiber.Ctx, chat *chatModel.Chat, member *chatMember.ChatMember) error {
	var members []chatMember.ChatMember
	if !member.IsRemoved() && !member.IsLeft() {
		var err error
		members, err = chat.GetChatMembers(member)
		if err != nil {
			return err
		}
	}
	membersDTOs := chatMemberDTO.CreateChatMembersDTOs(members)

	lastMessage, err := message.GetLastMessage(chat, member)
	if err != nil {
		return err
	}

	var lastMessageDTO *messageDTO.MessageDTO
	if lastMessage != nil {
		lastMessageDTO = messageDTO.CreateMessageDTO(lastMessage)
	}

	return c.JSON(fiber.Map{
		"chat":         chatDTO.CreateChatDTO(chat),
		"you":          chatMemberDTO.CreateChatMemberDTO(member),
		"members":      membersDTOs,
		"last_message": lastMessageDTO,
	})
}
