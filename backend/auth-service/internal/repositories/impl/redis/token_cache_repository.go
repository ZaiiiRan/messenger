package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

const (
	tokenKeyPrefix = "token"
)

type TokenCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewTokenCacheRepository(redis *rediscl.RedisClient) interfaces.TokenCacheRepository {
	return &TokenCacheRepository{redis: redis}
}

func (r *TokenCacheRepository) key(tokenStr string) string {
	return fmt.Sprintf("%s:%s", tokenKeyPrefix, tokenStr)
}

func (r *TokenCacheRepository) GetToken(ctx context.Context, tokenStr string) (*token.Token, error) {
	cached, err := get[models.V1RefreshTokenDal](ctx, r.redis, r.key(tokenStr))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *TokenCacheRepository) SetToken(ctx context.Context, t *token.Token) error {
	ttl := time.Until(t.GetExpiresAt())
	if ttl <= 0 {
		return nil
	}
	return set(ctx, r.redis, r.key(t.GetToken()), models.V1RefreshTokenFromDomain(t), ttl)
}

func (r *TokenCacheRepository) DelToken(ctx context.Context, tokenStr string) error {
	return del(ctx, r.redis, r.key(tokenStr))
}
