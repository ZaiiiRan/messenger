package messageController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
)

// parse request
func parseRequest[T any](request interface{}) (*T, error) {
	req, ok := request.(*T)
	if !ok || req == nil {
		return nil, appErr.BadRequest("invalid payload")
	}
	return req, nil
}

// verify access and get message
func verifyAccessAndGetMessage(chatID, actorID, messageID uint64) (*chatModel.Chat, *chatMember.ChatMember, *message.Message, error) {
	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(chatID, actorID)
	if err != nil {
		return nil, nil, nil, err
	}

	message, err := message.GetMessage(messageID)
	if err != nil {
		return nil, nil, nil, err
	}

	return chat, requestSendingMember, message, nil
}

// create response
func createResponse(chat *chatModel.Chat, requestSendingMember *chatMember.ChatMember, message *message.Message) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	messageDto := messageDTO.CreateMessageDTO(message)

	members, err := chat.GetChatMembers(requestSendingMember)
	if err != nil {
		return nil, nil, err
	}
	membersDTOs := chatMemberDTO.CreateChatMembersDTOs(members)

	return messageDto, membersDTOs, nil
}
