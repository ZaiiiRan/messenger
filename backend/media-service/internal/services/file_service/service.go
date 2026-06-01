package fileservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	filepb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/file/v1"
	filegrpcclient "github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/client/grpc/file_client"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileService interface {
	GetFiles(ctx context.Context, fileIDs []string) ([]*filepb.File, error)
	DeleteFiles(ctx context.Context, fileIDs []string) error
	UpdateFiles(ctx context.Context, req *filepb.UpdateFilesRequest) ([]*filepb.File, error)
}

type service struct {
	fileClient *filegrpcclient.Client
	log        *zap.SugaredLogger
}

func New(fileClient *filegrpcclient.Client, log *zap.SugaredLogger) FileService {
	return &service{
		fileClient: fileClient,
		log:        log,
	}
}

func (s *service) GetFiles(ctx context.Context, fileIDs []string) ([]*filepb.File, error) {
	l := s.log.With("op", "get_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	if len(fileIDs) == 0 {
		return nil, nil
	}

	resp, err := s.fileClient.FileClient().GetFiles(ctx, &filepb.GetFilesRequest{
		FileIds:                fileIDs,
		IncludeProprietaryInfo: true,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("file.get_files_failed.invalid_argument", "err", err)
		} else if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			l.Errorw("file.get_files_failed.get_files_error", "err", err)
		}
		return nil, err
	}

	l.Infow("file.get_files.success", "file_ids", fileIDs)
	return resp.Files, nil
}

func (s *service) DeleteFiles(ctx context.Context, fileIDs []string) error {
	l := s.log.With("op", "delete_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	if len(fileIDs) == 0 {
		return nil
	}

	_, err := s.fileClient.FileClient().DeleteFiles(ctx, &filepb.DeleteFilesRequest{
		FileIds: fileIDs,
	})
	if err != nil {
		l.Errorw("file.delete_files_failed.delete_files_error", "err", err)
		return err
	}

	l.Infow("file.delete_files.success", "file_ids", fileIDs)
	return nil
}

func (s *service) UpdateFiles(ctx context.Context, req *filepb.UpdateFilesRequest) ([]*filepb.File, error) {
	l := s.log.With("op", "update_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	if len(req.Files) == 0 {
		return nil, nil
	}

	resp, err := s.fileClient.FileClient().UpdateFiles(ctx, req)
	if err != nil {
		l.Errorw("file.update_files_failed.update_files_error", "err", err)
		return nil, err
	}

	l.Infow("file.update_files.success", "updated_files_count", len(resp.Files))
	return resp.Files, nil
}
