package fileservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	filepb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/file/v1"
	filegrpcclient "github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/client/grpc/file_client"
	"go.uber.org/zap"
)

type FileService interface {
	GetFiles(ctx context.Context, fileIDs []string) ([]*filepb.File, error)
}

type service struct {
	fileClient *filegrpcclient.Client
	log        *zap.SugaredLogger
}

func New(fileClient *filegrpcclient.Client, log *zap.SugaredLogger) FileService {
	return &service{fileClient: fileClient, log: log}
}

func (s *service) GetFiles(ctx context.Context, fileIDs []string) ([]*filepb.File, error) {
	l := s.log.With("op", "file_service.get_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(fileIDs) == 0 {
		return nil, nil
	}

	resp, err := s.fileClient.FileClient().GetFiles(ctx, &filepb.GetFilesRequest{
		FileIds:                fileIDs,
		IncludeProprietaryInfo: false,
	})
	if err != nil {
		l.Warnw("file_service.get_files_failed", "err", err)
		return nil, err
	}

	return resp.Files, nil
}
