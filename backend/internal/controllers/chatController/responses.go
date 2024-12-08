package chatController

import (
	"backend/internal/models/chat/chatDTO"
	"backend/internal/models/chat/chatMember/chatMemberDTO"
	"backend/internal/models/chat/message/messageDTO"
)

type ChatResponse struct {
	Chat        *chatDTO.ChatDTO               `json:"chat"`
	Members     []*chatMemberDTO.ChatMemberDTO `json:"members"`
	You         *chatMemberDTO.ChatMemberDTO   `json:"you"`
	LastMessage *messageDTO.MessageDTO         `json:"last_message"`
}
