package wsController

import (
	"backend/internal/requests"
	"encoding/json"
)

// send message request validation
func validateSendMessageRequest(content interface{}) (*requests.SendMessageRequest, bool) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, false
	}

	var req requests.SendMessageRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, false
	}

	if req.ChatID == 0 || req.MessageContent == "" {
		return nil, false
	}

	return &req, true
}

// edit message request validation
func validateEditMessageRequest(content interface{}) (*requests.EditMessageRequest, bool) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, false
	}

	var req requests.EditMessageRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, false
	}

	if req.ChatID == 0 || req.MessageID == 0 || req.MessageContent == "" {
		return nil, false
	}

	return &req, true
}

// remove message request validation
func validateRemoveMessageRequest(content interface{}) (*requests.RemoveMessageRequest, bool) {
	data, err := json.Marshal(content)
	if err != nil {
		return nil, false
	}

	var req requests.RemoveMessageRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, false
	}

	if req.ChatID == 0 || req.MessageID == 0 {
		return nil, false
	}

	return &req, true
}
