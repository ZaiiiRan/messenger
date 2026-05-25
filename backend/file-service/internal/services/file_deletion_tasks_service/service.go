package filedeletiontasksservice

import (
	"context"
	"encoding/json"
	"time"

	consumersmodels "github.com/ZaiiiRan/messenger/backend/file-service/internal/consumers/models"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/event"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	fileservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_service"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type FileDeletionTasksService interface {
	CreateFileDeletionTasks(ctx context.Context, workerID string, taskMessages []*consumersmodels.FileDeletionTask, uow *uow.UnitOfWork) error
	ProcessFileDeletionTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *fileDeletionTasksDataProvider
	fileService  fileservice.FileService
}

func New(
	pgClient *postgres.PostgresClient,
	fileService fileservice.FileService,
	log *zap.SugaredLogger,
) FileDeletionTasksService {
	return &service{
		log:          log,
		dataProvider: newFileDeletionTasksDataProvider(pgClient),
		fileService:  fileService,
	}
}

func (s *service) CreateFileDeletionTasks(
	ctx context.Context,
	workerID string,
	taskMessages []*consumersmodels.FileDeletionTask,
	uow *uow.UnitOfWork,
) error {
	l := s.log.With("op", "create_file_deletion_tasks", "worker_id", workerID)

	inboxEvents := make([]*event.Event, 0, len(taskMessages))
	inboxEventsByID := make(map[string]*event.Event)
	for _, taskMessage := range taskMessages {
		if taskMessage.FileIds == nil {
			continue
		}
		evt, err := s.createFileDeletionTask(taskMessage)
		if err != nil {
			l.Warnw(
				"file_deletion_tasks.create_file_deletion_tasks_failed.marshal_payload_error",
				"err", err,
				"task_message", taskMessage,
			)
			continue
		}
		if value, ok := inboxEventsByID[evt.GetID()]; ok {
			l.Warnw(
				"file_deletion_tasks.create_file_deletion_tasks_failed.duplicate_event_error",
				"task_message", value.GetPayload(),
				"duplicate_task_message", taskMessage,
				"action", "using_last_duplicate",
			)
		}
		inboxEventsByID[evt.GetID()] = evt
	}

	for _, value := range inboxEventsByID {
		inboxEvents = append(inboxEvents, value)
	}
	if len(inboxEvents) == 0 {
		return nil
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		_, err := uow.BeginTransaction(ctx)
		if err != nil {
			l.Errorw("file_deletion_tasks.create_file_deletion_tasks_failed.begin_transaction_error", "err", err)
			return ErrCreateFileDeletionTask
		}
	}

	if err := s.dataProvider.createFileDeletionTasks(ctx, inboxEvents, uow); err != nil {
		l.Errorw("file_deletion_tasks.create_file_deletion_tasks_failed.create_tasks_error", "err", err)
		return ErrCreateFileDeletionTask
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file_deletion_tasks.create_file_deletion_tasks_failed.commit_error", "err", err)
			return ErrCreateFileDeletionTask
		}
	}

	l.Infow("file_deletion_tasks.create_file_deletion_tasks_success", "count", len(inboxEvents))
	return nil
}

func (s *service) ProcessFileDeletionTasks(
	ctx context.Context,
	workerID string,
	retryIntervalMS uint,
	batchSize int,
	uow *uow.UnitOfWork,
) error {
	l := s.log.With("op", "process_file_deletion_tasks", "worker_id", workerID)

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		_, err := uow.BeginTransaction(ctx)
		if err != nil {
			l.Errorw("file_deletion_tasks.process_file_deletion_tasks_failed.begin_transaction_error", "err", err)
			return ErrProcessFileDeletionTasks
		}
		needToCommit = true
	}

	now := time.Now()
	retryAfter := now.Add(-1 * time.Duration(retryIntervalMS) * time.Millisecond)
	inboxEvents, err := s.dataProvider.getFileDeletionTasksLocked(ctx, batchSize, retryAfter, uow)
	if err != nil {
		l.Errorw("file_deletion_tasks.process_file_deletion_tasks_failed.get_tasks_error", "err", err)
		return ErrProcessFileDeletionTasks
	}
	if inboxEvents == nil {
		return nil
	}

	inboxEventsSuccess := make([]*event.Event, 0, len(inboxEvents))
	inboxEventsFailed := make([]*event.Event, 0, len(inboxEvents))

	for _, evt := range inboxEvents {
		var taskMessage *consumersmodels.FileDeletionTask
		if err := json.Unmarshal(evt.GetPayload(), &taskMessage); err != nil {
			l.Errorw(
				"file_deletion_tasks.process_file_deletion_tasks_failed.unmarshal_payload_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markFileDeletionTaskFailed(evt, now, l, "file_deletion_tasks.process_file_deletion_tasks_failed")
			if err != nil {
				continue
			}
			inboxEventsFailed = append(inboxEventsFailed, evt)
			continue
		}

		err := s.processFileDeletionTask(ctx, workerID, taskMessage, uow)
		if err != nil {
			l.Errorw(
				"file_deletion_tasks.process_file_deletion_tasks_failed.process_task_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markFileDeletionTaskFailed(evt, now, l, "file_deletion_tasks.process_file_deletion_tasks_failed")
			if err != nil {
				continue
			}
			inboxEventsFailed = append(inboxEventsFailed, evt)
			continue
		}

		inboxEventsSuccess = append(inboxEventsSuccess, evt)
	}

	if len(inboxEventsSuccess) > 0 {
		if err := s.dataProvider.deleteFileDeletionTasks(ctx, inboxEventsSuccess, uow); err != nil {
			l.Errorw("file_deletion_tasks.process_file_deletion_tasks_failed.delete_tasks_error", "err", err)
			return ErrProcessFileDeletionTasks
		}
	}
	if len(inboxEventsFailed) > 0 {
		if err := s.dataProvider.updateFileDeletionTasks(ctx, inboxEventsFailed, uow); err != nil {
			l.Errorw("file_deletion_tasks.process_file_deletion_tasks_failed.update_tasks_error", "err", err)
			return ErrProcessFileDeletionTasks
		}
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("file_deletion_tasks.process_file_deletion_tasks_failed.commit_error", "err", err)
			return ErrProcessFileDeletionTasks
		}
	}

	if len(inboxEventsSuccess) > 0 {
		l.Infow("file_deletion_tasks.process_file_deletion_tasks_success", "successfully_processed", len(inboxEventsSuccess))
	}
	if len(inboxEventsFailed) > 0 {
		l.Warnw("file_deletion_tasks.process_file_deletion_tasks_failed", "failed", len(inboxEventsFailed))
	}

	return nil
}

func (s *service) processFileDeletionTask(
	ctx context.Context,
	workerID string,
	taskMessage *consumersmodels.FileDeletionTask,
	uow *uow.UnitOfWork,
) error {
	return s.fileService.DeleteFilesByIDs(ctx, workerID, taskMessage.FileIds, uow)
}

func (s *service) markFileDeletionTaskFailed(
	evt *event.Event,
	now time.Time,
	log *zap.SugaredLogger,
	logPrefix string,
) error {
	err := evt.IncrementAttempts()
	if err != nil {
		log.Errorw(logPrefix+".mark_task_as_failed_error",
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
			logPrefix+".mark_task_as_failed_error",
			"err", err,
			"event", evt.GetID(),
			"attempts", evt.GetAttempts(),
			"status", evt.GetStatus(),
		)
		return err
	}
	return nil
}

func (s *service) createFileDeletionTask(taskMessage *consumersmodels.FileDeletionTask) (*event.Event, error) {
	jsonPayload, err := json.Marshal(taskMessage)
	if err != nil {
		return nil, err
	}
	return event.New(taskMessage.Id, jsonPayload), nil
}
