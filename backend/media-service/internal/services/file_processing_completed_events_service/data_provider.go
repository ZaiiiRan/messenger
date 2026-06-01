package fileprocessingcompletedeventsservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/postgres"
)

type fileProcessingCompletedEventsDataProvider struct {
	pg *postgres.PostgresClient
}

func newFileProcessingCompletedEventsDataProvider(pg *postgres.PostgresClient) *fileProcessingCompletedEventsDataProvider {
	return &fileProcessingCompletedEventsDataProvider{
		pg: pg,
	}
}

func (dp *fileProcessingCompletedEventsDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(dp.pg)
}

func (dp *fileProcessingCompletedEventsDataProvider) createEvents(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileProcessingCompletedOutboxRepository(pgConn)
	return dbRepo.CreateEvents(ctx, events)
}

func (dp *fileProcessingCompletedEventsDataProvider) updateEvents(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileProcessingCompletedOutboxRepository(pgConn)
	return dbRepo.UpdateEvents(ctx, events)
}

func (dp *fileProcessingCompletedEventsDataProvider) deleteEvents(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileProcessingCompletedOutboxRepository(pgConn)
	return dbRepo.DeleteEvents(ctx, events)
}

func (dp *fileProcessingCompletedEventsDataProvider) getEventsLocked(
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

	dbRepo := postgresimpl.NewFileProcessingCompletedOutboxRepository(pgConn)
	return dbRepo.QueryEventsLocked(ctx, query)
}
