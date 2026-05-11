package userdatadeletiontasksservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	outboxevent "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/outbox_event"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	producersinterfaces "github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/interfaces"
	producersmodels "github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/models"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type UserDataDeletionTasksService interface {
	CreateUserDataDeletionTasks(ctx context.Context, workerID string, users []*user.User, uow *uow.UnitOfWork) error
	SendUserDataDeletionTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error
}

type service struct {
	log                           *zap.SugaredLogger
	dataProvider                  *userDataDeletionTasksDataProvider
	userDataDeletionTasksProducer producersinterfaces.UserDataDeletionTasksProducer
}

func New(
	pgClient *postgres.PostgresClient,
	userDataDeletionTasksProducer producersinterfaces.UserDataDeletionTasksProducer,
	log *zap.SugaredLogger,
) UserDataDeletionTasksService {
	return &service{
		log:                           log,
		userDataDeletionTasksProducer: userDataDeletionTasksProducer,
		dataProvider:                  newUserDataDeletionTasksDataProvider(pgClient),
	}
}

func (s *service) CreateUserDataDeletionTasks(ctx context.Context, workerID string, users []*user.User, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "create_user_data_deletion_tasks", "worker_id", workerID)

	outboxEvents := make([]*outboxevent.OutboxEvent, 0, len(users))
	for _, user := range users {
		event, err := s.createUserDataDeletionTask(user)
		if err != nil {
			l.Errorw(
				"user_data_deletion_tasks.create_user_data_deletion_tasks_failed.json_marshal_error",
				"err", err,
				"user", user,
			)
			return ErrMarshalPayload
		}
		outboxEvents = append(outboxEvents, event)
	}

	if err := s.dataProvider.createUserDataDeletionTasks(ctx, outboxEvents, uow); err != nil {
		l.Errorw("user_data_deletion_tasks.create_user_data_deletion_tasks_failed.create_error", "err", err)
		return ErrCreateUserDataDeletionTask
	}

	if len(outboxEvents) > 0 {
		l.Infow("user_data_deletion_tasks.create_user_data_deletion_tasks.success", "count", len(outboxEvents))
	}

	return nil
}

func (s *service) SendUserDataDeletionTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error {
	l := s.log.With("op", "send_user_data_deletion_tasks", "worker_id", workerID)

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_data_deletion_tasks.send_user_data_deletion_tasks_failed.begin_transaction_error", "err", err)
		return ErrSendUserDataDeletionTasks
	}

	now := time.Now()
	retryAfter := now.Add(-1 * (time.Duration(retryIntervalMS) * time.Millisecond))
	outboxEvents, err := s.dataProvider.getUserDataDeletionTasksLocked(ctx, batchSize, retryAfter, uow)
	if err != nil {
		l.Errorw("user_data_deletion_tasks.send_user_data_deletion_tasks_failed.get_tasks_error", "err", err)
		return ErrSendUserDataDeletionTasks
	}
	if outboxEvents == nil {
		return nil
	}

	outboxEventsFailed := make([]*outboxevent.OutboxEvent, 0, len(outboxEvents))
	outboxEventsSuccess := make([]*outboxevent.OutboxEvent, 0, len(outboxEvents))

	for _, event := range outboxEvents {
		var payload producersmodels.UserDataDeletionTask
		if err := json.Unmarshal(event.GetPayload(), &payload); err != nil {
			l.Errorw(
				"user_data_deletion_tasks.send_user_data_deletion_tasks_failed.unmarshal_payload_error",
				"err", err,
				"event", event.GetID(),
				"payload", event.GetPayload(),
				"attempts", event.GetAttempts(),
				"status", event.GetStatus(),
			)
			err = s.markUserDataDeletionTaskFailed(event, now, l, "user_data_deletion_tasks.send_user_data_deletion_tasks_failed")
			if err != nil {
				continue
			}
			outboxEventsFailed = append(outboxEventsFailed, event)
			continue
		}
		payload.Id = event.GetID()

		if err := s.userDataDeletionTasksProducer.ProduceUserDataDeletionTask(ctx, &payload); err != nil {
			l.Errorw(
				"user_data_deletion_tasks.send_user_data_deletion_tasks_failed.produce_error",
				"err", err,
				"event", event.GetID(),
				"payload", event.GetPayload(),
				"attempts", event.GetAttempts(),
				"status", event.GetStatus(),
			)
			err = s.markUserDataDeletionTaskFailed(event, now, l, "user_data_deletion_tasks.send_user_data_deletion_tasks_failed")
			if err != nil {
				continue
			}
			outboxEventsFailed = append(outboxEventsFailed, event)
			continue
		}

		outboxEventsSuccess = append(outboxEventsSuccess, event)
	}

	if len(outboxEventsSuccess) > 0 {
		if err := s.dataProvider.deleteUserDataDeletionTasks(ctx, outboxEventsSuccess, uow); err != nil {
			l.Errorw("user_data_deletion_tasks.send_user_data_deletion_tasks_failed.delete_error", "err", err)
			return ErrSendUserDataDeletionTasks
		}
	}
	if len(outboxEventsFailed) > 0 {
		if err := s.dataProvider.updateUserDataDeletionTasks(ctx, outboxEventsFailed, uow); err != nil {
			l.Errorw("user_data_deletion_tasks.send_user_data_deletion_tasks_failed.update_error", "err", err)
			return ErrSendUserDataDeletionTasks
		}
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_data_deletion_tasks.send_user_data_deletion_tasks_failed.commit_error", "err", err)
		return ErrSendUserDataDeletionTasks
	}

	if len(outboxEventsSuccess) > 0 {
		l.Infow("user_data_deletion_tasks.send_user_data_deletion_tasks.success", "successfully_sended", len(outboxEventsSuccess))
	}
	if len(outboxEventsFailed) > 0 {
		l.Warnw("user_data_deletion_tasks.send_user_data_deletion_tasks.success", "not_sended", len(outboxEventsFailed))
	}

	return nil
}

func (s *service) markUserDataDeletionTaskFailed(
	event *outboxevent.OutboxEvent,
	now time.Time,
	log *zap.SugaredLogger,
	logPrefix string,
) error {
	err := event.IncrementAttempts()
	if err != nil {
		log.Errorw(
			fmt.Sprintf("%s.mark_task_as_failed_error", logPrefix),
			"err", err,
			"event", event.GetID(),
			"attempts", event.GetAttempts(),
			"status", event.GetStatus(),
		)
		return err
	}
	event.SetUpdatedAt(&now)
	if err := event.SetStatus(outboxevent.OutboxEventStatusFailed); err != nil {
		log.Errorw(
			fmt.Sprintf("%s.mark_task_as_failed_error", logPrefix),
			"err", err,
			"event", event.GetID(),
			"attempts", event.GetAttempts(),
			"status", event.GetStatus(),
		)
		return err
	}
	return nil
}

func (s *service) createUserDataDeletionTask(user *user.User) (*outboxevent.OutboxEvent, error) {
	payload := producersmodels.UserDataDeletionTask{
		UserId:      user.GetID(),
		Username:    user.GetUsername(),
		Email:       user.GetEmail(),
		IsConfirmed: user.GetStatus().IsConfirmed(),
		IsDeleted:   user.GetStatus().IsDeleted(),
		CreatedAt:   user.GetCreatedAt(),
		UpdatedAt:   user.GetUpdatedAt(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	event := outboxevent.New(jsonPayload)
	return event, nil
}
