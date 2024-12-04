package messageController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
	"backend/internal/models/user/userDTO"
	"strings"
)

type SendMessageReq struct {
	ChatID         uint64 `json:"chat_id"`
	MessageContent string `json:"message_content"`
}

// Send message request trim spaces
func (r *SendMessageReq) TrimSpaces() {
	r.MessageContent = strings.TrimSpace(r.MessageContent)
}

func SendMessage(userDto *userDTO.UserDTO, request interface{}) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	req, ok := request.(*SendMessageReq)
	if !ok || req == nil {
		return nil, nil, appErr.BadRequest("invalid send_message payload")
	}
	if req.MessageContent == "" {
		return nil, nil, appErr.BadRequest("message is empty")
	}

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(req.ChatID, userDto.ID)
	if err != nil {
		return nil, nil, err
	}

	message, err := message.NewMessage(chat, requestSendingMember, req.MessageContent)
	if err != nil {
		return nil, nil, err
	}
	messageDto := messageDTO.NewMessageDTO(message)

	members, err := chat.GetChatMembers(requestSendingMember)
	if err != nil {
		return nil, nil, err
	}
	membersDTOs := chatMemberDTO.GetChatMembersDTOs(members)

	return messageDto, membersDTOs, nil
}
