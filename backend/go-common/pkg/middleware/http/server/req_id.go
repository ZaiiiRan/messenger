package middleware

import (
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestIdMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(ctxmetadata.RequestIDKey)
		if id == "" {
			id = uuid.NewString()
		}
		ctx := ctxmetadata.WithReqId(c.UserContext(), id)
		c.SetUserContext(ctx)
		c.Set("X-Request-ID", id)
		return c.Next()
	}
}
