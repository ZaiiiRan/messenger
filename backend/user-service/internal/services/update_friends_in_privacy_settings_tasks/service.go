package updatefriendsinprivacysettingstasksservice

import (
	"context"
	"encoding/json"
	"time"

	consumersmodels "github.com/ZaiiiRan/messenger/backend/user-service/internal/consumers/models"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type UpdateFriendsInPrivacySettingsTasksService interface {
	CreateUpdateFriendsInPrivacySettingsTasks(
		ctx context.Context, workerID string, taskMessages []*consumersmodels.UserRelationshipChangeTask, uow *uow.UnitOfWork,
	) error
	ProcessUpdateFriendsInPrivacySettingsTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *updateFriendsInPrivacySettingsTasksDataProvider
	userService  userservice.UserService
}

func New(
	pgClient *postgres.PostgresClient,
	userService userservice.UserService,
	log *zap.SugaredLogger,
) UpdateFriendsInPrivacySettingsTasksService {
	return &service{
		log:          log,
		dataProvider: newUpdateFriendsInPrivacySettingsTasksDataProvider(pgClient),
		userService:  userService,
	}
}

func (s *service) CreateUpdateFriendsInPrivacySettingsTasks(
	ctx context.Context, workerID string, taskMessages []*consumersmodels.UserRelationshipChangeTask, uow *uow.UnitOfWork,
) error {
	l := s.log.With("op", "create_update_friends_in_privacy_settings_tasks", "worker_id", workerID)

	inboxEvents := make([]*event.Event, 0, len(taskMessages))
	inboxEventsByID := make(map[string]*event.Event)
	for _, taskMessage := range taskMessages {
		if taskMessage.Action != RemoveFromFriendsActionType && taskMessage.Action != BlockActionType {
			continue
		}
		evt, err := s.createUpdateFriendsInPrivacySettingsTask(taskMessage)
		if err != nil {
			l.Warnw(
				"update_friends_in_privacy_settings_tasks.create_update_friends_in_privacy_settings_tasks_failed",
				"err", err,
				"task_message", taskMessage,
			)
			continue
		}
		if value, ok := inboxEventsByID[evt.GetID()]; ok {
			l.Warnw(
				"update_friends_in_privacy_settings_tasks.create_update_friends_in_privacy_settings_tasks_failed",
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
			l.Errorw("update_friends_in_privacy_settings_tasks.create_update_friends_in_privacy_settings_tasks_failed.begin_transaction_error", "err", err)
			return ErrCreateUpdateFriendsInPrivacySettingsTasksError
		}
	}

	if err := s.dataProvider.createUpdateFriendsInPrivacySettingsTasks(ctx, inboxEvents, uow); err != nil {
		l.Errorw("update_friends_in_privacy_settings_tasks.create_update_friends_in_privacy_settings_tasks_failed.create_error", "err", err)
		return ErrCreateUpdateFriendsInPrivacySettingsTasksError
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("update_friends_in_privacy_settings_tasks.create_update_friends_in_privacy_settings_tasks_failed.commit_error", "err", err)
			return ErrCreateUpdateFriendsInPrivacySettingsTasksError
		}
	}

	l.Infow("update_friends_in_privacy_settings_tasks.create_update_friends_in_privacy_settings_tasks.success", "count", len(inboxEvents))
	return nil
}

func (s *service) ProcessUpdateFriendsInPrivacySettingsTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "process_update_friends_in_privacy_settings_tasks", "worker_id", workerID)

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		_, err := uow.BeginTransaction(ctx)
		if err != nil {
			l.Errorw("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.begin_transaction_error", "err", err)
			return ErrProcessUpdateFriendsInPrivacySettingsTasksError
		}
		needToCommit = true
	}

	now := time.Now()
	retryAfter := now.Add(-1 * (time.Duration(retryIntervalMS) * time.Millisecond))
	inboxEvents, err := s.dataProvider.getUpdateFriendsInPrivacySettingsTasksLocked(ctx, batchSize, retryAfter, uow)
	if err != nil {
		l.Errorw("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.get_error", "err", err)
		return ErrProcessUpdateFriendsInPrivacySettingsTasksError
	}
	if inboxEvents == nil {
		return nil
	}

	inboxEventsSuccess := make([]*event.Event, 0, len(inboxEvents))
	inboxEventsFailed := make([]*event.Event, 0, len(inboxEvents))

	for _, evt := range inboxEvents {
		var taskMessage *consumersmodels.UserRelationshipChangeTask
		if err := json.Unmarshal(evt.GetPayload(), &taskMessage); err != nil {
			l.Errorw(
				"update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.unmarshal_payload_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markUserDataDeletionTaskFailed(evt, now, l, "update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed")
			if err != nil {
				continue
			}
			inboxEventsFailed = append(inboxEventsFailed, evt)
			continue
		}

		if err := s.processUpdateFriendsInPrivacySettingsTask(ctx, workerID, taskMessage); err != nil {
			l.Errorw(
				"update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.process_task_error",
				"err", err,
				"event", evt.GetID(),
				"payload", evt.GetPayload(),
				"attempts", evt.GetAttempts(),
				"status", evt.GetStatus(),
			)
			err = s.markUserDataDeletionTaskFailed(evt, now, l, "update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed")
			if err != nil {
				continue
			}
			inboxEventsFailed = append(inboxEventsFailed, evt)
			continue
		}

		inboxEventsSuccess = append(inboxEventsSuccess, evt)
	}

	if len(inboxEventsSuccess) > 0 {
		if err := s.dataProvider.deleteUpdateFriendsInPrivacySettingsTasks(ctx, inboxEventsSuccess, uow); err != nil {
			l.Errorw("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.delete_error", "err", err)
			return ErrProcessUpdateFriendsInPrivacySettingsTasksError
		}
	}
	if len(inboxEventsFailed) > 0 {
		if err := s.dataProvider.updateUpdateFriendsInPrivacySettingsTasks(ctx, inboxEventsFailed, uow); err != nil {
			l.Errorw("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.update_error", "err", err)
			return ErrProcessUpdateFriendsInPrivacySettingsTasksError
		}
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks_failed.commit_error", "err", err)
			return ErrProcessUpdateFriendsInPrivacySettingsTasksError
		}
	}

	if len(inboxEventsSuccess) > 0 {
		l.Infow("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks.success", "successfully_processed", len(inboxEventsSuccess))
	}
	if len(inboxEventsFailed) > 0 {
		l.Warnw("update_friends_in_privacy_settings_tasks.process_update_friends_in_privacy_settings_tasks.success", "failed", len(inboxEventsFailed))
	}

	return nil
}

func (s *service) processUpdateFriendsInPrivacySettingsTask(
	ctx context.Context,
	workerID string,
	taskMessage *consumersmodels.UserRelationshipChangeTask,
) error {
	if err := s.userService.DeleteUserFromPrivacySettingsLists(ctx, taskMessage.User1Id, taskMessage.User2Id, workerID, nil); err != nil {
		return err
	}
	if err := s.userService.DeleteUserFromPrivacySettingsLists(ctx, taskMessage.User2Id, taskMessage.User1Id, workerID, nil); err != nil {
		return err
	}
	return nil
}

func (s *service) markUserDataDeletionTaskFailed(
	evt *event.Event,
	now time.Time,
	log *zap.SugaredLogger,
	logPrefix string,
) error {
	err := evt.IncrementAttempts()
	if err != nil {
		log.Errorw(
			logPrefix+".mark_task_as_failed_error",
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

func (s *service) createUpdateFriendsInPrivacySettingsTask(taskMessage *consumersmodels.UserRelationshipChangeTask) (*event.Event, error) {
	jsonPayload, err := json.Marshal(taskMessage)
	if err != nil {
		return nil, err
	}
	return event.New(taskMessage.Id, jsonPayload), nil
}
