package middleware

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func LogMiddleware(log *zap.SugaredLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqId := ctxmetadata.GetReqIdFromContext(c.UserContext())
		start := time.Now()
		err := c.Next()
		log.Infow(
			"http.request",
			"req_id", reqId,
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}
}
