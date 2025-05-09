package messageController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
	"backend/internal/models/user/userDTO"
	"backend/internal/requests"
)

// Send message
func SendMessage(userDto *userDTO.UserDTO, request interface{}) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	req, err := requests.ParseSendMessageRequest(request)
	if err != nil {
		return nil, nil, err
	}
	if req.MessageContent == "" {
		return nil, nil, appErr.BadRequest("message content is empty")
	}

	chat, requestSendingMember, err := chatModel.GetChatAndVerifyAccess(req.ChatID, userDto.ID)
	if err != nil {
		return nil, nil, err
	}

	if chat.IsGroupChat {
		sendMessageToPrivateChat(chat, requestSendingMember, req)
	}

	message, err := message.NewMessage(chat, requestSendingMember, req.MessageContent)
	if err != nil {
		return nil, nil, err
	}

	return createResponse(chat, requestSendingMember, message)
}

// check if user can send message to private chat
func sendMessageToPrivateChat(chat *chatModel.Chat, requestSendingMember *chatMember.ChatMember,
	req *requests.SendMessageRequest) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	members, err := chat.GetChatMembers(requestSendingMember)
	if err != nil {
		return nil, nil, err
	}

	err = checkUserAccess(requestSendingMember.User.ID, members[0].User.ID)
	if err != nil {
		return nil, nil, err
	}

	message, err := message.NewMessage(chat, requestSendingMember, req.MessageContent)
	if err != nil {
		return nil, nil, err
	}

	messageDto := messageDTO.CreateMessageDTO(message)
	membersDTOs := chatMemberDTO.CreateChatMembersDTOs(members)
	return messageDto, membersDTOs, nil
}

// Edit message
func EditMessage(userDto *userDTO.UserDTO, request interface{}) (*messageDTO.MessageDTO, []*chatMemberDTO.ChatMemberDTO, error) {
	req, err := requests.ParseEditMessageReq(request)
	if err != nil {
		return nil, nil, err
	}
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
	req, err := requests.ParseRemoveMessageReq(request)
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
