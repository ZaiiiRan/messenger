package middleware

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func UserAuthMiddleware(secretKey []byte, shouldProtect PathMatcher) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if shouldProtect == nil || !shouldProtect(c.Method(), c.Path()) {
			return c.Next()
		}

		tokenStr, err := extractBearerToken(c)
		if err != nil {
			return commonerror.ErrUnauthorized
		}

		claims, err := jwt.ParseUserToken(tokenStr, secretKey)
		if err != nil {
			return commonerror.ErrUnauthorized
		}

		ctx := ctxmetadata.WithUserClaims(c.UserContext(), claims)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

func extractBearerToken(c *fiber.Ctx) (string, error) {
	authorizationHeader := c.Get("Authorization")
	if len(authorizationHeader) == 0 {
		return "", commonerror.ErrUnauthorized
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", commonerror.ErrUnauthorized
	}

	return parts[1], nil
}
