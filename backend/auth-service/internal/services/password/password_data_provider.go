package passwordservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

type passwordDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newPasswordDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *passwordDataProvider {
	return &passwordDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (pdp *passwordDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(pdp.pg)
}

func (pdp *passwordDataProvider) getByUserID(ctx context.Context, userID string, uow *uow.UnitOfWork) (*password.Password, error) {
	cacheRepo := redisimpl.NewPasswordCacheRepository(pdp.redis)
	password, err := cacheRepo.GetPasswordByUserId(ctx, userID)
	if err == nil && password != nil {
		return password, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewPasswordRepository(pgConn)
	query := models.NewQueryPasswordDal(nil, &userID)
	password, err = dbRepo.QueryPassword(ctx, query)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return nil, nil
	}

	cacheRepo.SetPasswordByUserId(ctx, password)
	return password, nil
}

func (pdp *passwordDataProvider) save(ctx context.Context, p *password.Password, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewPasswordRepository(pgConn)

	if p.GetID() == 0 {
		if err := dbRepo.CreatePassword(ctx, p); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdatePassword(ctx, p); err != nil {
			return err
		}
	}

	cacheRepo := redisimpl.NewPasswordCacheRepository(pdp.redis)
	if err := cacheRepo.SetPasswordByUserId(ctx, p); err != nil {
		return err
	}

	return nil
}
