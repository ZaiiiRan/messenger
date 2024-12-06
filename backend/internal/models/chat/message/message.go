package message

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/models/chat"
	"backend/internal/models/chat/chatMember"
	"backend/internal/utils"
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
	if err := validateContent(content); err != nil {
		return nil, err
	}
	message := &Message{
		Content:    content,
		Chat:       chat,
		ChatMember: member,
	}

	err := message.Save()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Edit message
func (m *Message) Edit(content string) error {
	if m.SentAt.Before(time.Now().Add(-12 * time.Hour)) {
		return appErr.BadRequest("the message was sent more than 12 hours ago")
	}
	if err := validateContent(content); err != nil {
		return err
	}
	if m.Content == content {
		return appErr.BadRequest("content match")
	}
	m.Content = content
	m.LastEdited = utils.TimePtr(time.Now())

	err := m.Save()
	if err != nil {
		return err
	}
	return nil
}

// Remove message
func (m *Message) RemoveForAll(actor *chatMember.ChatMember) error {
	if m.IsDeleted {
		return appErr.BadRequest("the message has already been deleted")
	}
	
	if !m.Chat.IsGroupChat {
		if m.ChatMember.User.ID != actor.User.ID {
			return appErr.BadRequest("you cannot access this message")
		}
	} else if actor.Role > m.ChatMember.Role {
		m.IsDeleted = true
	} else if m.ChatMember.User.ID != actor.User.ID {
		return appErr.Forbidden("you cannot access this message")
	} else {
		if m.SentAt.Before(time.Now().Add(-12 * time.Hour)) {
			return appErr.BadRequest("the message was sent more than 12 hours ago")
		}
	}

	m.IsDeleted = true

	err := m.Save()
	if err != nil {
		return err
	}
	return nil
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

// Get message by ID
func GetMessage(id uint64) (*Message, error) {
	return getMessageByIDFromDB(id)
}

// Get messages
func GetMessages(chat *chat.Chat, actor *chatMember.ChatMember, limit, offset int) ([]Message, error) {
	return getMessagesFromDB(chat, actor, limit, offset)
}

// validate message content
func validateContent(content string) error {
	if content == "" {
		return appErr.BadRequest("message content is empty")
	}
	return nil
}
