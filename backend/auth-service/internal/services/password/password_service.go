package passwordservice

import (
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"go.uber.org/zap"
)

type PasswordService interface {
}

type passwordService struct {
	passwordDataProvider *passwordDataProvider
	log                  *zap.SugaredLogger
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) PasswordService {
	return &passwordService{
		passwordDataProvider: newPasswordDataProvider(pgClient, redisClient),
		log:                  log,
	}
}
