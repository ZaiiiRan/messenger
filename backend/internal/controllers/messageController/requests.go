package messageController

import (
	"strings"
)

// Send message request format
type SendMessageReq struct {
	ChatID         uint64 `json:"chat_id"`
	MessageContent string `json:"message_content"`
}

// Send message request trim spaces
func (r *SendMessageReq) TrimSpaces() {
	r.MessageContent = strings.TrimSpace(r.MessageContent)
}

// Edit message request format
type EditMessageReq struct {
	ChatID         uint64 `json:"chat_id"`
	MessageID      uint64 `json:"message_id"`
	MessageContent string `json:"message_content"`
}

// Edit message request trim spaces
func (r *EditMessageReq) TrimSpaces() {
	r.MessageContent = strings.TrimSpace(r.MessageContent)
}

// Remove message request format
type RemoveMessageReq struct {
	ChatID    uint64 `json:"chat_id"`
	MessageID uint64 `json:"message_id"`
}
