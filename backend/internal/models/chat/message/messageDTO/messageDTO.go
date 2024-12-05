package messageDTO

import (
	"backend/internal/models/chat/message"
	"time"
)

type MessageDTO struct {
	ID         uint64     `json:"id"`
	ChatID     uint64     `json:"chat_id"`
	MemberID   uint64     `json:"member_id"`
	Content    string     `json:"content"`
	SentAt     time.Time  `json:"sent_at"`
	LastEdited *time.Time `json:"last_edited"`
}

func NewMessageDTO(message *message.Message) *MessageDTO {
	return &MessageDTO{
		ID:         message.ID,
		ChatID:     message.Chat.ID,
		MemberID:   message.ChatMember.User.ID,
		Content:    message.Content,
		SentAt:     message.SentAt,
		LastEdited: message.LastEdited,
	}
}

func GetMessagesDTOs(messages []message.Message) []*MessageDTO {
	messagessDTOs := make([]*MessageDTO, len(messages))
	for index, message := range messages {
		messagessDTOs[index] = NewMessageDTO(&message)
	}
	return messagessDTOs
}
