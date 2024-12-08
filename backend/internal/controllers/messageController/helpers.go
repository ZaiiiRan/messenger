package messageController

import (
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
)

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
