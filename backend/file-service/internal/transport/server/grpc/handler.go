package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fileHandler struct {
	pb.UnimplementedFileServiceServer
}

func newFileHandler() *fileHandler {
	return &fileHandler{}
}

func (h *fileHandler) CreateFilesByUser(ctx context.Context, req *pb.CreateFilesByUserRequest) (*pb.CreateFilesByUserResponse, error) {
	utils.SanitizeCreateFilesByUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "not_implemented")
}

func (h *fileHandler) ConfirmUploadsByUser(ctx context.Context, req *pb.ConfirmUploadsByUserRequest) (*pb.ConfirmUploadsByUserResponse, error) {
	utils.SanitizeConfirmUploadsByUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "not_implemented")
}

func (h *fileHandler) GetFilesByUser(ctx context.Context, req *pb.GetFilesByUserRequest) (*pb.GetFilesByUserResponse, error) {
	utils.SanitizeGetFilesByUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "not_implemented")
}

func (h *fileHandler) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	utils.SanitizeGetFilesRequest(req)
	return nil, status.Error(codes.Unimplemented, "not_implemented")
}

func (h *fileHandler) DeleteFiles(ctx context.Context, req *pb.DeleteFilesRequest) (*pb.DeleteFilesResponse, error) {
	utils.SanitizeDeleteFilesRequest(req)
	return nil, status.Error(codes.Unimplemented, "not_implemented")
}

func (h *fileHandler) UpdateFiles(ctx context.Context, req *pb.UpdateFilesRequest) (*pb.UpdateFilesResponse, error) {
	utils.SanitizeUpdateFilesRequest(req)
	return nil, status.Error(codes.Unimplemented, "not_implemented")
}
