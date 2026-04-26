package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/redis"
	goredis "github.com/redis/go-redis/v9"
)

const (
	userKeyPrefix           = "user"
	userKeyPrefixByEmail    = "user:email"
	userKeyPrefixByUsername = "user:username"
	userListKeyPrefix       = "user:list"
	userListIndexKey        = "user:list:index"
	userTTL                 = 10 * time.Minute
	userListTTL             = 5 * time.Minute
)

type cachedUserDal struct {
	User    models.V1UserDal    `json:"user"`
	Profile models.V1ProfileDal `json:"profile"`
	Status  models.V1StatusDal  `json:"status"`
}

func toCachedUserDal(u *user.User) cachedUserDal {
	return cachedUserDal{
		User:    models.V1UserDalFromDomain(u),
		Profile: models.V1ProfileDalFromDomain(u.GetID(), u.GetProfile()),
		Status:  models.V1StatusDalFromDomain(u.GetID(), u.GetStatus()),
	}
}

func (c cachedUserDal) toDomain() *user.User {
	return c.User.ToDomain(c.Profile.ToDomain(), c.Status.ToDomain())
}

type UserCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewUserCacheRepository(redis *rediscl.RedisClient) interfaces.UserCacheRepository {
	return &UserCacheRepository{redis: redis}
}

func (r *UserCacheRepository) keyByID(id string) string {
	return fmt.Sprintf("%s:%s", userKeyPrefix, id)
}

func (r *UserCacheRepository) keyByEmail(email string) string {
	return fmt.Sprintf("%s:%s", userKeyPrefixByEmail, email)
}

func (r *UserCacheRepository) keyByUsername(username string) string {
	return fmt.Sprintf("%s:%s", userKeyPrefixByUsername, username)
}

func (r *UserCacheRepository) listKey(hash string) string {
	return fmt.Sprintf("%s:%s", userListKeyPrefix, hash)
}

func (r *UserCacheRepository) SetUser(ctx context.Context, u *user.User) error {
	if err := set(ctx, r.redis, r.keyByID(u.GetID()), toCachedUserDal(u), userTTL); err != nil {
		return err
	}
	cl := r.redis.GetClient()
	if err := cl.Set(ctx, r.keyByEmail(u.GetEmail()), u.GetID(), userTTL).Err(); err != nil {
		return err
	}
	return cl.Set(ctx, r.keyByUsername(u.GetUsername()), u.GetID(), userTTL).Err()
}

func (r *UserCacheRepository) GetUser(ctx context.Context, id string) (*user.User, error) {
	cached, err := get[cachedUserDal](ctx, r.redis, r.keyByID(id))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.toDomain(), nil
}

func (r *UserCacheRepository) DeleteUser(ctx context.Context, id string) error {
	cached, err := get[cachedUserDal](ctx, r.redis, r.keyByID(id))
	if err != nil {
		return err
	}

	cl := r.redis.GetClient()

	if cached != nil {
		if err := cl.Del(ctx,
			r.keyByEmail(cached.User.Email),
			r.keyByUsername(cached.User.Username),
		).Err(); err != nil {
			return err
		}
	}

	return cl.Del(ctx, r.keyByID(id)).Err()
}

func (r *UserCacheRepository) SetUserByUsername(ctx context.Context, u *user.User) error {
	return r.SetUser(ctx, u)
}

func (r *UserCacheRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.getByPointerKey(ctx, r.keyByUsername(username))
}

func (r *UserCacheRepository) DeleteUserByUsername(ctx context.Context, username string) error {
	return r.deleteByPointerKey(ctx, r.keyByUsername(username))
}

func (r *UserCacheRepository) SetUserByEmail(ctx context.Context, u *user.User) error {
	return r.SetUser(ctx, u)
}

func (r *UserCacheRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.getByPointerKey(ctx, r.keyByEmail(email))
}

func (r *UserCacheRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	return r.deleteByPointerKey(ctx, r.keyByEmail(email))
}

func (r *UserCacheRepository) getByPointerKey(ctx context.Context, pointerKey string) (*user.User, error) {
	id, err := r.redis.GetClient().Get(ctx, pointerKey).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return r.GetUser(ctx, id)
}

func (r *UserCacheRepository) deleteByPointerKey(ctx context.Context, pointerKey string) error {
	id, err := r.redis.GetClient().Get(ctx, pointerKey).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil
		}
		return err
	}
	return r.DeleteUser(ctx, id)
}

func (r *UserCacheRepository) SetUserList(ctx context.Context, query *models.QueryUsersDal, users []*user.User) error {
	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	cached := make([]cachedUserDal, len(users))
	for i, u := range users {
		cached[i] = toCachedUserDal(u)
	}

	key := r.listKey(hash)
	if err := set(ctx, r.redis, key, cached, userListTTL); err != nil {
		return err
	}

	return r.redis.GetClient().SAdd(ctx, userListIndexKey, key).Err()
}

func (r *UserCacheRepository) GetUserList(ctx context.Context, query *models.QueryUsersDal) ([]*user.User, error) {
	hash, err := queryHash(query)
	if err != nil {
		return nil, err
	}

	cachedList, err := get[[]cachedUserDal](ctx, r.redis, r.listKey(hash))
	if err != nil || cachedList == nil {
		return nil, err
	}

	users := make([]*user.User, len(*cachedList))
	for i, c := range *cachedList {
		users[i] = c.toDomain()
	}
	return users, nil
}

func (r *UserCacheRepository) InvalidateUserList(ctx context.Context, query *models.QueryUsersDal) error {
	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	key := r.listKey(hash)
	cl := r.redis.GetClient()

	if err := cl.Del(ctx, key).Err(); err != nil {
		return err
	}
	return cl.SRem(ctx, userListIndexKey, key).Err()
}
