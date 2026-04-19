package requests

import (
	"strings"
)

type EditMessageRequest struct {
	SendMessageRequest
	MessageID uint64 `json:"message_id"`
}

// Edit message request trim spaces
func (r *EditMessageRequest) TrimSpaces() {
	r.MessageContent = strings.TrimSpace(r.MessageContent)
}

// Parse edit message request for web socket
func ParseEditMessageReq(request interface{}) (*EditMessageRequest, error) {
	return ParseWebSocketRequest[EditMessageRequest](request)
}
