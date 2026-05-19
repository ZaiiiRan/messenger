package userrelationshipchangestasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/event"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	producersinterfaces "github.com/ZaiiiRan/messenger/backend/social-service/internal/producers/interfaces"
	producersmodels "github.com/ZaiiiRan/messenger/backend/social-service/internal/producers/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type UserRelationshipChangesTasksService interface {
	CreateUserRelationshipChangesTasks(ctx context.Context, userRelationships []*userrelationship.UserRelationship, uow *uow.UnitOfWork) error
	SendUserRelationshipChangesTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error
}

type service struct {
	log                                  *zap.SugaredLogger
	dataProvider                         *userRelationshipChangesTasksDataProvider
	userRelationshipChangesTasksProducer producersinterfaces.UserRelationshipChangesTasksProducer
}

func New(
	pgClient *postgres.PostgresClient,
	userRelationshipChangesTasksProducer producersinterfaces.UserRelationshipChangesTasksProducer,
	log *zap.SugaredLogger,
) UserRelationshipChangesTasksService {
	return &service{
		log:                                  log,
		userRelationshipChangesTasksProducer: userRelationshipChangesTasksProducer,
		dataProvider:                         newUserRelationshipChangesTasksDataProvider(pgClient),
	}
}

func (s *service) CreateUserRelationshipChangesTasks(ctx context.Context, userRelationships []*userrelationship.UserRelationship, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "create_user_relationship_changes_tasks", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship_changes_tasks.create_user_relationship_changes_tasks_failed.begin_transaction_error", "err", err)
			return ErrCreateUserDataDeletionTask
		}
	}

	outboxEvents := make([]*event.Event, 0, len(userRelationships))
	for _, ur := range userRelationships {
		evt, err := s.createUserRelationshipChangeTask(ur)
		if err != nil {
			l.Errorw(
				"user_relationship_changes_tasks.create_user_relationship_changes_tasks_failed.json_marshal_error",
				"err", err,
				"user_relationship_change_task", ur,
			)
			return ErrMarshalPayload
		}
		outboxEvents = append(outboxEvents, evt)
	}

	if err := s.dataProvider.createUserRelationshipChangesTasks(ctx, outboxEvents, uow); err != nil {
		l.Errorw("user_relationship_changes_tasks.create_user_relationship_changes_tasks_failed.create_error", "err", err)
		return ErrCreateUserDataDeletionTask
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship_changes_tasks.create_user_relationship_changes_tasks_failed.commit_error", "err", err)
			return ErrCreateUserDataDeletionTask
		}
	}

	l.Infow("user_relationship_changes_tasks.create_user_relationship_changes_task.success", "count", len(outboxEvents))
	return nil
}

func (s *service) SendUserRelationshipChangesTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "send_user_relationship_changes_tasks", "worker_id", workerID)
	now := time.Now()
	retryAfter := now.Add(-1 * time.Duration(retryIntervalMS) * time.Millisecond)

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.begin_transaction_error", "err", err)
			return ErrSendUserDataDeletionTasks
		}
	}

	outboxEvents, err := s.dataProvider.getUserRelationshipChangesTasksLocked(ctx, batchSize, retryAfter, uow)
	if err != nil {
		l.Errorw("user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.get_tasks_error", "err", err)
		return ErrSendUserDataDeletionTasks
	}
	if outboxEvents == nil {
		return nil
	}

	outboxEventsFailed := make([]*event.Event, 0, len(outboxEvents))
	outboxEventsSuccess := make([]*event.Event, 0, len(outboxEvents))

	for _, evt := range outboxEvents {
		var payload producersmodels.UserRelationshipChangeTask
		if err := json.Unmarshal(evt.GetPayload(), &payload); err != nil {
			l.Errorw(
				"user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.unmarshal_payload_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markUserRelationshipChangeTaskFailed(evt, now, l, "user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed")
			if err != nil {
				continue
			}
			outboxEventsFailed = append(outboxEventsFailed, evt)
			continue
		}
		payload.Id = evt.GetID()

		if err := s.userRelationshipChangesTasksProducer.ProduceUserRelationshipChangesTask(ctx, &payload); err != nil {
			l.Errorw(
				"user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.produce_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markUserRelationshipChangeTaskFailed(evt, now, l, "user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed")
			if err != nil {
				continue
			}
			outboxEventsFailed = append(outboxEventsFailed, evt)
			continue
		}

		outboxEventsSuccess = append(outboxEventsSuccess, evt)
	}

	if len(outboxEventsSuccess) > 0 {
		if err := s.dataProvider.deleteUserRelationshipChangesTasks(ctx, outboxEventsSuccess, uow); err != nil {
			l.Errorw("user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.delete error", "err", err)
			return ErrSendUserDataDeletionTasks
		}
	}
	if len(outboxEventsFailed) > 0 {
		if err := s.dataProvider.updatetUserRelationshipChangesTasks(ctx, outboxEventsFailed, uow); err != nil {
			l.Errorw("user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.update error", "err", err)
			return ErrSendUserDataDeletionTasks
		}
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship_changes_tasks.send_user_relationship_changes_tasks_failed.commit error", "err", err)
			return ErrSendUserDataDeletionTasks
		}
	}

	if len(outboxEvents) > 0 {
		l.Infow("user_relationship_changes_tasks.send_user_relationship_changes_tasks.success", "successfully_sended", len(outboxEvents))
	}
	if len(outboxEventsFailed) > 0 {
		l.Warnw("user_relationship_changes_tasks.send_user_relationship_changes_tasks.success", "not_sended", len(outboxEventsFailed))
	}

	return nil
}

func (s *service) markUserRelationshipChangeTaskFailed(
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

func (s *service) createUserRelationshipChangeTask(ur *userrelationship.UserRelationship) (*event.Event, error) {
	payload := producersmodels.UserRelationshipChangeTask{
		User1Id:                ur.GetUserID1(),
		User2Id:                ur.GetUserID2(),
		UserRelationshipStatus: ur.GetStatus().String(),
		CreatedAt:              ur.GetCreatedAt(),
		UpdatedAt:              ur.GetUpdatedAt(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return event.New("", jsonPayload), nil
}
