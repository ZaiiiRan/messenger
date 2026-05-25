package filedeletiontasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
)

type fileDeletionTasksDataProvider struct {
	pg *postgres.PostgresClient
}

func newFileDeletionTasksDataProvider(pg *postgres.PostgresClient) *fileDeletionTasksDataProvider {
	return &fileDeletionTasksDataProvider{
		pg: pg,
	}
}

func (fdp *fileDeletionTasksDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(fdp.pg)
}

func (fdp *fileDeletionTasksDataProvider) createFileDeletionTasks(
	ctx context.Context,
	tasks []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileDeletionTasksInboxRepository(pgConn)
	return dbRepo.CreateEvents(ctx, tasks)
}

func (fdp *fileDeletionTasksDataProvider) updateFileDeletionTasks(
	ctx context.Context,
	tasks []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileDeletionTasksInboxRepository(pgConn)
	return dbRepo.UpdateEvents(ctx, tasks)
}

func (fdp *fileDeletionTasksDataProvider) deleteFileDeletionTasks(
	ctx context.Context,
	tasks []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileDeletionTasksInboxRepository(pgConn)
	return dbRepo.DeleteEvents(ctx, tasks)
}

func (fdp *fileDeletionTasksDataProvider) getFileDeletionTasksLocked(
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

	dbRepo := postgresimpl.NewFileDeletionTasksInboxRepository(pgConn)
	return dbRepo.QueryEventsLocked(ctx, query)
}
