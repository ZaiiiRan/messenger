package userdatadeletiontasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
)

type userDataDeletionTasksDataProvider struct {
	pg *postgres.PostgresClient
}

func newUserDataDeletionTasksDataProvider(pg *postgres.PostgresClient) *userDataDeletionTasksDataProvider {
	return &userDataDeletionTasksDataProvider{
		pg: pg,
	}
}

func (udp *userDataDeletionTasksDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(udp.pg)
}

func (udp *userDataDeletionTasksDataProvider) createUserDataDeletionTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksInboxRepository(pgConn)
	return dbRepo.CreateEvents(ctx, events)
}

func (udp *userDataDeletionTasksDataProvider) updatetUserDataDeletionTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksInboxRepository(pgConn)
	return dbRepo.UpdateEvents(ctx, events)
}

func (udp *userDataDeletionTasksDataProvider) deleteUserDataDeletionTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksInboxRepository(pgConn)
	return dbRepo.DeleteEvents(ctx, events)
}

func (udp *userDataDeletionTasksDataProvider) getUserDataDeletionTasksLocked(
	ctx context.Context,
	batchSize int,
	retryAfter time.Time,
	uow *uow.UnitOfWork,
) ([]*event.Event, error) {
	query := models.NewQueryEventsLockedDal(retryAfter, batchSize)

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksInboxRepository(pgConn)
	return dbRepo.QueryEventsLocked(ctx, query)
}
