package filedeletionfroms3tasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
)

type fileDeletionFromS3TasksDataProvider struct {
	pg *postgres.PostgresClient
}

func newFileDeletionFromS3TasksDataProvider(pg *postgres.PostgresClient) *fileDeletionFromS3TasksDataProvider {
	return &fileDeletionFromS3TasksDataProvider{
		pg: pg,
	}
}

func (fdp *fileDeletionFromS3TasksDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(fdp.pg)
}

func (fdp *fileDeletionFromS3TasksDataProvider) createFileDeletionFromS3Tasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileDeletionFromS3TasksRepository(pgConn)
	return dbRepo.CreateEvents(ctx, events)
}

func (fdp *fileDeletionFromS3TasksDataProvider) updateFileDeletionFromS3Tasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileDeletionFromS3TasksRepository(pgConn)
	return dbRepo.UpdateEvents(ctx, events)
}

func (fdp *fileDeletionFromS3TasksDataProvider) deleteFileDeletionFromS3Tasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileDeletionFromS3TasksRepository(pgConn)
	return dbRepo.DeleteEvents(ctx, events)
}

func (fdp *fileDeletionFromS3TasksDataProvider) getFileDeletionFromS3TasksLocked(
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

	dbRepo := postgresimpl.NewFileDeletionFromS3TasksRepository(pgConn)
	return dbRepo.QueryEventsLocked(ctx, query)
}
