package requests

import (
	"strings"
)

type SendMessageRequest struct {
	ChatID         uint64 `json:"chat_id"`
	MessageContent string `json:"message_content"`
}

// Send message request trim spaces
func (r *SendMessageRequest) TrimSpaces() {
	r.MessageContent = strings.TrimSpace(r.MessageContent)
}

// Parse send message request for web socket
func ParseSendMessageRequest(request interface{}) (*SendMessageRequest, error) {
	return ParseWebSocketRequest[SendMessageRequest](request)
}
