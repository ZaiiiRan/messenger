package tokenservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"

	postgresimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

type tokenDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newTokenDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *tokenDataProvider {
	return &tokenDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (tdp *tokenDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(tdp.pg)
}

func (tdp *tokenDataProvider) getByToken(ctx context.Context, token string, uow *uow.UnitOfWork) (*token.Token, error) {
	cacheRepo := redisimpl.NewTokenCacheRepository(tdp.redis)
	t, err := cacheRepo.GetToken(ctx, token)
	if err == nil && t != nil {
		return t, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewTokenRepository(pgConn)
	query := models.NewQueryTokenDal(nil, nil, &token, nil)
	t, err = dbRepo.QueryToken(ctx, query)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}

	cacheRepo.SetToken(ctx, t)
	return t, nil
}

func (tdp *tokenDataProvider) save(ctx context.Context, t *token.Token, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewTokenRepository(pgConn)

	if t.GetID() == 0 {
		if err := dbRepo.CreateToken(ctx, t); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateToken(ctx, t); err != nil {
			return err
		}
	}

	cacheRepo := redisimpl.NewTokenCacheRepository(tdp.redis)
	cacheRepo.SetToken(ctx, t)

	return nil
}

func (tdp *tokenDataProvider) delete(ctx context.Context, token string, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewTokenRepository(pgConn)

	err = dbRepo.DeleteToken(ctx, token)
	if err != nil {
		return err
	}

	cacheRepo := redisimpl.NewTokenCacheRepository(tdp.redis)
	cacheRepo.DelToken(ctx, token)
	return nil
}

func (tdp *tokenDataProvider) deleteFromCache(ctx context.Context, token string) error {
	cacheRepo := redisimpl.NewTokenCacheRepository(tdp.redis)
	return cacheRepo.DelToken(ctx, token)
}

func (tdp *tokenDataProvider) deleteExpiredTokens(ctx context.Context, batchSize uint, uow *uow.UnitOfWork) ([]*token.Token, error) {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewTokenRepository(pgConn)

	tokens, err := dbRepo.DeleteExpiredTokens(ctx, models.NewQueryExpiredTokensDal(int(batchSize)))
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
