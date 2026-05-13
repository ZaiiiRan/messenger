package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	emailchangecode "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code/email_change_code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

type EmailChangeCodeCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewEmailChangeCodeCacheRepository(redis *rediscl.RedisClient) interfaces.EmailChangeCodeCacheRepository {
	return &EmailChangeCodeCacheRepository{redis: redis}
}

func (r *EmailChangeCodeCacheRepository) keyByID(id int64) string {
	return fmt.Sprintf("%s:%d", codeKeyPrefix, id)
}

func (r *EmailChangeCodeCacheRepository) keyByUserID(userId string) string {
	return fmt.Sprintf("%s:%s:%s", codeKeyByUserPrefix, userId, string(code.CodeTypeEmailChange))
}

func (r *EmailChangeCodeCacheRepository) GetCodeById(ctx context.Context, id int64) (*emailchangecode.EmailChangeCode, error) {
	// GetCodeById is not used with codeType context; callers that need typed domain
	// objects should use GetCodeByUserId instead.
	return nil, nil
}

func (r *EmailChangeCodeCacheRepository) SetCodeById(ctx context.Context, c *emailchangecode.EmailChangeCode) error {
	ttl := time.Until(c.GetExpiresAt())
	if ttl <= 0 {
		return nil
	}
	return set(ctx, r.redis, r.keyByID(c.GetID()), models.V1EmailCodeDalFromDomain(c), ttl)
}

func (r *EmailChangeCodeCacheRepository) DelCodeById(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyByID(id))
}

func (r *EmailChangeCodeCacheRepository) GetCodeByUserId(ctx context.Context, userId string) (*emailchangecode.EmailChangeCode, error) {
	cached, err := get[models.V1EmailCodeDal](ctx, r.redis, r.keyByUserID(userId))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *EmailChangeCodeCacheRepository) SetCodeByUserId(ctx context.Context, c *emailchangecode.EmailChangeCode) error {
	ttl := time.Until(c.GetExpiresAt())
	if ttl <= 0 {
		return nil
	}
	return set(ctx, r.redis, r.keyByUserID(c.GetUserID()), models.V1EmailCodeDalFromDomain(c), ttl)
}

func (r *EmailChangeCodeCacheRepository) DelCodeByUserId(ctx context.Context, userId string) error {
	return del(ctx, r.redis, r.keyByUserID(userId))
}
