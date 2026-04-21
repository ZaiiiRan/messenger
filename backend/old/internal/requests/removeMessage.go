package requests

type RemoveMessageRequest struct {
	ChatID    uint64 `json:"chat_id"`
	MessageID uint64 `json:"message_id"`
}

// Parse remove message request for web socket
func ParseRemoveMessageReq(request interface{}) (*RemoveMessageRequest, error) {
	return ParseWebSocketRequest[RemoveMessageRequest](request)
}
