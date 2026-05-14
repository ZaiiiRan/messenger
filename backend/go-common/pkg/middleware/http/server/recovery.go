package middleware

import (
	"runtime/debug"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RecoveryMiddleware(log *zap.SugaredLogger) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorw(
					"http.panic",
					"method", c.Method(),
					"path", c.Path(),
					"panic", r,
					"stack", string(debug.Stack()),
				)
				err = commonerror.ErrInternal
			}
		}()
		return c.Next()
	}
}
