package codeservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	pgclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	redisclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

type codeDataProvider struct {
	pg    *pgclient.PostgresClient
	redis *redisclient.RedisClient
}

func newCodeDataProvider(pg *pgclient.PostgresClient, redis *redisclient.RedisClient) *codeDataProvider {
	return &codeDataProvider{pg: pg, redis: redis}
}

func (cdp *codeDataProvider) getByUserID(ctx context.Context, userID string, codeType code.CodeType, uow *uow.UnitOfWork) (*code.Code, error) {
	cacheRepo := redisimpl.NewCodeCacheRepository(cdp.redis)
	if c, err := cacheRepo.GetCodeByUserId(ctx, userID, codeType); err == nil && c != nil {
		return c, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	var dbRepo interfaces.CodeRepository
	if codeType == code.CodeTypePasswordReset {
		dbRepo = postgresimpl.NewPasswordResetCodeRepository(pgConn)
	} else {
		dbRepo = postgresimpl.NewActivationCodeRepository(pgConn)
	}

	c, err := dbRepo.QueryCode(ctx, models.NewQueryCodeDal(nil, &userID))
	if err != nil {
		return nil, err
	}
	if c != nil {
		cacheRepo.SetCodeByUserId(ctx, c)
	}
	return c, nil
}

func (cdp *codeDataProvider) getByUserIDLocked(ctx context.Context, userID string, codeType code.CodeType, uow *uow.UnitOfWork) (*code.Code, error) {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	var dbRepo interfaces.CodeRepository
	if codeType == code.CodeTypePasswordReset {
		dbRepo = postgresimpl.NewPasswordResetCodeRepository(pgConn)
	} else {
		dbRepo = postgresimpl.NewActivationCodeRepository(pgConn)
	}

	q := models.NewQueryCodeDal(nil, &userID)
	q.ForUpdate = true
	return dbRepo.QueryCode(ctx, q)
}

func (cdp *codeDataProvider) getByLinkTokenLocked(ctx context.Context, linkToken string, codeType code.CodeType, uow *uow.UnitOfWork) (*code.Code, error) {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	var dbRepo interfaces.CodeRepository
	if codeType == code.CodeTypePasswordReset {
		dbRepo = postgresimpl.NewPasswordResetCodeRepository(pgConn)
	} else {
		dbRepo = postgresimpl.NewActivationCodeRepository(pgConn)
	}

	q := models.NewQueryCodeDal(nil, nil)
	q.LinkToken = &linkToken
	q.ForUpdate = true
	return dbRepo.QueryCode(ctx, q)
}

func (cdp *codeDataProvider) save(ctx context.Context, c *code.Code, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var dbRepo interfaces.CodeRepository
	if c.GetCodeType() == code.CodeTypePasswordReset {
		dbRepo = postgresimpl.NewPasswordResetCodeRepository(pgConn)
	} else {
		dbRepo = postgresimpl.NewActivationCodeRepository(pgConn)
	}

	if c.GetID() == 0 {
		if err := dbRepo.CreateCode(ctx, c); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateCode(ctx, c); err != nil {
			return err
		}
	}

	cacheRepo := redisimpl.NewCodeCacheRepository(cdp.redis)
	cacheRepo.SetCodeByUserId(ctx, c)
	return nil
}

func (cdp *codeDataProvider) delete(ctx context.Context, c *code.Code, uow *uow.UnitOfWork) error {
	cacheRepo := redisimpl.NewCodeCacheRepository(cdp.redis)
	cacheRepo.DelCodeByUserId(ctx, c.GetUserID(), c.GetCodeType())

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	var dbRepo interfaces.CodeRepository
	if c.GetCodeType() == code.CodeTypePasswordReset {
		dbRepo = postgresimpl.NewPasswordResetCodeRepository(pgConn)
	} else {
		dbRepo = postgresimpl.NewActivationCodeRepository(pgConn)
	}

	return dbRepo.DeleteCode(ctx, c)
}
