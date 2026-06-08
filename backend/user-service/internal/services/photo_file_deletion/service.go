package photofiledeletionsvc

import (
	"context"
	"encoding/json"
	"time"

	filepb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/file/v1"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	filegrpcclient "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/client/grpc/file_client"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type photoFileDeletionPayload struct {
	FileID string `json:"file_id"`
}

type PhotoFileDeletionService interface {
	CreateTasks(ctx context.Context, fileIDs []string, u *uow.UnitOfWork) error
	SendTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error
}

type service struct {
	pg         *postgres.PostgresClient
	fileClient *filegrpcclient.Client
	log        *zap.SugaredLogger
}

func New(pg *postgres.PostgresClient, fileClient *filegrpcclient.Client, log *zap.SugaredLogger) PhotoFileDeletionService {
	return &service{pg: pg, fileClient: fileClient, log: log}
}

func (s *service) CreateTasks(ctx context.Context, fileIDs []string, u *uow.UnitOfWork) error {
	if len(fileIDs) == 0 {
		return nil
	}

	needToCommit := u == nil
	if needToCommit {
		u = uow.New(s.pg)
		defer u.Close()
		if _, err := u.BeginTransaction(ctx); err != nil {
			return ErrCreateTasks
		}
	}

	conn, err := u.GetConn(ctx)
	if err != nil {
		return ErrCreateTasks
	}

	events := make([]*event.Event, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		payload, err := json.Marshal(photoFileDeletionPayload{FileID: fileID})
		if err != nil {
			return ErrMarshalPayload
		}
		events = append(events, event.New(uuid.NewString(), payload))
	}

	if err := postgresimpl.NewPhotoFileDeletionOutboxRepository(conn).Create(ctx, events, false); err != nil {
		return ErrCreateTasks
	}

	if needToCommit {
		if err := u.Commit(ctx); err != nil {
			return ErrCreateTasks
		}
	}
	return nil
}

func (s *service) SendTasks(ctx context.Context, workerID string, retryIntervalMS uint, batchSize int) error {
	l := s.log.With("op", "photo_file_deletion.send", "worker_id", workerID)

	u := uow.New(s.pg)
	defer u.Close()

	if _, err := u.BeginTransaction(ctx); err != nil {
		l.Errorw("photo_file_deletion.send_failed.begin_tx", "err", err)
		return ErrSendTasks
	}

	conn, err := u.GetConn(ctx)
	if err != nil {
		return ErrSendTasks
	}

	retryAfter := time.Now().Add(-time.Duration(retryIntervalMS) * time.Millisecond)
	repo := postgresimpl.NewPhotoFileDeletionOutboxRepository(conn)

	events, err := repo.QueryLocked(ctx, models.NewQueryEventsLockedDal(retryAfter, batchSize))
	if err != nil {
		l.Errorw("photo_file_deletion.send_failed.query_locked", "err", err)
		return ErrSendTasks
	}
	if len(events) == 0 {
		return nil
	}

	now := time.Now()
	for _, e := range events {
		var payload photoFileDeletionPayload
		if err := json.Unmarshal(e.GetPayload(), &payload); err != nil {
			l.Warnw("photo_file_deletion.send.unmarshal_failed", "err", err)
			_ = e.SetStatus(event.EventStatusFailed)
			e.SetUpdatedAt(&now)
			continue
		}

		_, callErr := s.fileClient.FileClient().DeleteFiles(ctx, &filepb.DeleteFilesRequest{
			FileIds: []string{payload.FileID},
		})

		_ = e.IncrementAttempts()
		if callErr != nil {
			l.Warnw("photo_file_deletion.send.delete_failed", "file_id", payload.FileID, "err", callErr)
			_ = e.SetStatus(event.EventStatusFailed)
		} else {
			_ = e.SetStatus(event.EventStatusSuccess)
		}
		e.SetUpdatedAt(&now)
	}

	if err := repo.Update(ctx, events); err != nil {
		l.Errorw("photo_file_deletion.send_failed.update", "err", err)
		return ErrSendTasks
	}

	if err := u.Commit(ctx); err != nil {
		l.Errorw("photo_file_deletion.send_failed.commit", "err", err)
		return ErrSendTasks
	}
	return nil
}
