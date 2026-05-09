package userdatadeletiontasksservice

import (
	"context"
	"encoding/json"

	consumersmodels "github.com/ZaiiiRan/messenger/backend/auth-service/internal/consumers/models"
	inboxevent "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/inbox_event"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type UserDataDeletionTasksService interface {
	CreateUserDataDeletionTasks(ctx context.Context, workerID string, taskMessages []*consumersmodels.UserDataDeletionTask) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *userDataDeletionTasksDataProvider
}

func New(
	pgClient *postgres.PostgresClient,
	log *zap.SugaredLogger,
) UserDataDeletionTasksService {
	return &service{
		log:          log,
		dataProvider: newUserDataDeletionTasksDataProvider(pgClient),
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
		return ErrCreateUserDataDeletionTask
	}

	if err := s.dataProvider.createUserDataDeletionTasks(ctx, inboxEvents, uow); err != nil {
		l.Errorw("user_data_deletion_tasks.create_user_data_deletion_tasks_failed.create_error", "err", err)
		return ErrCreateUserDataDeletionTask
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_data_deletion_tasks.create_user_data_deletion_tasks_failed.commit_error", "err", err)
		return ErrCreateUserDataDeletionTask
	}

	l.Infow("user_data_deletion_tasks.create_user_data_deletion_tasks.success", "count", len(inboxEvents))
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
