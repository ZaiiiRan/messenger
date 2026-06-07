package handlers

import (
	filepb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/file/v1"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/httputil"
	"github.com/gofiber/fiber/v2"
)

// CreateFiles godoc
// @Summary      Create file upload tickets
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body filepb.CreateFilesByUserRequest true "Files to create"
// @Success      200 {object} filepb.CreateFilesByUserResponse
// @Router       /api/files [post]
func (h *Handlers) CreateFiles(c *fiber.Ctx) error {
	req := &filepb.CreateFilesByUserRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.File.CreateFilesByUser(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// ConfirmUploads godoc
// @Summary      Confirm uploaded files
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body filepb.ConfirmUploadsByUserRequest true "File ids to confirm"
// @Success      200 {object} filepb.ConfirmUploadsByUserResponse
// @Router       /api/files/confirm [post]
func (h *Handlers) ConfirmUploads(c *fiber.Ctx) error {
	req := &filepb.ConfirmUploadsByUserRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.File.ConfirmUploadsByUser(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}

// GetFiles godoc
// @Summary      Get files owned by the current user
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body filepb.GetFilesByUserRequest true "File ids"
// @Success      200 {object} filepb.GetFilesByUserResponse
// @Router       /api/files/query [post]
func (h *Handlers) GetFiles(c *fiber.Ctx) error {
	req := &filepb.GetFilesByUserRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.File.GetFilesByUser(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}
