package fileprocessingtasksservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/event"
	domainfile "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file"
	filevariant "github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/processing"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/mediaprocessor"
	producersmodels "github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/models"
	uow "github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/unitofwork/postgres"
	fileprocessingcompletedeventsservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_completed_events_service"
	fileservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_service"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/s3"
	"go.uber.org/zap"
)

type FileProcessingTasksService interface {
	CreateFileProcessingTasks(ctx context.Context, payloads []*TaskPayload, uow *uow.UnitOfWork) error
	ProcessFileProcessingTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error
}

type service struct {
	log                    *zap.SugaredLogger
	dataProvider           *fileProcessingTasksDataProvider
	fileService            fileservice.FileService
	processor              mediaprocessor.Processor
	fileStorage            *fileStorage
	completedEventsService fileprocessingcompletedeventsservice.FileProcessingCompletedEventsService
}

func New(
	pgClient *postgres.PostgresClient,
	fileService fileservice.FileService,
	processor mediaprocessor.Processor,
	userPhotosS3Client *s3.S3Client,
	chatFilesS3Client *s3.S3Client,
	completedEventsService fileprocessingcompletedeventsservice.FileProcessingCompletedEventsService,
	log *zap.SugaredLogger,
) FileProcessingTasksService {
	return &service{
		log:                    log,
		dataProvider:           newFileProcessingTasksDataProvider(pgClient),
		fileService:            fileService,
		processor:              processor,
		fileStorage:            newFileStorage(userPhotosS3Client, chatFilesS3Client),
		completedEventsService: completedEventsService,
	}
}

func (s *service) CreateFileProcessingTasks(ctx context.Context, payloads []*TaskPayload, uowArg *uow.UnitOfWork) error {
	l := s.log.With("op", "create_file_processing_tasks", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(payloads) == 0 {
		return nil
	}

	events := make([]*event.Event, 0, len(payloads))
	for _, p := range payloads {
		jsonPayload, err := json.Marshal(p)
		if err != nil {
			l.Errorw("file_processing_tasks.create_file_processing_tasks_failed.marshal_error", "err", err)
			return ErrCreateFileProcessingTasks
		}
		events = append(events, event.New("", jsonPayload))
	}

	needToCommit := false
	transaction := uowArg
	if transaction == nil {
		transaction = s.dataProvider.newUOW()
		defer transaction.Close()
		needToCommit = true
		if _, err := transaction.BeginTransaction(ctx); err != nil {
			l.Errorw("file_processing_tasks.create_file_processing_tasks_failed.begin_transaction_error", "err", err)
			return ErrCreateFileProcessingTasks
		}
	}

	if err := s.dataProvider.createTasks(ctx, events, transaction); err != nil {
		l.Errorw("file_processing_tasks.create_file_processing_tasks_failed.create_error", "err", err)
		return ErrCreateFileProcessingTasks
	}

	if needToCommit {
		if err := transaction.Commit(ctx); err != nil {
			l.Errorw("file_processing_tasks.create_file_processing_tasks_failed.commit_error", "err", err)
			return ErrCreateFileProcessingTasks
		}
	}

	l.Infow("file_processing_tasks.create_file_processing_tasks.success", "count", len(events))
	return nil
}

func (s *service) ProcessFileProcessingTasks(
	ctx context.Context,
	workerID string,
	retryIntervalMS uint,
	batchSize int,
	uowArg *uow.UnitOfWork,
) error {
	l := s.log.With("op", "process_file_processing_tasks", "worker_id", workerID)

	now := time.Now()
	retryAfter := now.Add(-1 * time.Duration(retryIntervalMS) * time.Millisecond)

	needToCommit := false
	transaction := uowArg
	if transaction == nil {
		transaction = s.dataProvider.newUOW()
		defer transaction.Close()
		needToCommit = true
		if _, err := transaction.BeginTransaction(ctx); err != nil {
			l.Errorw("file_processing_tasks.process_file_processing_tasks_failed.begin_transaction_error", "err", err)
			return ErrProcessFileProcessingTasks
		}
	}

	tasks, err := s.dataProvider.getTasksLocked(ctx, batchSize, retryAfter, transaction)
	if err != nil {
		l.Errorw("file_processing_tasks.process_file_processing_tasks_failed.get_tasks_error", "err", err)
		return ErrProcessFileProcessingTasks
	}
	if tasks == nil {
		return nil
	}

	tasksSuccess := make([]*event.Event, 0, len(tasks))
	tasksFailed := make([]*event.Event, 0, len(tasks))
	completedPayloads := make([]*producersmodels.FileProcessingCompletedEvent, 0, len(tasks))

	for _, evt := range tasks {
		var payload TaskPayload
		if err := json.Unmarshal(evt.GetPayload(), &payload); err != nil {
			l.Errorw(
				"file_processing_tasks.process_file_processing_tasks_failed.unmarshal_payload_error",
				"err", err, "event", evt.GetID(), "attempts", evt.GetAttempts(), "status", evt.GetStatus(),
			)
			if err := s.markTaskFailed(evt, now, l); err != nil {
				continue
			}
			tasksFailed = append(tasksFailed, evt)
			continue
		}

		completed, err := s.processTask(ctx, &payload)
		if err != nil {
			l.Errorw(
				"file_processing_tasks.process_file_processing_tasks_failed.process_error",
				"err", err, "event", evt.GetID(), "file_id", payload.FileID,
				"attempts", evt.GetAttempts(), "status", evt.GetStatus(),
			)
			if err := s.markTaskFailed(evt, now, l); err != nil {
				continue
			}
			tasksFailed = append(tasksFailed, evt)
			continue
		}

		completedPayloads = append(completedPayloads, completed)
		tasksSuccess = append(tasksSuccess, evt)
	}

	if len(completedPayloads) > 0 {
		if err := s.completedEventsService.CreateFileProcessingCompletedEvents(ctx, completedPayloads, transaction); err != nil {
			l.Errorw("file_processing_tasks.process_file_processing_tasks_failed.create_completed_events_error", "err", err)
			return ErrProcessFileProcessingTasks
		}
	}
	if len(tasksSuccess) > 0 {
		if err := s.dataProvider.deleteTasks(ctx, tasksSuccess, transaction); err != nil {
			l.Errorw("file_processing_tasks.process_file_processing_tasks_failed.delete_tasks_error", "err", err)
			return ErrProcessFileProcessingTasks
		}
	}
	if len(tasksFailed) > 0 {
		if err := s.dataProvider.updateTasks(ctx, tasksFailed, transaction); err != nil {
			l.Errorw("file_processing_tasks.process_file_processing_tasks_failed.update_tasks_error", "err", err)
			return ErrProcessFileProcessingTasks
		}
	}
	if needToCommit {
		if err := transaction.Commit(ctx); err != nil {
			l.Errorw("file_processing_tasks.process_file_processing_tasks_failed.commit_error", "err", err)
			return ErrProcessFileProcessingTasks
		}
	}

	if len(tasksSuccess) > 0 {
		l.Infow("file_processing_tasks.process_file_processing_tasks.success", "processed", len(tasksSuccess))
	}
	if len(tasksFailed) > 0 {
		l.Warnw("file_processing_tasks.process_file_processing_tasks.partial", "failed", len(tasksFailed))
	}

	return nil
}

func (s *service) processTask(ctx context.Context, payload *TaskPayload) (*producersmodels.FileProcessingCompletedEvent, error) {
	bucket := filevariant.BucketType(payload.Bucket)

	data, err := s.fileStorage.get(ctx, bucket, payload.OriginalS3Key)
	if err != nil {
		return nil, fmt.Errorf("download original: %w", err)
	}

	pt, err := processing.ProcessingTypeFromString(payload.ProcessingType)
	if err != nil {
		return nil, err
	}

	variants, err := s.processor.Process(ctx, mediaprocessor.Request{
		ProcessingType: pt,
		SourceMime:     payload.OriginalMime,
		SourceData:     data,
		SourceWidth:    payload.Width,
		SourceHeight:   payload.Height,
		SourceDuration: payload.Duration,
		Crop:           payload.Crop,
	})
	if err != nil {
		return nil, fmt.Errorf("process media: %w", err)
	}
	if len(variants) == 0 {
		return nil, fmt.Errorf("no variants produced")
	}

	uploaded := make([]uploadedVariant, 0, len(variants))
	for _, v := range variants {
		key := variantKey(payload.FileID, v.Variant, v.MimeType)
		if v.Variant == filevariant.VariantOriginal {
			key = payload.OriginalS3Key
		}
		if err := s.fileStorage.put(ctx, bucket, key, v.Data, v.MimeType); err != nil {
			return nil, fmt.Errorf("upload variant %s: %w", v.Variant.String(), err)
		}
		uploaded = append(uploaded, uploadedVariant{
			variant:  v.Variant,
			s3Key:    key,
			mimeType: v.MimeType,
			meta:     filevariant.Meta{Width: v.Width, Height: v.Height, Duration: v.Duration},
		})
	}

	updateReq, err := buildFinalizeRequest(payload, uploaded)
	if err != nil {
		return nil, err
	}
	if _, err := s.fileService.UpdateFiles(ctx, updateReq); err != nil {
		return nil, fmt.Errorf("register variants: %w", err)
	}

	completedVariants := make([]producersmodels.FileProcessingCompletedVariant, 0, len(variants))
	for _, v := range variants {
		completedVariants = append(completedVariants, producersmodels.FileProcessingCompletedVariant{
			VariantType: v.Variant.String(),
			MimeType:    v.MimeType,
			Width:       v.Width,
			Height:      v.Height,
			Duration:    v.Duration,
		})
	}

	return &producersmodels.FileProcessingCompletedEvent{
		FileId:         payload.FileID,
		OwnerId:        payload.OwnerID,
		ProcessingType: payload.ProcessingType,
		Purpose:        payload.Purpose,
		Status:         domainfile.FileStatusReady.String(),
		Variants:       completedVariants,
		CompletedAt:    time.Now(),
	}, nil
}

func (s *service) markTaskFailed(evt *event.Event, now time.Time, log *zap.SugaredLogger) error {
	if err := evt.IncrementAttempts(); err != nil {
		log.Errorw(
			"file_processing_tasks.mark_task_as_failed_error",
			"err", err, "event", evt.GetID(), "attempts", evt.GetAttempts(), "status", evt.GetStatus(),
		)
		return err
	}
	evt.SetUpdatedAt(&now)
	if err := evt.SetStatus(event.EventStatusFailed); err != nil {
		log.Errorw(
			"file_processing_tasks.mark_task_as_failed_error",
			"err", err, "event", evt.GetID(), "attempts", evt.GetAttempts(), "status", evt.GetStatus(),
		)
		return err
	}
	return nil
}
