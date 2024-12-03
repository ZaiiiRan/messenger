package webSocketManager

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type WebSocketManager struct {
	clients map[uint64]map[*websocket.Conn]bool
	mu      sync.Mutex
}

// Add ws connection
func (wm *WebSocketManager) AddConnection(key uint64, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if wm.clients[key] == nil {
		wm.clients[key] = make(map[*websocket.Conn]bool)
	}
	wm.clients[key][conn] = true
}

// Remove ws connection
func (wm *WebSocketManager) RemoveConnection(key uint64, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if _, ok := wm.clients[key]; ok {
		delete(wm.clients[key], conn)
		if len(wm.clients[key]) == 0 {
			delete(wm.clients, key)
		}
	}
}

// Broadcast to client
func (wm *WebSocketManager) BroadcastToClient(key uint64, message []byte) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for conn := range wm.clients[key] {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			conn.Close()
			delete(wm.clients[key], conn)
		}
	}
}

var instance *WebSocketManager
var once sync.Once

// Get instance of WS Manager
func GetInstance() *WebSocketManager {
	once.Do(func() {
		instance = &WebSocketManager{
			clients: make(map[uint64]map[*websocket.Conn]bool),
		}
	})
	return instance
}
