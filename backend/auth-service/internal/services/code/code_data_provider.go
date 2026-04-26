package codeservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

type codeDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newCodeDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *codeDataProvider {
	return &codeDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (cdp *codeDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(cdp.pg)
}

func (cdp *codeDataProvider) getByUserID(ctx context.Context, userID string, uow *uow.UnitOfWork) (*code.Code, error) {
	cacheRepo := redisimpl.NewCodeCacheRepository(cdp.redis)
	code, err := cacheRepo.GetCodeByUserId(ctx, userID)
	if err == nil && code != nil {
		return code, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewCodeRepository(pgConn)
	query := models.NewQueryCodeDal(nil, &userID)
	code, err = dbRepo.QueryCode(ctx, query)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, nil
	}

	cacheRepo.SetCodeByUserId(ctx, code)
	return code, nil
}

func (cdp *codeDataProvider) save(ctx context.Context, code *code.Code, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewCodeRepository(pgConn)

	if code.GetID() == 0 {
		if err := dbRepo.CreateCode(ctx, code); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateCode(ctx, code); err != nil {
			return err
		}
	}

	cachRepo := redisimpl.NewCodeCacheRepository(cdp.redis)
	if err := cachRepo.SetCodeByUserId(ctx, code); err != nil {
		return err
	}

	return nil
}

func (cdp *codeDataProvider) delete(ctx context.Context, code *code.Code, uow *uow.UnitOfWork) error {
	cacheRepo := redisimpl.NewCodeCacheRepository(cdp.redis)
	if err := cacheRepo.DelCodeByUserId(ctx, code.GetUserID()); err != nil {
		return err
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewCodeRepository(pgConn)
	err = dbRepo.DeleteCode(ctx, code)

	return err
}
