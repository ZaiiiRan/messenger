package grpcserver

import (
	pb "github.com/ZaiiiRan/messenger/backend/file-service/gen/go/file/v1"
)

type fileHandler struct {
	pb.UnimplementedFileServiceServer
}

func newFileHandler() *fileHandler {
	return &fileHandler{}
}
