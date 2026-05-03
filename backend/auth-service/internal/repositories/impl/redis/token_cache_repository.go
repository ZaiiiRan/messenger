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
	tokenKeyPrefix          = "token"
	tokenListKeyPrefix      = "token:list"
	tokenListIndexKeyPrefix = "token:list:index"
	tokenListTTL            = 1 * time.Minute
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

func (r *TokenCacheRepository) listKey(query *models.QueryTokensDal) (string, error) {
	hash, err := queryHash(query)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", tokenListKeyPrefix, hash), nil
}

func (r *TokenCacheRepository) GetTokenList(ctx context.Context, query *models.QueryTokensDal) ([]*token.Token, error) {
	key, err := r.listKey(query)
	if err != nil {
		return nil, err
	}
	cached, err := get[[]models.V1RefreshTokenDal](ctx, r.redis, key)
	if err != nil || cached == nil {
		return nil, err
	}
	result := make([]*token.Token, 0, len(*cached))
	for _, dal := range *cached {
		result = append(result, dal.ToDomain())
	}
	return result, nil
}

func (r *TokenCacheRepository) listIndexKey(userId string) string {
	return fmt.Sprintf("%s:%s", tokenListIndexKeyPrefix, userId)
}

func (r *TokenCacheRepository) SetTokenList(ctx context.Context, query *models.QueryTokensDal, tokens []*token.Token) error {
	key, err := r.listKey(query)
	if err != nil {
		return err
	}
	dals := make([]models.V1RefreshTokenDal, 0, len(tokens))
	for _, t := range tokens {
		dals = append(dals, models.V1RefreshTokenFromDomain(t))
	}
	if err := set(ctx, r.redis, key, dals, tokenListTTL); err != nil {
		return err
	}
	indexKey := r.listIndexKey(query.UserId)
	r.redis.GetClient().SAdd(ctx, indexKey, key)
	r.redis.GetClient().Expire(ctx, indexKey, tokenListTTL)
	return nil
}

func (r *TokenCacheRepository) DelTokenListsByUserId(ctx context.Context, userId string) error {
	indexKey := r.listIndexKey(userId)
	keys, err := r.redis.GetClient().SMembers(ctx, indexKey).Result()
	if err != nil || len(keys) == 0 {
		return nil
	}
	toDelete := append(keys, indexKey)
	return r.redis.GetClient().Del(ctx, toDelete...).Err()
}
