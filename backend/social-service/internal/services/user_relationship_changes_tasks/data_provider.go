package userrelationshipchangestasks

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
)

type userRelationshipChangesTasksDataProvider struct {
	pg *postgres.PostgresClient
}

func newUserRelationshipChangesTasksDataProvider(pg *postgres.PostgresClient) *userRelationshipChangesTasksDataProvider {
	return &userRelationshipChangesTasksDataProvider{
		pg: pg,
	}
}

func (udp *userRelationshipChangesTasksDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(udp.pg)
}

func (udp *userRelationshipChangesTasksDataProvider) createUserRelationshipChangesTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipChangesTasksOutboxRepository(pgConn)
	return dbRepo.CreateEvents(ctx, events)
}

func (udp *userRelationshipChangesTasksDataProvider) updatetUserRelationshipChangesTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipChangesTasksOutboxRepository(pgConn)
	return dbRepo.UpdateEvents(ctx, events)
}

func (udp *userRelationshipChangesTasksDataProvider) deleteUserRelationshipChangesTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipChangesTasksOutboxRepository(pgConn)
	return dbRepo.DeleteEvents(ctx, events)
}

func (udp *userRelationshipChangesTasksDataProvider) getUserRelationshipChangesTasksLocked(
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

	dbRepo := postgresimpl.NewUserRelationshipChangesTasksOutboxRepository(pgConn)
	return dbRepo.QueryEventsLocked(ctx, query)
}
