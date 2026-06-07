package handlers

import (
	userpb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/httputil"
	"github.com/gofiber/fiber/v2"
)

// GetMe godoc
// @Summary      Get the current user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        include_privacy_settings query bool false "Include privacy settings"
// @Success      200 {object} userpb.GetMeByUserResponse
// @Router       /api/users/me [get]
func (h *Handlers) GetMe(c *fiber.Ctx) error {
	req := &userpb.GetMeByUserRequest{
		IncludePrivacySettings: c.QueryBool("include_privacy_settings", false),
	}

	resp, err := h.clients.User.GetMeByUser(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// UpdateMe godoc
// @Summary      Update the current user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body userpb.UpdateMeByUserRequest true "Fields to update"
// @Success      200 {object} userpb.UpdateMeByUserResponse
// @Router       /api/users/me [patch]
func (h *Handlers) UpdateMe(c *fiber.Ctx) error {
	req := &userpb.UpdateMeByUserRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.User.UpdateMeByUser(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// UpdateMyPrivacySettings godoc
// @Summary      Update the current user privacy settings
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body userpb.UpdateMyPrivacySettingsByUserRequest true "Privacy settings"
// @Success      200 {object} userpb.UpdateMyPrivacySettingsByUserResponse
// @Router       /api/users/me/privacy [patch]
func (h *Handlers) UpdateMyPrivacySettings(c *fiber.Ctx) error {
	req := &userpb.UpdateMyPrivacySettingsByUserRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.User.UpdateMyPrivacySettingsByUser(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// DeleteMe godoc
// @Summary      Delete the current user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} userpb.DeleteMeByUserResponse
// @Router       /api/users/me [delete]
func (h *Handlers) DeleteMe(c *fiber.Ctx) error {
	resp, err := h.clients.User.DeleteMeByUser(h.grpcCtx(c), &userpb.DeleteMeByUserRequest{})
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}
