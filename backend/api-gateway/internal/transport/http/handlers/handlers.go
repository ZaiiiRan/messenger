package handlers

import (
	"context"
	"encoding/json"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config"
	grpcclient "github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/client/grpc"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/httputil"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/proto"
)

// Handlers wires the upstream gRPC clients and gateway configuration into the
// HTTP request handlers.
type Handlers struct {
	clients *grpcclient.Clients
	cfg     *config.ServerConfig
}

func New(clients *grpcclient.Clients, cfg *config.ServerConfig) *Handlers {
	return &Handlers{
		clients: clients,
		cfg:     cfg,
	}
}

// grpcCtx builds the outgoing gRPC context carrying auth, cookie and request
// metadata extracted from the incoming HTTP request.
func (h *Handlers) grpcCtx(c *fiber.Ctx) context.Context {
	return httputil.GRPCContext(c, h.cfg.Auth.RefreshTokenCookieName)
}

// authPayload is the response shape returned for every auth flow that issues
// tokens, mirroring the legacy backend: the user object plus the access token,
// while the refresh token travels in an HttpOnly cookie.
type authPayload struct {
	User        json.RawMessage `json:"user"`
	AccessToken string          `json:"accessToken,omitempty"`
}

// respondWithTokens sets the refresh-token cookie and returns the user object
// together with the access token as JSON.
func (h *Handlers) respondWithTokens(c *fiber.Ctx, user proto.Message, accessToken, refreshToken string) error {
	httputil.SetRefreshTokenCookie(c, h.cfg.Auth, refreshToken)

	rawUser, err := httputil.RawProto(user)
	if err != nil {
		return err
	}

	return httputil.JSON(c, authPayload{User: rawUser, AccessToken: accessToken})
}
