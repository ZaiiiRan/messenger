package handlers

import (
	authpb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/auth/v1"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/httputil"
	"github.com/gofiber/fiber/v2"
)

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authpb.RegisterRequest true "Registration data"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/register [post]
func (h *Handlers) Register(c *fiber.Ctx) error {
	req := &authpb.RegisterRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.Register(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authpb.LoginRequest true "Login credentials"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/login [post]
func (h *Handlers) Login(c *fiber.Ctx) error {
	req := &authpb.LoginRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.Login(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// Refresh godoc
// @Summary      Refresh access and refresh tokens
// @Tags         auth
// @Produce      json
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/refresh [post]
func (h *Handlers) Refresh(c *fiber.Ctx) error {
	resp, err := h.clients.Auth.Refresh(h.grpcCtx(c), &authpb.RefreshRequest{})
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// Logout godoc
// @Summary      Logout
// @Tags         auth
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /api/auth/logout [post]
func (h *Handlers) Logout(c *fiber.Ctx) error {
	_, err := h.clients.Auth.Logout(h.grpcCtx(c), &authpb.LogoutRequest{})
	if err != nil {
		return err
	}
	httputil.ClearRefreshTokenCookie(c, h.cfg.Auth)
	return httputil.JSON(c, fiber.Map{"message": "ok"})
}

// GetNewConfirmationCode godoc
// @Summary      Request a new account confirmation code
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} authpb.GetNewConfirmationCodeResponse
// @Router       /api/auth/confirmation-code [post]
func (h *Handlers) GetNewConfirmationCode(c *fiber.Ctx) error {
	resp, err := h.clients.Auth.GetNewConfirmationCode(h.grpcCtx(c), &authpb.GetNewConfirmationCodeRequest{})
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// ConfirmByCode godoc
// @Summary      Confirm account by code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body authpb.ConfirmByCodeRequest true "Confirmation code"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/confirm/code [post]
func (h *Handlers) ConfirmByCode(c *fiber.Ctx) error {
	req := &authpb.ConfirmByCodeRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ConfirmByCode(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// ConfirmByLink godoc
// @Summary      Confirm account by link token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authpb.ConfirmByLinkRequest true "Confirmation token"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/confirm/link [post]
func (h *Handlers) ConfirmByLink(c *fiber.Ctx) error {
	req := &authpb.ConfirmByLinkRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ConfirmByLink(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// ChangePassword godoc
// @Summary      Change password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body authpb.ChangePasswordRequest true "Old and new password"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/password/change [post]
func (h *Handlers) ChangePassword(c *fiber.Ctx) error {
	req := &authpb.ChangePasswordRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ChangePassword(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// ForgotPassword godoc
// @Summary      Request a password reset
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authpb.ForgotPasswordRequest true "Login"
// @Success      200 {object} authpb.ForgotPasswordResponse
// @Router       /api/auth/password/forgot [post]
func (h *Handlers) ForgotPassword(c *fiber.Ctx) error {
	req := &authpb.ForgotPasswordRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ForgotPassword(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// ResetPasswordByCode godoc
// @Summary      Reset password by code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authpb.ResetPasswordByCodeRequest true "Reset data"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/password/reset/code [post]
func (h *Handlers) ResetPasswordByCode(c *fiber.Ctx) error {
	req := &authpb.ResetPasswordByCodeRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ResetPasswordByCode(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// ResetPasswordByLink godoc
// @Summary      Reset password by link token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body authpb.ResetPasswordByLinkRequest true "Reset data"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/password/reset/link [post]
func (h *Handlers) ResetPasswordByLink(c *fiber.Ctx) error {
	req := &authpb.ResetPasswordByLinkRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ResetPasswordByLink(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// ChangeEmail godoc
// @Summary      Request an email change
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body authpb.ChangeEmailRequest true "New email"
// @Success      200 {object} authpb.ChangeEmailResponse
// @Router       /api/auth/email/change [post]
func (h *Handlers) ChangeEmail(c *fiber.Ctx) error {
	req := &authpb.ChangeEmailRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ChangeEmail(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// ConfirmNewEmailByCode godoc
// @Summary      Confirm a new email by code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body authpb.ConfirmNewEmailByCodeRequest true "Confirmation code"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/email/confirm/code [post]
func (h *Handlers) ConfirmNewEmailByCode(c *fiber.Ctx) error {
	req := &authpb.ConfirmNewEmailByCodeRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ConfirmNewEmailByCode(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// ConfirmNewEmailByLink godoc
// @Summary      Confirm a new email by link token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body authpb.ConfirmNewEmailByLinkRequest true "Confirmation token"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/email/confirm/link [post]
func (h *Handlers) ConfirmNewEmailByLink(c *fiber.Ctx) error {
	req := &authpb.ConfirmNewEmailByLinkRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.ConfirmNewEmailByLink(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
}

// GetActiveSessions godoc
// @Summary      List active sessions
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page"
// @Param        page_size query int false "Page size"
// @Success      200 {object} authpb.GetActiveSessionsResponse
// @Router       /api/auth/sessions [get]
func (h *Handlers) GetActiveSessions(c *fiber.Ctx) error {
	req := &authpb.GetActiveSessionsRequest{
		Page:     int32(c.QueryInt("page", 0)),
		PageSize: int32(c.QueryInt("page_size", 0)),
	}

	resp, err := h.clients.Auth.GetActiveSessions(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// InvalidateSessions godoc
// @Summary      Invalidate sessions
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body authpb.InvalidateSessionsRequest true "Session ids to invalidate"
// @Success      200 {object} handlers.authPayload
// @Router       /api/auth/sessions/invalidate [post]
func (h *Handlers) InvalidateSessions(c *fiber.Ctx) error {
	req := &authpb.InvalidateSessionsRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Auth.InvalidateSessions(h.grpcCtx(c), req)
	if err != nil {
		return err
	}

	// When the current session is invalidated the service rotates the tokens;
	// otherwise no token fields are returned.
	if resp.GetAccessToken() != "" || resp.GetRefreshToken() != "" {
		return h.respondWithTokens(c, resp.GetUser(), resp.GetAccessToken(), resp.GetRefreshToken())
	}
	return httputil.JSON(c, fiber.Map{"message": "ok"})
}
