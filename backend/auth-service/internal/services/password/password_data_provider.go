package passwordservice

import (
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
