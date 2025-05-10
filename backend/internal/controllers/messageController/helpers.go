package messageController

import (
	appErr "backend/internal/errors/appError"
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message"
	"backend/internal/models/chat/message/messageDTO"
	"backend/internal/models/socialUser"
	"errors"
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

	var appError *appErr.AppError
	members, err := chat.GetChatMembers(requestSendingMember)
	if err != nil && errors.As(err, &appError) && appError.StatusCode != 404 {
		return nil, nil, err
	}
	membersDTOs := chatMemberDTO.CreateChatMembersDTOs(members)

	return messageDto, membersDTOs, nil
}

// check user access
func checkUserAccess(requestSenderID, targetID uint64) error {
	target, err := socialUser.GetTargetByID(requestSenderID, targetID)
	if err != nil {
		return err
	}

	if target.User.IsBanned {
		return appErr.BadRequest("your interlocutor is blocked")
	}
	if !target.User.IsActivated || target.User.IsDeleted {
		return appErr.NotFound("user not found")
	}

	if target.FriendStatus != nil && *target.FriendStatus == "blocked by target" {
		return appErr.BadRequest("you are blocked by your interlocutor")
	}

	if target.FriendStatus != nil && *target.FriendStatus == "blockedr" {
		return appErr.BadRequest("your interlocutor is blocked by you")
	}

	return nil
}
