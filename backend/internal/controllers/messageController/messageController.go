package messageController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
	"backend/internal/models/user/userDTO"
)

// Send message
func SendMessage(userDto *userDTO.UserDTO, request interface{}) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	req, err := parseRequest[SendMessageReq](request)
	if err != nil {
		return nil, nil, err
	}
	req.TrimSpaces()
	if req.MessageContent == "" {
		return nil, nil, appErr.BadRequest("message content is empty")
	}

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(req.ChatID, userDto.ID)
	if err != nil {
		return nil, nil, err
	}

	message, err := message.NewMessage(chat, requestSendingMember, req.MessageContent)
	if err != nil {
		return nil, nil, err
	}

	return createResponse(chat, requestSendingMember, message)
}

// Edit message
func EditMessage(userDto *userDTO.UserDTO, request interface{}) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	req, err := parseRequest[EditMessageReq](request)
	if err != nil {
		return nil, nil, err
	}
	req.TrimSpaces()
	if req.MessageContent == "" {
		return nil, nil, appErr.BadRequest("message content is empty")
	}

	chat, requestSendingMember, message, err := verifyAccessAndGetMessage(req.ChatID, userDto.ID, req.MessageID)
	if err != nil {
		return nil, nil, err
	}

	if message.Chat.ID != req.ChatID || message.ChatMember.User.ID != userDto.ID {
		return nil, nil, appErr.Forbidden("you cannot access this message")
	}

	err = message.Edit(req.MessageContent)
	if err != nil {
		return nil, nil, err
	}

	return createResponse(chat, requestSendingMember, message)
}

// Remove message for all
func RemoveMessageForAll(userDto *userDTO.UserDTO, request interface{}) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	req, err := parseRequest[RemoveMessageReq](request)
	if err != nil {
		return nil, nil, err
	}

	chat, requestSendingMember, message, err := verifyAccessAndGetMessage(req.ChatID, userDto.ID, req.MessageID)
	if err != nil {
		return nil, nil, err
	}

	if message.Chat.ID != req.ChatID {
		return nil, nil, appErr.Forbidden("you cannot access this message")
	}

	err = message.RemoveForAll(requestSendingMember)
	if err != nil {
		return nil, nil, err
	}
	message.Content = ""

	return createResponse(chat, requestSendingMember, message)
}
