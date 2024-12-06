package wsController

import (
	"backend/internal/controllers/messageController"
	"encoding/json"
)

// send message request validation
func validateSendMessageRequest(content interface{}) (*messageController.SendMessageReq, bool) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, false
	}

	var req messageController.SendMessageReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, false
	}

	if req.ChatID == 0 || req.MessageContent == "" {
		return nil, false
	}

	return &req, true
}

// edit message request validation
func validateEditMessageRequest(content interface{}) (*messageController.EditMessageReq, bool) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, false
	}

	var req messageController.EditMessageReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, false
	}

	if req.ChatID == 0 || req.MessageID == 0 || req.MessageContent == "" {
		return nil, false
	}

	return &req, true
}

// remove message request validation
func validateRemoveMessageRequest(content interface{}) (*messageController.RemoveMessageReq, bool) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, false
	}

	var req messageController.RemoveMessageReq
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, false
	}

	if req.ChatID == 0 || req.MessageID == 0 {
		return nil, false
	}

	return &req, true
}
