package redisimpl

import (
	"context"
	"fmt"
	"time"

	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

const (
	userVersionKeyByUserPrefix = "user_version:user"
	userVersionTTL             = 10 * time.Minute
)

type UserVersionCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewUserVersionCacheRepository(redis *rediscl.RedisClient) interfaces.UserVersionCacheRepository {
	return &UserVersionCacheRepository{redis: redis}
}

func (r *UserVersionCacheRepository) keyByUserID(userId string) string {
	return fmt.Sprintf("%s:%s", userVersionKeyByUserPrefix, userId)
}

func (r *UserVersionCacheRepository) GetByUserId(ctx context.Context, userId string) (*userversion.UserVersion, error) {
	cached, err := get[models.V1UserVersionDal](ctx, r.redis, r.keyByUserID(userId))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *UserVersionCacheRepository) SetByUserId(ctx context.Context, uv *userversion.UserVersion) error {
	return set(ctx, r.redis, r.keyByUserID(uv.GetUserID()), models.V1UserVersionDalFromDomain(uv), userVersionTTL)
}

func (r *UserVersionCacheRepository) DelByUserId(ctx context.Context, userId string) error {
	return del(ctx, r.redis, r.keyByUserID(userId))
}
