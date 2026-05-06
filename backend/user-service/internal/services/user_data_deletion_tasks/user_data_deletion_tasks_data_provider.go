package userdatadeletiontasks

import (
	"context"

	outboxevent "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/outbox_event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/postgres"
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
	events []*outboxevent.OutboxEvent,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserDataDeletionTasksRepository(pgConn)
	err = dbRepo.Create(ctx, events)

	return err
}
