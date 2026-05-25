package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
	fileservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_service"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/utils"
)

type fileHandler struct {
	pb.UnimplementedFileServiceServer
	fileService fileservice.FileService
}

func newFileHandler(fileService fileservice.FileService) *fileHandler {
	return &fileHandler{
		fileService: fileService,
	}
}

func (h *fileHandler) CreateFilesByUser(ctx context.Context, req *pb.CreateFilesByUserRequest) (*pb.CreateFilesByUserResponse, error) {
	utils.SanitizeCreateFilesByUserRequest(req)
	return h.fileService.CreateFilesByUser(ctx, req)
}

func (h *fileHandler) ConfirmUploadsByUser(ctx context.Context, req *pb.ConfirmUploadsByUserRequest) (*pb.ConfirmUploadsByUserResponse, error) {
	utils.SanitizeConfirmUploadsByUserRequest(req)
	return h.fileService.ConfirmUploadsByUser(ctx, req)
}

func (h *fileHandler) GetFilesByUser(ctx context.Context, req *pb.GetFilesByUserRequest) (*pb.GetFilesByUserResponse, error) {
	utils.SanitizeGetFilesByUserRequest(req)
	return h.fileService.GetFilesByUser(ctx, req)
}

func (h *fileHandler) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	utils.SanitizeGetFilesRequest(req)
	return h.fileService.GetFiles(ctx, req)
}

func (h *fileHandler) DeleteFiles(ctx context.Context, req *pb.DeleteFilesRequest) (*pb.DeleteFilesResponse, error) {
	utils.SanitizeDeleteFilesRequest(req)
	return h.fileService.DeleteFiles(ctx, req)
}

func (h *fileHandler) UpdateFiles(ctx context.Context, req *pb.UpdateFilesRequest) (*pb.UpdateFilesResponse, error) {
	utils.SanitizeUpdateFilesRequest(req)
	return h.fileService.UpdateFiles(ctx, req)
}
