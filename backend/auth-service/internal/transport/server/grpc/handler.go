package grpcserver

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	authservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/auth"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/utils"
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
	return h.authService.Register(ctx, req)
}

func (h *authHandler) GetNewConfirmationCode(ctx context.Context, req *pb.GetNewConfirmationCodeRequest) (*pb.GetNewConfirmationCodeResponse, error) {
	return h.authService.GetNewConfirmationCode(ctx, req)
}

func (h *authHandler) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	utils.SanitizeConfirmRequest(req)
	return h.authService.Confirm(ctx, req)
}

func (h *authHandler) ConfirmByLink(ctx context.Context, req *pb.ConfirmByLinkRequest) (*pb.ConfirmByLinkResponse, error) {
	utils.SanitizeConfirmByLinkRequest(req)
	return h.authService.ConfirmByLink(ctx, req)
}

func (h *authHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	utils.SanitizeLoginRequest(req)
	return h.authService.Login(ctx, req)
}

func (h *authHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	return h.authService.Refresh(ctx, req)
}

func (h *authHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return h.authService.Logout(ctx, req)
}

func (h *authHandler) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	utils.SanitizeChangePasswordRequest(req)
	return h.authService.ChangePassword(ctx, req)
}

func (h *authHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	utils.SanitizeForgotPasswordRequest(req)
	return h.authService.ForgotPassword(ctx, req)
}

func (h *authHandler) ResetPasswordByCode(ctx context.Context, req *pb.ResetPasswordByCodeRequest) (*pb.ResetPasswordByCodeResponse, error) {
	utils.SanitizeResetPasswordByCodeRequest(req)
	return h.authService.ResetPasswordByCode(ctx, req)
}

func (h *authHandler) ResetPasswordByLink(ctx context.Context, req *pb.ResetPasswordByLinkRequest) (*pb.ResetPasswordByLinkResponse, error) {
	utils.SanitizeResetPasswordByLinkRequest(req)
	return h.authService.ResetPasswordByLink(ctx, req)
}

func (h *authHandler) GetActiveSessions(ctx context.Context, req *pb.GetActiveSessionsRequest) (*pb.GetActiveSessionsResponse, error) {
	return h.authService.GetActiveSessions(ctx, req)
}

func (h *authHandler) InvalidateSessions(ctx context.Context, req *pb.InvalidateSessionsRequest) (*pb.InvalidateSessionsResponse, error) {
	return h.authService.InvalidateSessions(ctx, req)
}
