package middleware

import (
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/gofiber/fiber/v2"
)

func UserPermissionMiddleware(shouldProtect PathMatcher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if shouldProtect == nil || !shouldProtect(c.Method(), c.Path()) {
			return c.Next()
		}

		claims, ok := ctxmetadata.GetUserClaimsFromContext(c.UserContext())
		if !ok || claims == nil || claims.IsDeleted {
			return commonerror.ErrUnauthorized
		}

		if claims.IsPermanentlyBanned || claims.IsTemporarilyBanned || !claims.IsConfirmed {
			return commonerror.ErrPermissionDenied
		}

		return c.Next()
	}
}
