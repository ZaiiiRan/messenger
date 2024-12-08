package chatController

import (
	chatModel "backend/internal/models/chat"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message/messageDTO"
)

type ChatResponse struct {
	Chat        *chatModel.Chat                `json:"chat"`
	Members     []*chatMemberDTO.ChatMemberDTO `json:"members"`
	You         *chatMemberDTO.ChatMemberDTO   `json:"you"`
	LastMessage *messageDTO.MessageDTO         `json:"last_message"`
}
