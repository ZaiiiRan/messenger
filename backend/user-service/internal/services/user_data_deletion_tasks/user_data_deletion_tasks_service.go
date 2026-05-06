package userdatadeletiontasks

import (
	"context"
	"encoding/json"

	outboxevent "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/outbox_event"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	producersmodels "github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/models"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserDataDeletionTasksService interface {
	CreateUserDataDeletionTasks(ctx context.Context, users []*user.User, uow *uow.UnitOfWork) error
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *userDataDeletionTasksDataProvider
}

func New(pgClient *postgres.PostgresClient, log *zap.SugaredLogger) UserDataDeletionTasksService {
	return &service{
		log:          log,
		dataProvider: newUserDataDeletionTasksDataProvider(pgClient),
	}
}

func (s *service) CreateUserDataDeletionTasks(ctx context.Context, users []*user.User, uow *uow.UnitOfWork) error {
	l := s.log.With("op", "create_user_data_deletion_tasks")

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

	return nil
}

func (s *service) createUserDataDeletionTask(user *user.User) (*outboxevent.OutboxEvent, error) {
	payload := producersmodels.UserDataDeletionTask{
		Id:          uuid.New().String(),
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
