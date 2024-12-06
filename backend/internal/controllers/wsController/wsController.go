package wsController

import (
	"backend/internal/controllers/messageController"
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user/userDTO"
	"backend/internal/webSocketManager"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Init WebSocket connection
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

// Handle WebSocket
func HandleWebSocket(conn *websocket.Conn, userDTO *userDTO.UserDTO) {
	manager := webSocketManager.GetInstance()
	manager.AddConnection(userDTO.ID, conn)
	defer manager.RemoveConnection(userDTO.ID, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMessage WebSocketMessage
		if err := json.Unmarshal(msg, &wsMessage); err != nil {
			sendError(conn, appErr.BadRequest("invalid request format"))
			continue
		}
		wsMessage.TrimSpaces()

		err = processRequest(conn, manager, userDTO, &wsMessage)
		if err != nil {
			sendError(conn, err)
		}
	}
}

// requests processing
func processRequest(conn *websocket.Conn, manager *webSocketManager.WebSocketManager, userDTO *userDTO.UserDTO, wsMessage *WebSocketMessage) error {
	switch wsMessage.Type {
	case "send_message":
		return handleSendMessage(manager, userDTO, wsMessage.Content)
	case "edit_message":
		return handleEditMessage(manager, userDTO, wsMessage.Content)
	case "remove_message_for_all":
		return handleRemoveMessageForAll(manager, userDTO, wsMessage.Content)
	case "ping":
		return sendPong(conn)
	default:
		return appErr.BadRequest("unknown request type")
	}
}

// handle send message
func handleSendMessage(manager *webSocketManager.WebSocketManager, userDTO *userDTO.UserDTO, content interface{}) error {
	req, ok := validateSendMessageRequest(content)
	if !ok {
		return appErr.BadRequest("invalid send_message payload")
	}

	message, members, err := messageController.SendMessage(userDTO, req)
	if err != nil {
		return err
	}

	wsMsg := WebSocketMessage{
		Type:    "new_message_notification",
		Content: message,
	}
	wsMsgJSON, err := json.Marshal(wsMsg)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}

	for _, member := range members {
		manager.BroadcastToClient(member.User.ID, wsMsgJSON)
	}
	manager.BroadcastToClient(userDTO.ID, wsMsgJSON)
	return nil
}

// handle edit message
func handleEditMessage(manager *webSocketManager.WebSocketManager, userDTO *userDTO.UserDTO, content interface{}) error {
	req, ok := validateEditMessageRequest(content)
	if !ok {
		return appErr.BadRequest("invalid edit_message payload")
	}

	message, members, err := messageController.EditMessage(userDTO, req)
	if err != nil {
		return err
	}

	wsMsg := WebSocketMessage{
		Type:    "edit_message_notification",
		Content: message,
	}
	wsMsgJSON, err := json.Marshal(wsMsg)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}

	for _, member := range members {
		manager.BroadcastToClient(member.User.ID, wsMsgJSON)
	}
	manager.BroadcastToClient(userDTO.ID, wsMsgJSON)
	return nil
}

// handle remove message for all members
func handleRemoveMessageForAll(manager *webSocketManager.WebSocketManager, userDTO *userDTO.UserDTO, content interface{}) error {
	req, ok := validateRemoveMessageRequest(content)
	if !ok {
		return appErr.BadRequest("invalid remove_message payload")
	}

	message, members, err := messageController.RemoveMessageForAll(userDTO, req)
	if err != nil {
		return err
	}

	wsMsg := WebSocketMessage{
		Type:    "remove_message_notification",
		Content: message,
	}
	wsMsgJSON, err := json.Marshal(wsMsg)
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}

	for _, member := range members {
		manager.BroadcastToClient(member.User.ID, wsMsgJSON)
	}
	manager.BroadcastToClient(userDTO.ID, wsMsgJSON)
	return nil
}

// send pong to self websocket
func sendPong(conn *websocket.Conn) error {
	response := WebSocketMessage{
		Type:    "pong",
		Content: "pong",
	}
	return sendResponse(conn, response)
}

// send error to self websocket
func sendError(conn *websocket.Conn, err error) {
	response := WebSocketMessage{
		Type:    "error",
		Content: err.Error(),
	}
	sendResponse(conn, response)
}

// send response to self websocket
func sendResponse(conn *websocket.Conn, response WebSocketMessage) error {
	respJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, respJSON)
}
