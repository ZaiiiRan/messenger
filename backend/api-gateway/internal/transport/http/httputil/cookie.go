package httputil

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config/settings"
	"github.com/gofiber/fiber/v2"
)

// SetRefreshTokenCookie stores the refresh token in an HttpOnly cookie, exactly
// like the legacy backend did. A no-op when token is empty.
func SetRefreshTokenCookie(c *fiber.Ctx, cfg settings.AuthSettings, token string) {
	if token == "" {
		return
	}
	c.Cookie(&fiber.Cookie{
		Name:     cfg.RefreshTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(time.Duration(cfg.RefreshTokenTTLHours) * time.Hour),
		HTTPOnly: true,
		Path:     cfg.CookiePath,
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		SameSite: cfg.CookieSameSite,
	})
}

// ClearRefreshTokenCookie removes the refresh token cookie.
func ClearRefreshTokenCookie(c *fiber.Ctx, cfg settings.AuthSettings) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.RefreshTokenCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     cfg.CookiePath,
		Domain:   cfg.CookieDomain,
		Secure:   cfg.CookieSecure,
		SameSite: cfg.CookieSameSite,
	})
}
