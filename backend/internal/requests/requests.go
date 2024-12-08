package requests

import (
	appErr "backend/internal/errors/appError"

	"github.com/gofiber/fiber/v2"
)

// Parse HTTP Request
func ParseRequest[T any](c *fiber.Ctx) (*T, error) {
	var req T
	if err := c.BodyParser(&req); err != nil {
		return nil, appErr.BadRequest("invalid request format")
	}

	if trimmable, ok := any(&req).(interface{ TrimSpaces() }); ok {
		trimmable.TrimSpaces()
	}

	return &req, nil
}

// Parse WebSocket request
func ParseWebSocketRequest[T any](request interface{}) (*T, error) {
	req, ok := request.(*T)
	if !ok || req == nil {
		return nil, appErr.BadRequest("invalid payload")
	}

	if trimmable, ok := any(&req).(interface{ TrimSpaces() }); ok {
		trimmable.TrimSpaces()
	}

	return req, nil
}
