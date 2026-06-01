package fileprocessingtasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/media-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/postgres"
)

type fileProcessingTasksDataProvider struct {
	pg *postgres.PostgresClient
}

func newFileProcessingTasksDataProvider(pg *postgres.PostgresClient) *fileProcessingTasksDataProvider {
	return &fileProcessingTasksDataProvider{
		pg: pg,
	}
}

func (dp *fileProcessingTasksDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(dp.pg)
}

func (dp *fileProcessingTasksDataProvider) createTasks(
	ctx context.Context,
	tasks []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileProcessingTasksRepository(pgConn)
	return dbRepo.CreateEvents(ctx, tasks)
}

func (dp *fileProcessingTasksDataProvider) updateTasks(
	ctx context.Context,
	tasks []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileProcessingTasksRepository(pgConn)
	return dbRepo.UpdateEvents(ctx, tasks)
}

func (dp *fileProcessingTasksDataProvider) deleteTasks(
	ctx context.Context,
	tasks []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileProcessingTasksRepository(pgConn)
	return dbRepo.DeleteEvents(ctx, tasks)
}

func (dp *fileProcessingTasksDataProvider) getTasksLocked(
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

	dbRepo := postgresimpl.NewFileProcessingTasksRepository(pgConn)
	return dbRepo.QueryEventsLocked(ctx, query)
}
