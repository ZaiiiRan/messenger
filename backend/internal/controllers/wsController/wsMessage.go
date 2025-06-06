package wsController

import "strings"

// WebSocket message format
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// Trim Spaces for WS Message
func (m *WebSocketMessage) TrimSpaces() {
	m.Type = strings.TrimSpace(m.Type)
}
