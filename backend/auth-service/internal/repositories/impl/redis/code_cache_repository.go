package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

const (
	codeKeyPrefix       = "code"
	codeKeyByUserPrefix = "code:user"
)

type CodeCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewCodeCacheRepository(redis *rediscl.RedisClient) interfaces.CodeCacheRepository {
	return &CodeCacheRepository{redis: redis}
}

func (r *CodeCacheRepository) keyByID(id int64) string {
	return fmt.Sprintf("%s:%d", codeKeyPrefix, id)
}

func (r *CodeCacheRepository) keyByUserID(userId string) string {
	return fmt.Sprintf("%s:%s", codeKeyByUserPrefix, userId)
}

func (r *CodeCacheRepository) GetCodeById(ctx context.Context, id int64) (*code.Code, error) {
	cached, err := get[models.V1CodeDal](ctx, r.redis, r.keyByID(id))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *CodeCacheRepository) SetCodeById(ctx context.Context, c *code.Code) error {
	ttl := time.Until(c.GetExpiresAt())
	if ttl <= 0 {
		return nil
	}
	return set(ctx, r.redis, r.keyByID(c.GetID()), models.V1CodeDalFromDomain(c), ttl)
}

func (r *CodeCacheRepository) DelCodeById(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyByID(id))
}

func (r *CodeCacheRepository) GetCodeByUserId(ctx context.Context, userId string) (*code.Code, error) {
	cached, err := get[models.V1CodeDal](ctx, r.redis, r.keyByUserID(userId))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *CodeCacheRepository) SetCodeByUserId(ctx context.Context, c *code.Code) error {
	ttl := time.Until(c.GetExpiresAt())
	if ttl <= 0 {
		return nil
	}
	return set(ctx, r.redis, r.keyByUserID(c.GetUserID()), models.V1CodeDalFromDomain(c), ttl)
}

func (r *CodeCacheRepository) DelCodeByUserId(ctx context.Context, userId string) error {
	return del(ctx, r.redis, r.keyByUserID(userId))
}
