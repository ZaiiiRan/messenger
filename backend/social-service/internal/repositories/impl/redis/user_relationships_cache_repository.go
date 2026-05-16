package redisimpl

import (
	"context"
	"fmt"
	"time"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/redis"
)

const (
	userRelationshipKeyPrefix          = "user_relationships:pair"
	userRelationshiplistKeyPrefix      = "user_relationships:list"
	userRelationshiplistIndexKeyPrefix = "user_relationships:list:index"

	userRelationshipTTL     = 10 * time.Minute
	userRelationshiplistTTL = 5 * time.Minute
)

type UserRelationshipsCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewUserRelationshipsCacheRepository(redis *rediscl.RedisClient) interfaces.UserRelationshipsCacheRepository {
	return &UserRelationshipsCacheRepository{redis: redis}
}

func (r *UserRelationshipsCacheRepository) pairKey(uid1, uid2 string) string {
	if uid1 > uid2 {
		uid1, uid2 = uid2, uid1
	}
	return fmt.Sprintf("%s:%s:%s", userRelationshipKeyPrefix, uid1, uid2)
}

func (r *UserRelationshipsCacheRepository) listKey(hash string) string {
	return fmt.Sprintf("%s:%s", userRelationshiplistKeyPrefix, hash)
}

func (r *UserRelationshipsCacheRepository) listIndexKey(userId string) string {
	return fmt.Sprintf("%s:%s", userRelationshiplistIndexKeyPrefix, userId)
}

func (r *UserRelationshipsCacheRepository) SetUserRelationships(ctx context.Context, ur *userrelationship.UserRelationship) error {
	cached := models.V1UserRelationshipDalFromDomain(ur)
	return set(ctx, r.redis, r.pairKey(ur.GetUserID1(), ur.GetUserID2()), cached, userRelationshipTTL)
}

func (r *UserRelationshipsCacheRepository) GetUserRelationships(
	ctx context.Context,
	firstUserId, secondUserId string,
) (*userrelationship.UserRelationship, error) {
	cached, err := get[models.V1UserRelationshipDal](ctx, r.redis, r.pairKey(firstUserId, secondUserId))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *UserRelationshipsCacheRepository) DelUserRelationships(
	ctx context.Context,
	firstUserId, secondUserId string,
) error {
	return del(ctx, r.redis, r.pairKey(firstUserId, secondUserId))
}

func (r *UserRelationshipsCacheRepository) SetUserRelationshipsList(
	ctx context.Context,
	query *models.QueryUserRelationshipsDal,
	urs []*userrelationship.UserRelationship,
) error {
	if query.FirstUserId == nil {
		return nil
	}

	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	cached := make([]models.V1UserRelationshipDal, len(urs))
	for i, ur := range urs {
		cached[i] = models.V1UserRelationshipDalFromDomain(ur)
	}

	key := r.listKey(hash)
	if err := set(ctx, r.redis, key, cached, userRelationshiplistTTL); err != nil {
		return err
	}

	return r.redis.GetClient().SAdd(ctx, r.listIndexKey(*query.FirstUserId), key).Err()
}

func (r *UserRelationshipsCacheRepository) GetUserRelationshipsList(
	ctx context.Context,
	query *models.QueryUserRelationshipsDal,
) ([]*userrelationship.UserRelationship, error) {
	hash, err := queryHash(query)
	if err != nil {
		return nil, err
	}

	cachedList, err := get[[]models.V1UserRelationshipDal](ctx, r.redis, r.listKey(hash))
	if err != nil || cachedList == nil {
		return nil, err
	}

	urs := make([]*userrelationship.UserRelationship, len(*cachedList))
	for i, c := range *cachedList {
		urs[i] = c.ToDomain()
	}
	return urs, nil
}

func (r *UserRelationshipsCacheRepository) InvalidateUserRelationshipsLists(
	ctx context.Context,
	userId string,
) error {
	indexKey := r.listIndexKey(userId)
	cl := r.redis.GetClient()

	keys, err := cl.SMembers(ctx, indexKey).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	return cl.Del(ctx, append(keys, indexKey)...).Err()
}
