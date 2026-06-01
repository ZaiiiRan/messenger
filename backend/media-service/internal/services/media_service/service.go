package mediaservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	commonerror "github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	filepb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/file/v1"
	mediapb "github.com/ZaiiiRan/messenger/backend/media-service/gen/go/media/v1"
	fileprocessingtasksservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_tasks_service"
	fileservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_service"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type MediaService interface {
	ProcessFiles(ctx context.Context, req *mediapb.ProcessFilesRequest) (*mediapb.ProcessFilesResponse, error)
}

type service struct {
	log          *zap.SugaredLogger
	fileService  fileservice.FileService
	tasksService fileprocessingtasksservice.FileProcessingTasksService
}

func New(
	fileService fileservice.FileService,
	tasksService fileprocessingtasksservice.FileProcessingTasksService,
	log *zap.SugaredLogger,
) MediaService {
	return &service{
		log:          log,
		fileService:  fileService,
		tasksService: tasksService,
	}
}

func (s *service) ProcessFiles(ctx context.Context, req *mediapb.ProcessFilesRequest) (*mediapb.ProcessFilesResponse, error) {
	l := s.log.With("op", "process_files", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	claims, ok := ctxmetadata.GetUserClaimsFromContext(ctx)
	if !ok || claims == nil {
		return nil, grpcstatus.Error(codes.Unauthenticated, commonerror.ErrUnauthorized.Error())
	}

	if req == nil || len(req.Files) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrNoFiles.Error())
	}
	if len(req.Files) > maxFilesPerRequest {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyFiles.Error())
	}

	fileIDs := make([]string, 0, len(req.Files))
	seen := make(map[string]struct{}, len(req.Files))
	for _, fr := range req.Files {
		id := fr.GetFileId()
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		fileIDs = append(fileIDs, id)
	}

	files, err := s.fileService.GetFiles(ctx, fileIDs)
	if err != nil {
		l.Errorw("media.process_files_failed.get_files_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	filesByID := make(map[string]*filepb.File, len(files))
	for _, f := range files {
		filesByID[f.GetId()] = f
	}

	payloads, verr := validateAndBuildPayloads(req, filesByID, claims)
	if len(verr) > 0 {
		return nil, verr.ToStatus()
	}

	if err := s.tasksService.CreateFileProcessingTasks(ctx, payloads, nil); err != nil {
		l.Errorw("media.process_files_failed.create_tasks_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if _, err := s.fileService.UpdateFiles(ctx, buildProcessingStatusRequest(fileIDs)); err != nil {
		l.Errorw("media.process_files_failed.set_processing_status_error", "err", err)
	}

	l.Infow("media.process_files.success", "count", len(payloads))
	return &mediapb.ProcessFilesResponse{}, nil
}
