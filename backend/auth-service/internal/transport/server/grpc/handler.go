package grpcserver

import (
	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
}

func newAuthHandler() *authHandler {
	return &authHandler{}
}
