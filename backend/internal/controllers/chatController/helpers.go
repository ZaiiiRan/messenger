package chatController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatDTO"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"errors"
)

// create Chat response array for chat list
func createReponseForChatList(chats []chatModel.Chat, you []*chatMember.ChatMember, messages []*message.Message) ([]ChatResponse, error) {
	var response []ChatResponse
	for index, chat := range chats {
		var membersDTOs []*chatMemberDTO.ChatMemberDTO

		if !you[index].IsRemoved() && !you[index].IsLeft() {
			members, err := chat.GetChatMembers(you[index])
			if err != nil {
				return nil, err
			}
			membersDTOs = chatMemberDTO.CreateChatMembersDTOs(members)
		} else {
			membersDTOs = nil
		}

		var msgDTO *messageDTO.MessageDTO
		if messages[index] != nil {
			msgDTO = messageDTO.CreateMessageDTO(messages[index])
		}

		responsePart := ChatResponse{
			Chat:        chatDTO.CreateChatDTO(&chat),
			Members:     membersDTOs,
			You:         chatMemberDTO.CreateChatMemberDTO(you[index]),
			LastMessage: msgDTO,
		}

		response = append(response, responsePart)
	}

	return response, nil
}

// get last messages by id for chat list
func getMessagesByID(messageIDs []*uint64) ([]*message.Message, error) {
	var messages []*message.Message
	for _, messageID := range messageIDs {
		var appError *appErr.AppError
		if messageID != nil {
			msg, err := message.GetMessage(*messageID)
			if err != nil && (errors.As(err, &appError) && appError.StatusCode == 404) {
				msg = nil
			} else if err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		} else {
			messages = append(messages, nil)
		}

	}

	return messages, nil
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
