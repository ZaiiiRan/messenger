package updatefriendsinprivacysettingstasksservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
)

type updateFriendsInPrivacySettingsTasksDataProvider struct {
	pg *postgres.PostgresClient
}

func newUpdateFriendsInPrivacySettingsTasksDataProvider(pg *postgres.PostgresClient) *updateFriendsInPrivacySettingsTasksDataProvider {
	return &updateFriendsInPrivacySettingsTasksDataProvider{
		pg: pg,
	}
}

func (udp *updateFriendsInPrivacySettingsTasksDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(udp.pg)
}

func (udp *updateFriendsInPrivacySettingsTasksDataProvider) createUpdateFriendsInPrivacySettingsTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUpdateFriendsInPrivacySettingsTasksInboxRepository(pgConn)
	return dbRepo.Create(ctx, events, true)
}

func (udp *updateFriendsInPrivacySettingsTasksDataProvider) updateUpdateFriendsInPrivacySettingsTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUpdateFriendsInPrivacySettingsTasksInboxRepository(pgConn)
	return dbRepo.Update(ctx, events)
}

func (udp *updateFriendsInPrivacySettingsTasksDataProvider) deleteUpdateFriendsInPrivacySettingsTasks(
	ctx context.Context,
	events []*event.Event,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUpdateFriendsInPrivacySettingsTasksInboxRepository(pgConn)
	return dbRepo.Delete(ctx, events)
}

func (udp *updateFriendsInPrivacySettingsTasksDataProvider) getUpdateFriendsInPrivacySettingsTasksLocked(
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

	dbRepo := postgresimpl.NewUpdateFriendsInPrivacySettingsTasksInboxRepository(pgConn)
	events, err := dbRepo.QueryLocked(ctx, query)
	if err != nil {
		return nil, err
	}
	return events, nil
}
