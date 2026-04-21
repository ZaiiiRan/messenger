package grpcserver

import (
	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
}

func newUserHandler() *userHandler {
	return &userHandler{}
}
