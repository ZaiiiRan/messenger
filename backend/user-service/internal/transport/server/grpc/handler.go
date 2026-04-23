package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
}

func newUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	utils.SanitizeCreateUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "CreateUser not implemented")
}

func (h *userHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	utils.SanitizeGetUsersRequest(req)
	return nil, status.Error(codes.Unimplemented, "GetUsers not implemented")
}

func (h *userHandler) ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	utils.SanitizeConfirmUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "ConfirmUser not implemented")
}

func (h *userHandler) BanUser(ctx context.Context, req *pb.BanUserRequest) (*pb.BanUserResponse, error) {
	utils.SanitizeBanUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "BanUser not implemented")
}

func (h *userHandler) UnbanUser(ctx context.Context, req *pb.UnbanUserRequest) (*pb.UnbanUserResponse, error) {
	utils.SanitizeUnbanUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "UnbanUser not implemented")
}

func (h *userHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	utils.SanitizeDeleteUserRequest(req)
	return nil, status.Error(codes.Unimplemented, "DeleteUser not implemented")
}
