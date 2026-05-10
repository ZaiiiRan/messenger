package userdatadeletiontasksservice

import (
	"context"
	"encoding/json"
	"time"

	consumersmodels "github.com/ZaiiiRan/messenger/backend/auth-service/internal/consumers/models"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	inboxevent "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/inbox_event"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	codeservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type UserDataDeletionTasksService interface {
	CreateUserDataDeletionTasks(ctx context.Context, workerID string, taskMessages []*consumersmodels.UserDataDeletionTask) error
	ProcessUserDataDeletionTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error
}

type service struct {
	log             *zap.SugaredLogger
	dataProvider    *userDataDeletionTasksDataProvider
	passwordService passwordservice.PasswordService
	codeService     codeservice.CodeService
	tokenService    tokenservice.TokenService
}

func New(
	pgClient *postgres.PostgresClient,
	passwordService passwordservice.PasswordService,
	codeService codeservice.CodeService,
	tokenService tokenservice.TokenService,
	log *zap.SugaredLogger,
) UserDataDeletionTasksService {
	return &service{
		log:             log,
		dataProvider:    newUserDataDeletionTasksDataProvider(pgClient),
		passwordService: passwordService,
		codeService:     codeService,
		tokenService:    tokenService,
	}
}

func (s *service) CreateUserDataDeletionTasks(ctx context.Context, workerID string, taskMessages []*consumersmodels.UserDataDeletionTask) error {
	l := s.log.With("op", "create_user_data_deletion_tasks", "worker_id", workerID)

	inboxEvents := make([]*inboxevent.InboxEvent, 0, len(taskMessages))
	inboxEventsByID := make(map[string]*inboxevent.InboxEvent)
	for _, taskMessage := range taskMessages {
		event, err := s.createUserDataDeletionTask(taskMessage)
		if err != nil {
			l.Warnw(
				"user_data_deletion_tasks.create_user_data_deletion_tasks_warning.json_marshal_error",
				"err", err,
				"task_message", taskMessage,
			)
			continue
		}
		if value, ok := inboxEventsByID[event.GetID()]; ok {
			l.Warnw(
				"user_data_deletion_tasks.create_user_data_deletion_tasks_warning.duplicate_task",
				"task_mesage", value.GetPayload(),
				"duplicate_task_message", taskMessage,
				"action", "using_last_duplicate",
			)
		}
		inboxEventsByID[event.GetID()] = event
	}

	for _, value := range inboxEventsByID {
		inboxEvents = append(inboxEvents, value)
	}
	if len(inboxEvents) == 0 {
		return nil
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_data_deletion_tasks.create_user_data_deletion_tasks_failed.begin_transaction_error", "err", err)
		return ErrCreateUserDataDeletionTasks
	}

	if err := s.dataProvider.createUserDataDeletionTasks(ctx, inboxEvents, uow); err != nil {
		l.Errorw("user_data_deletion_tasks.create_user_data_deletion_tasks_failed.create_error", "err", err)
		return ErrCreateUserDataDeletionTasks
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_data_deletion_tasks.create_user_data_deletion_tasks_failed.commit_error", "err", err)
		return ErrCreateUserDataDeletionTasks
	}

	l.Infow("user_data_deletion_tasks.create_user_data_deletion_tasks.success", "count", len(inboxEvents))
	return nil
}

func (s *service) ProcessUserDataDeletionTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error {
	l := s.log.With("op", "process_user_data_deletion_tasks", "worker_id", workerID)

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_data_deletion_tasks.process_user_data_deletion_tasks_failed.begin_transaction_error", "err", err)
		return ErrProcessUserDataDeletionTasks
	}

	now := time.Now()
	retryAfter := now.Add(-1 * (time.Duration(retryIntervalMS) * time.Millisecond))
	createdAfter := now.Add(-5 * time.Minute)
	inboxEvents, err := s.dataProvider.getUserDataDeletionTasksLocked(ctx, batchSize, retryAfter, &createdAfter, uow)
	if err != nil {
		l.Errorw("user_data_deletion_tasks.process_user_data_deletion_tasks_failed.get_tasks_error", "err", err)
		return ErrProcessUserDataDeletionTasks
	}
	if inboxEvents == nil {
		return nil
	}

	inboxEventsSuccess := make([]*inboxevent.InboxEvent, 0, len(inboxEvents))
	inboxEventsFailed := make([]*inboxevent.InboxEvent, 0, len(inboxEvents))

	for _, event := range inboxEvents {
		var taskMessage *consumersmodels.UserDataDeletionTask
		if err := json.Unmarshal(event.GetPayload(), &taskMessage); err != nil {
			l.Errorw(
				"user_data_deletion_tasks.process_user_data_deletion_tasks_failed.unmarshal_payload_error",
				"err", err,
				"event", event.GetID(),
				"payload", event.GetPayload(),
				"attempts", event.GetAttempts(),
				"status", event.GetStatus(),
			)
			err = s.markUserDataDeletionTaskFailed(event, now, l, "user_data_deletion_tasks.process_user_data_deletion_tasks_failed")
			if err != nil {
				continue
			}
			inboxEventsFailed = append(inboxEventsFailed, event)
			continue
		}

		if err := s.processUserDataDeletionTask(ctx, workerID, taskMessage, uow); err != nil {
			l.Errorw(
				"user_data_deletion_tasks.process_user_data_deletion_tasks_failed.process_task_error",
				"err", err,
				"event", event.GetID(),
				"payload", event.GetPayload(),
				"attempts", event.GetAttempts(),
				"status", event.GetStatus(),
			)
			err = s.markUserDataDeletionTaskFailed(event, now, l, "user_data_deletion_tasks.process_user_data_deletion_tasks_failed")
			if err != nil {
				continue
			}
			inboxEventsFailed = append(inboxEventsFailed, event)
			continue
		}

		inboxEventsSuccess = append(inboxEventsSuccess, event)
	}

	if len(inboxEventsSuccess) > 0 {
		if err := s.dataProvider.deleteUserDataDeletionTasks(ctx, inboxEventsSuccess, uow); err != nil {
			l.Errorw("user_data_deletion_tasks.process_user_data_deletion_tasks_failed.delete_error", "err", err)
			return ErrProcessUserDataDeletionTasks
		}
	}
	if len(inboxEventsFailed) > 0 {
		if err := s.dataProvider.updateUserDataDeletionTasks(ctx, inboxEventsFailed, uow); err != nil {
			l.Errorw("user_data_deletion_tasks.process_user_data_deletion_tasks_failed.update_error", "err", err)
			return ErrProcessUserDataDeletionTasks
		}
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_data_deletion_tasks.process_user_data_deletion_tasks_failed.commit_error", "err", err)
		return ErrProcessUserDataDeletionTasks
	}

	if len(inboxEventsSuccess) > 0 {
		l.Infow("user_data_deletion_tasks.process_user_data_deletion_tasks.success", "successfully_processed", len(inboxEventsSuccess))
	}
	if len(inboxEventsFailed) > 0 {
		l.Warnw("user_data_deletion_tasks.process_user_data_deletion_tasks.success", "failed", len(inboxEventsFailed))
	}

	return nil
}

func (s *service) processUserDataDeletionTask(
	ctx context.Context,
	workerID string,
	taskMessage *consumersmodels.UserDataDeletionTask,
	uow *uow.UnitOfWork,
) error {
	if err := s.passwordService.DeletePasswordByUserID(ctx, workerID, uow, taskMessage.UserId); err != nil {
		return err
	}
	if err := s.codeService.DeleteCodeByUserID(ctx, workerID, uow, taskMessage.UserId, code.CodeTypePasswordReset); err != nil {
		return err
	}
	if err := s.codeService.DeleteCodeByUserID(ctx, workerID, uow, taskMessage.UserId, code.CodeTypeActivation); err != nil {
		return err
	}
	if err := s.tokenService.DeleteUserVersionAndTokensByUserID(ctx, workerID, uow, taskMessage.UserId); err != nil {
		return err
	}
	return nil
}

func (s *service) markUserDataDeletionTaskFailed(
	event *inboxevent.InboxEvent,
	now time.Time,
	log *zap.SugaredLogger,
	logPrefix string,
) error {
	err := event.IncrementAttempts()
	if err != nil {
		log.Errorw(
			logPrefix+".mark_task_as_failed_error",
			"err", err,
			"event", event.GetID(),
			"attempts", event.GetAttempts(),
			"status", event.GetStatus(),
		)
		return err
	}
	event.SetUpdatedAt(&now)
	if err := event.SetStatus(inboxevent.InboxEventStatusFailed); err != nil {
		log.Errorw(
			logPrefix+".mark_task_as_failed_error",
			"err", err,
			"event", event.GetID(),
			"attempts", event.GetAttempts(),
			"status", event.GetStatus(),
		)
		return err
	}
	return nil
}

func (s *service) createUserDataDeletionTask(taskMessage *consumersmodels.UserDataDeletionTask) (*inboxevent.InboxEvent, error) {
	jsonPayload, err := json.Marshal(taskMessage)
	if err != nil {
		return nil, err
	}

	event := inboxevent.New(taskMessage.Id, jsonPayload)
	return event, nil
}
