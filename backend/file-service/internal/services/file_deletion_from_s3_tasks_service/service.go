package filedeletionfroms3tasksservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/event"
	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/s3"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"go.uber.org/zap"
)

type FileDeletionFromS3TasksService interface {
	CreateFileDeletionFromS3Tasks(ctx context.Context, workerID string, tasks []FileDeletionFromS3Task, uow *uow.UnitOfWork) error
	ProcessFileDeletionFromS3Tasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *fileDeletionFromS3TasksDataProvider
	fileStorage  *fileStorage
}

func New(
	pgClient *postgres.PostgresClient,
	userPhotosS3Client *s3.S3Client,
	chatFilesS3Client *s3.S3Client,
	log *zap.SugaredLogger,
) FileDeletionFromS3TasksService {
	return &service{
		log:          log,
		dataProvider: newFileDeletionFromS3TasksDataProvider(pgClient),
		fileStorage:  newFileStorage(userPhotosS3Client, chatFilesS3Client),
	}
}

func (s *service) CreateFileDeletionFromS3Tasks(ctx context.Context, workerID string, tasks []FileDeletionFromS3Task, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "create_file_deletion_from_s3_task", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "worker_id", workerID)

	if len(tasks) == 0 {
		return nil
	}

	events := make([]*event.Event, 0, len(tasks))
	for _, task := range tasks {
		if len(task.S3Keys) == 0 {
			continue
		}
		jsonPayload, err := json.Marshal(task)
		if err != nil {
			l.Errorw("file_deletion_from_s3_tasks.create_file_deletion_from_s3_tasks_failed.marshal_error", "err", err, "payload", task)
			continue
		}
		e := event.New("", jsonPayload)
		events = append(events, e)
	}

	if len(events) == 0 {
		return nil
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("file_deletion_from_s3_tasks.create_file_deletion_from_s3_tasks_failed.begin_transaction_error", "err", err)
			return ErrCreateFileDeletionFromS3Task
		}
	}

	if err := s.dataProvider.createFileDeletionFromS3Tasks(ctx, events, uow); err != nil {
		l.Errorw("file_deletion_from_s3_tasks.create_file_deletion_from_s3_tasks_failed.create_error", "err", err)
		return ErrCreateFileDeletionFromS3Task
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file_deletion_from_s3_tasks.create_file_deletion_from_s3_tasks_failed.commit_error", "err", err)
			return ErrCreateFileDeletionFromS3Task
		}
	}

	l.Infow("file_deletion_from_s3_tasks.create_file_deletion_from_s3_tasks.success", "count", len(events))
	return nil
}

func (s *service) ProcessFileDeletionFromS3Tasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "process_file_deletion_from_s3_tasks", "worker_id", workerID)
	now := time.Now()
	retryAfter := now.Add(-1 * time.Duration(retryIntervalMS) * time.Millisecond)

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.begin_transaction_error", "err", err)
			return ErrProcessFileDeletionFromS3Tasks
		}
	}

	events, err := s.dataProvider.getFileDeletionFromS3TasksLocked(ctx, batchSize, retryAfter, uow)
	if err != nil {
		l.Errorw("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.get_tasks_error", "err", err)
		return ErrProcessFileDeletionFromS3Tasks
	}
	if events == nil {
		return nil
	}

	eventsFailed := make([]*event.Event, 0, len(events))
	eventsSuccess := make([]*event.Event, 0, len(events))

	for _, evt := range events {
		var payload FileDeletionFromS3Task
		if err := json.Unmarshal(evt.GetPayload(), &payload); err != nil {
			l.Errorw(
				"file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.unmarshal_payload_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markFileDeletionFromS3TaskFailed(evt, now, l, "file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed")
			if err != nil {
				continue
			}
			eventsFailed = append(eventsFailed, evt)
			continue
		}

		if err := s.processFileDeletionFromS3Task(ctx, payload); err != nil {
			l.Errorw(
				"file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.process_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markFileDeletionFromS3TaskFailed(evt, now, l, "file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed")
			if err != nil {
				continue
			}
			eventsFailed = append(eventsFailed, evt)
			continue
		}

		eventsSuccess = append(eventsSuccess, evt)
	}

	if len(eventsSuccess) > 0 {
		if err := s.dataProvider.deleteFileDeletionFromS3Tasks(ctx, eventsSuccess, uow); err != nil {
			l.Errorw("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.delete_error", "err", err)
			return ErrProcessFileDeletionFromS3Tasks
		}
	}
	if len(eventsFailed) > 0 {
		if err := s.dataProvider.updateFileDeletionFromS3Tasks(ctx, eventsFailed, uow); err != nil {
			l.Errorw("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.update_error", "err", err)
			return ErrProcessFileDeletionFromS3Tasks
		}
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks_failed.commit_error", "err", err)
			return ErrProcessFileDeletionFromS3Tasks
		}
	}

	if len(eventsSuccess) > 0 {
		l.Infow("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks.success", "successfully_processed", len(eventsSuccess))
	}
	if len(eventsFailed) > 0 {
		l.Warnw("file_deletion_from_s3_tasks.process_file_deletion_from_s3_tasks.success", "not_processed", len(eventsFailed))
	}

	return nil
}

func (s *service) processFileDeletionFromS3Task(
	ctx context.Context,
	task FileDeletionFromS3Task,
) error {
	bucket := filevariant.BucketType(task.Bucket)
	switch bucket {
	case filevariant.UserPhotosBucket:
		return s.fileStorage.userPhotosDeleteMany(ctx, task.S3Keys)
	case filevariant.ChatFilesBucket:
		return s.fileStorage.chatFilesDeleteMany(ctx, task.S3Keys)
	default:
		return fmt.Errorf("unknown bucket type: %d", task.Bucket)
	}
}

func (s *service) markFileDeletionFromS3TaskFailed(
	evt *event.Event,
	now time.Time,
	log *zap.SugaredLogger,
	logPrefix string,
) error {
	err := evt.IncrementAttempts()
	if err != nil {
		log.Errorw(
			fmt.Sprintf("%s.mark_task_as_failed_error", logPrefix),
			"err", err,
			"event", evt.GetID(),
			"attempts", evt.GetAttempts(),
			"status", evt.GetStatus(),
		)
		return err
	}
	evt.SetUpdatedAt(&now)
	if err := evt.SetStatus(event.EventStatusFailed); err != nil {
		log.Errorw(
			fmt.Sprintf("%s.mark_task_as_failed_error", logPrefix),
			"err", err,
			"event", evt.GetID(),
			"attempts", evt.GetAttempts(),
			"status", evt.GetStatus(),
		)
		return err
	}
	return nil
}
