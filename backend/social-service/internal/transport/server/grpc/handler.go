package grpcserver

import (
	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
)

type socialHandler struct {
	pb.UnimplementedSocialServiceServer
}

func newSocialHandler() *socialHandler {
	return &socialHandler{}
}
