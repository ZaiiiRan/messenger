package middleware

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func I18nMiddleware(getLocalizerFunc func(lang string) *i18n.Localizer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := extractLang(c.Get("Accept-Language"))
		localizer := getLocalizerFunc(lang)
		ctx := ctxmetadata.WithLocalizer(c.UserContext(), localizer)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func extractLang(acceptLang string) string {
	s := strings.SplitN(acceptLang, ",", 2)[0]
	s = strings.TrimSpace(strings.SplitN(s, ";", 2)[0])
	if idx := strings.IndexByte(s, '-'); idx != -1 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}
