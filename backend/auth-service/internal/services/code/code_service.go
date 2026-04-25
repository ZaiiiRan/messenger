package codeservice

import (
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"go.uber.org/zap"
)

type CodeService interface {
}

type codeService struct {
	codeDataProvider *codeDataProvider
	log              *zap.SugaredLogger
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) CodeService {
	return &codeService{
		codeDataProvider: newCodeDataProvider(pgClient, redisClient),
		log:              log,
	}
}
