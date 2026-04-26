package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	authservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/auth"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	authService authservice.AuthService
}

func newAuthHandler(authService authservice.AuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	utils.SanitizeRegisterRequest(req)
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

func (h *authHandler) GetNewConfirmationCode(ctx context.Context, req *pb.GetNewConfirmationCodeRequest) (*pb.GetNewConfirmationCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

func (h *authHandler) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	utils.SanitizeConfirmRequest(req)
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

func (h *authHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	utils.SanitizeLoginRequest(req)
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

func (h *authHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

func (h *authHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}

func (h *authHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	utils.SanitizeChangePasswordRequest(req)
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
