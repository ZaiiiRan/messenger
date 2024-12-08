package chatDTO

import (
	"backend/internal/models/chat"
)

type ChatDTO struct {
	ID          uint64  `json:"id"`
	Name        *string `json:"name"`
	IsGroupChat bool    `json:"is_group_chat"`
}

// Create Chat DTO from chat object
func CreateChatDTO(chat *chat.Chat) *ChatDTO {
	return &ChatDTO{
		ID:          chat.ID,
		Name:        chat.Name,
		IsGroupChat: chat.IsGroupChat,
	}
}
