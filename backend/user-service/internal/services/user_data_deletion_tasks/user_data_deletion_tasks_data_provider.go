package userdatadeletiontasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
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

	dbRepo := postgresimpl.NewUserDataDeletionTasksOutboxRepository(pgConn)
	err = dbRepo.Create(ctx, events)

	return err
}

func (udp *userDataDeletionTasksDataProvider) updateUserDataDeletionTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksOutboxRepository(pgConn)
	err = dbRepo.Update(ctx, events)

	return err
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

	dbRepo := postgresimpl.NewUserDataDeletionTasksOutboxRepository(pgConn)
	err = dbRepo.Delete(ctx, events)

	return err
}

func (udp *userDataDeletionTasksDataProvider) getUserDataDeletionTasksLocked(
	ctx context.Context,
	batch_size int,
	retryAfter time.Time,
	uow *uow.UnitOfWork,
) ([]*event.Event, error) {
	query := models.NewQueryEventsLockedDal(retryAfter, batch_size)

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksOutboxRepository(pgConn)
	events, err := dbRepo.QueryLocked(ctx, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}
