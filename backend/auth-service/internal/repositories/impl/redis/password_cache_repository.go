package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
)

const (
	passwordKeyPrefix       = "password"
	passwordKeyByUserPrefix = "password:user"
	passwordTTL             = 10 * time.Minute
)

type PasswordCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewPasswordCacheRepository(redis *rediscl.RedisClient) interfaces.PasswordCacheRepository {
	return &PasswordCacheRepository{redis: redis}
}

func (r *PasswordCacheRepository) keyByID(id int64) string {
	return fmt.Sprintf("%s:%d", passwordKeyPrefix, id)
}

func (r *PasswordCacheRepository) keyByUserID(userId string) string {
	return fmt.Sprintf("%s:%s", passwordKeyByUserPrefix, userId)
}

func (r *PasswordCacheRepository) GetPasswordById(ctx context.Context, id int64) (*password.Password, error) {
	cached, err := get[models.V1PasswordDal](ctx, r.redis, r.keyByID(id))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *PasswordCacheRepository) SetPasswordById(ctx context.Context, p *password.Password) error {
	return set(ctx, r.redis, r.keyByID(p.GetID()), models.V1PasswordDalFromDomain(p), passwordTTL)
}

func (r *PasswordCacheRepository) DelPasswordById(ctx context.Context, id int64) error {
	return del(ctx, r.redis, r.keyByID(id))
}

func (r *PasswordCacheRepository) GetPasswordByUserId(ctx context.Context, userId string) (*password.Password, error) {
	cached, err := get[models.V1PasswordDal](ctx, r.redis, r.keyByUserID(userId))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *PasswordCacheRepository) SetPasswordByUserId(ctx context.Context, p *password.Password) error {
	return set(ctx, r.redis, r.keyByUserID(p.GetUserID()), models.V1PasswordDalFromDomain(p), passwordTTL)
}

func (r *PasswordCacheRepository) DelPasswordByUserId(ctx context.Context, userId string) error {
	return del(ctx, r.redis, r.keyByUserID(userId))
}
