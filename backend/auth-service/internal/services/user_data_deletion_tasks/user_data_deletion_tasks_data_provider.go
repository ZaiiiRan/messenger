package userdatadeletiontasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
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

	dbRepo := postgresimpl.NewUserDataDeletionTasksRepository(pgConn)
	err = dbRepo.CreateInboxEvents(ctx, events)

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

	dbRepo := postgresimpl.NewUserDataDeletionTasksRepository(pgConn)
	err = dbRepo.UpdateInboxEvents(ctx, events)

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

	dbRepo := postgresimpl.NewUserDataDeletionTasksRepository(pgConn)
	err = dbRepo.DeleteInboxEvents(ctx, events)

	return err
}

func (udp *userDataDeletionTasksDataProvider) getUserDataDeletionTasksLocked(
	ctx context.Context,
	batch_size int,
	retryAfter time.Time,
	createdAfter *time.Time,
	uow *uow.UnitOfWork,
) ([]*event.Event, error) {
	query := models.NewQueryInboxEventsLockedDal(retryAfter, createdAfter, batch_size)

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksRepository(pgConn)
	events, err := dbRepo.QueryInboxEventsLocked(ctx, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}
