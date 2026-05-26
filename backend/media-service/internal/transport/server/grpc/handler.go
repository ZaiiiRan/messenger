package grpcserver

import (
	pb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/media/v1"
)

type mediaHandler struct {
	pb.UnimplementedMediaServiceServer
}

func newMediaHandler() *mediaHandler {
	return &mediaHandler{}
}
