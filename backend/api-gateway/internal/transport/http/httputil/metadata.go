package httputil

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"
)

// Optional headers that an upstream reverse proxy / CDN may attach with
// geo information about the client.
const (
	countryHeader = "X-Country-Name"
	cityHeader    = "X-City"
)

// GRPCContext builds an outgoing gRPC context carrying all the request-scoped
// data the downstream microservices expect as metadata: the access token
// (Authorization header), the refresh token (cookie), the request id, the
// client IP / user-agent / language and optional geo headers.
func GRPCContext(c *fiber.Ctx, refreshCookieName string) context.Context {
	pairs := make([]string, 0, 16)

	if auth := c.Get(fiber.HeaderAuthorization); auth != "" {
		pairs = append(pairs, ctxmetadata.AuthorizationKey, auth)
	}

	if refresh := c.Cookies(refreshCookieName); refresh != "" {
		pairs = append(pairs, ctxmetadata.RefreshTokenKey, refresh)
	}

	if reqID := ctxmetadata.GetReqIdFromContext(c.UserContext()); reqID != "" {
		pairs = append(pairs, ctxmetadata.RequestIDKey, reqID)
	}

	if ip := c.IP(); ip != "" {
		pairs = append(pairs, ctxmetadata.RealIPKey, ip)
	}

	if ua := c.Get(fiber.HeaderUserAgent); ua != "" {
		pairs = append(pairs, ctxmetadata.UserAgentKey, ua)
	}

	if lang := extractLang(c.Get(fiber.HeaderAcceptLanguage)); lang != "" {
		pairs = append(pairs, ctxmetadata.AcceptLanguageKey, lang)
	}

	if country := c.Get(countryHeader); country != "" {
		pairs = append(pairs, ctxmetadata.CountryNameKey, country)
	}

	if city := c.Get(cityHeader); city != "" {
		pairs = append(pairs, ctxmetadata.CityKey, city)
	}

	md := metadata.Pairs(pairs...)
	return metadata.NewOutgoingContext(c.UserContext(), md)
}
