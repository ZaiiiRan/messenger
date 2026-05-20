package tokenservice

import (
	"context"

	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

type userVersionDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newUserVersionDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *userVersionDataProvider {
	return &userVersionDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (uvdp *userVersionDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(uvdp.pg)
}

func (uvdp *userVersionDataProvider) getByUserId(ctx context.Context, userID string, uow *uow.UnitOfWork) (*userversion.UserVersion, error) {
	cacheRepo := redisimpl.NewUserVersionCacheRepository(uvdp.redis)
	uv, err := cacheRepo.GetByUserId(ctx, userID)
	if err == nil && uv != nil {
		return uv, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserVersionRepository(pgConn)
	query := models.NewQueryUserVersionDal(nil, &userID)
	uv, err = dbRepo.QueryUserVersion(ctx, query)
	if err != nil {
		return nil, err
	}
	if uv == nil {
		return nil, nil
	}

	cacheRepo.SetByUserId(ctx, uv)
	return uv, nil
}

func (uvdp *userVersionDataProvider) save(ctx context.Context, uv *userversion.UserVersion, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserVersionRepository(pgConn)

	if uv.GetID() == 0 {
		if err := dbRepo.CreateUserVersion(ctx, uv); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateUserVersion(ctx, uv); err != nil {
			return err
		}
	}

	uow.OnCommit(func() {
		redisimpl.NewUserVersionCacheRepository(uvdp.redis).SetByUserId(ctx, uv)
	})
	return nil
}

func (uvdp *userVersionDataProvider) delete(ctx context.Context, uv *userversion.UserVersion, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserVersionRepository(pgConn)
	if err := dbRepo.DeleteUserVersion(ctx, uv); err != nil {
		return err
	}

	userID := uv.GetUserID()
	uow.OnCommit(func() {
		redisimpl.NewUserVersionCacheRepository(uvdp.redis).DelByUserId(ctx, userID)
	})
	return nil
}
