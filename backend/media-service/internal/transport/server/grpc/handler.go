package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/media/v1"
	mediaservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/media_service"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/utils"
)

type mediaHandler struct {
	pb.UnimplementedMediaServiceServer

	mediaService mediaservice.MediaService
}

func newMediaHandler(mediaService mediaservice.MediaService) *mediaHandler {
	return &mediaHandler{
		mediaService: mediaService,
	}
}

func (h *mediaHandler) ProcessFiles(ctx context.Context, req *pb.ProcessFilesRequest) (*pb.ProcessFilesResponse, error) {
	utils.SanitizeProcessFilesRequest(req)
	return h.mediaService.ProcessFiles(ctx, req)
}
