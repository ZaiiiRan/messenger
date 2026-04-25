package grpcserver

import (
	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	authservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/auth"
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
