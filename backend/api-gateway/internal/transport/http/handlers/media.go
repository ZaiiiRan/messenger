package handlers

import (
	mediapb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/media/v1"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/transport/http/httputil"
	"github.com/gofiber/fiber/v2"
)

// ProcessFiles godoc
// @Summary      Request media processing for uploaded files
// @Tags         media
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body mediapb.ProcessFilesRequest true "Files to process"
// @Success      200 {object} mediapb.ProcessFilesResponse
// @Router       /api/media/process [post]
func (h *Handlers) ProcessFiles(c *fiber.Ctx) error {
	req := &mediapb.ProcessFilesRequest{}
	if err := httputil.Bind(c, req); err != nil {
		return err
	}

	resp, err := h.clients.Media.ProcessFiles(h.grpcCtx(c), req)
	if err != nil {
		return err
	}
	return httputil.Proto(c, resp)
}
