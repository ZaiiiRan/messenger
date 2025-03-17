package wsController

import "strings"

// Auth WebSocket message format
type AuthWebSocketMessage struct {
	Token string `json:"token"`
}

// Trim Spaces for Auth WebSocket message
func (m *AuthWebSocketMessage) TrimSpaces() {
	m.Token = strings.TrimSpace(m.Token)
}
