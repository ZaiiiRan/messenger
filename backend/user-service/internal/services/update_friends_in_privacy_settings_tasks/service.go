package updatefriendsinprivacysettingstasks

import (
	"context"
	"encoding/json"

	consumersmodels "github.com/ZaiiiRan/messenger/backend/user-service/internal/consumers/models"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type UpdateFriendsInPrivacySettingsTasksService interface {
	CreateUpdateFriendsInPrivacySettingsTasks(
		ctx context.Context, workerID string, taskMessages []*consumersmodels.UserRelationshipChangeTask, uow *uow.UnitOfWork,
	) error
	// ProcessUpdateFriendsInPrivacySettingsTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int, uow *uow.UnitOfWork) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *updateFriendsInPrivacySettingsTasksDataProvider
}

func New(
	pgClient *postgres.PostgresClient,
	log *zap.SugaredLogger,
) UpdateFriendsInPrivacySettingsTasksService {
	return &service{
		log:          log,
		dataProvider: newUpdateFriendsInPrivacySettingsTasksDataProvider(pgClient),
	}
}

func (s *service) CreateUpdateFriendsInPrivacySettingsTasks(
	ctx context.Context, workerID string, taskMessages []*consumersmodels.UserRelationshipChangeTask, uow *uow.UnitOfWork,
) error {
	l := s.log.With("op", "create_update_friends_in_privacy_settings_tasks", "worker_id", workerID)

	inboxEvents := make([]*event.Event, 0, len(taskMessages))
	inboxEventsByID := make(map[string]*event.Event)
	for _, taskMessage := range taskMessages {
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

func (s *service) createUpdateFriendsInPrivacySettingsTask(taskMessage *consumersmodels.UserRelationshipChangeTask) (*event.Event, error) {
	jsonPayload, err := json.Marshal(taskMessage)
	if err != nil {
		return nil, err
	}
	return event.New(taskMessage.Id, jsonPayload), nil
}
