package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	userService userservice.UserService
}

func newUserHandler(userService userservice.UserService) *userHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	utils.SanitizeCreateUserRequest(req)
	return h.userService.CreateUser(ctx, req)
}

func (h *userHandler) ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	utils.SanitizeConfirmUserRequest(req)
	return h.userService.ConfirmUser(ctx, req)
}

func (h *userHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	utils.SanitizeGetUserByIDRequest(req)
	return h.userService.GetUserByID(ctx, req)
}

func (h *userHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	utils.SanitizeGetUserByUsernameRequest(req)
	return h.userService.GetUserByUsername(ctx, req)
}

func (h *userHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	utils.SanitizeGetUserByEmailRequest(req)
	return h.userService.GetUserByEmail(ctx, req)
}

func (h *userHandler) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	utils.SanitizeGetUsersRequest(req)
	return h.userService.GetUsers(ctx, req)
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
