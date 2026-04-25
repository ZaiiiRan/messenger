package codeservice

import (
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
