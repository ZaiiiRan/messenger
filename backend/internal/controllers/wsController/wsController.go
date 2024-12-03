package wsController

import (
	"backend/internal/models/user/userDTO"
	"backend/internal/webSocketManager"
	appErr "backend/internal/errors/appError"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func InitConnection(c *fiber.Ctx) error {
	user, ok := c.Locals("userDTO").(*userDTO.UserDTO)
	if !ok || user == nil {
		return appErr.Unauthorized("unauthorized")
	}
	if websocket.IsWebSocketUpgrade(c) {
		return websocket.New(func(conn *websocket.Conn) {
			HandleWebSocket(conn, user)
		})(c)
	}
	return fiber.ErrUpgradeRequired
} 

func HandleWebSocket(conn *websocket.Conn, userDTO *userDTO.UserDTO) {
	manager := webSocketManager.GetInstance()
	manager.AddConnection(userDTO.ID, conn)
	defer manager.RemoveConnection(userDTO.ID, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		manager.BroadcastToClient(userDTO.ID, msg)
	}
}