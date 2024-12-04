package message

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"time"
)

type Message struct {
	ID         uint64
	Chat       *chat.Chat
	ChatMember *chatMember.ChatMember
	Content    string
	SentAt     time.Time
	LastEdited *time.Time
	IsDeleted  bool
}

// New message
func NewMessage(chat *chat.Chat, member *chatMember.ChatMember, content string) (*Message, error) {
	if content == "" {
		return nil, appErr.BadRequest("message is empty")
	}
	message := &Message{
		Content:    content,
		Chat:       chat,
		ChatMember: member,
	}

	message.Save()
	return message, nil
}

// Save message
func (m *Message) Save() error {
	if m.ID == 0 {
		err := insertMessageToDB(m)
		if err != nil {
			return err
		}
	} else {
		err := updateMessageInDB(m)
		if err != nil {
			return err
		}
	}
	return nil
}
