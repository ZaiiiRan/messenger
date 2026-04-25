package tokenservice

import (
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
