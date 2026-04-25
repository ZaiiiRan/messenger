package tokenservice

import (
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"go.uber.org/zap"
)

type TokenService interface {
}

type tokenService struct {
	tokenDataProvider *tokenDataProvider
	log              *zap.SugaredLogger
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) TokenService {
	return &tokenService{
		tokenDataProvider: newTokenDataProvider(pgClient, redisClient),
		log:              log,
	}
}
