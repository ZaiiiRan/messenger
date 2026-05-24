package redisimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
)

const (
	fileKeyPrefix         = "file"
	fileListPrefix        = "file:list"
	fileListIndexKey      = "file:list:index"
	fileListOwnerPrefix   = "file:list:index:owner"
	fileTTL               = 10 * time.Minute
	fileListTTL           = 5 * time.Minute
)

type FileCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewFileCacheRepository(redis *rediscl.RedisClient) interfaces.FileCacheRepository {
	return &FileCacheRepository{redis: redis}
}

func (r *FileCacheRepository) key(id string) string {
	return fmt.Sprintf("%s:%s", fileKeyPrefix, id)
}

func (r *FileCacheRepository) listKey(hash string) string {
	return fmt.Sprintf("%s:%s", fileListPrefix, hash)
}

func (r *FileCacheRepository) ownerIndexKey(ownerID string) string {
	return fmt.Sprintf("%s:%s", fileListOwnerPrefix, ownerID)
}

func (r *FileCacheRepository) GetFileByID(ctx context.Context, id string) (*file.File, error) {
	cached, err := get[models.V1FileDal](ctx, r.redis, r.key(id))
	if err != nil || cached == nil {
		return nil, err
	}
	return cached.ToDomain(), nil
}

func (r *FileCacheRepository) SetFileByID(ctx context.Context, f *file.File) error {
	return set(ctx, r.redis, r.key(f.GetID()), models.V1FileFromDomain(f), fileTTL)
}

func (r *FileCacheRepository) DelFileByID(ctx context.Context, id string) error {
	return del(ctx, r.redis, r.key(id))
}

func (r *FileCacheRepository) GetFileList(ctx context.Context, query *models.QueryFilesDal) ([]*file.File, error) {
	hash, err := queryHash(query)
	if err != nil {
		return nil, err
	}

	cachedList, err := get[[]models.V1FileDal](ctx, r.redis, r.listKey(hash))
	if err != nil || cachedList == nil {
		return nil, err
	}

	files := make([]*file.File, len(*cachedList))
	for i, c := range *cachedList {
		files[i] = c.ToDomain()
	}
	return files, nil
}

func (r *FileCacheRepository) SetFileList(ctx context.Context, query *models.QueryFilesDal, files []*file.File) error {
	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	cached := make([]models.V1FileDal, len(files))
	for i, f := range files {
		cached[i] = models.V1FileFromDomain(f)
	}

	key := r.listKey(hash)
	if err := set(ctx, r.redis, key, cached, fileListTTL); err != nil {
		return err
	}

	cl := r.redis.GetClient()
	if err := cl.SAdd(ctx, fileListIndexKey, key).Err(); err != nil {
		return err
	}
	for _, ownerID := range query.Filter.OwnerIds {
		if err := cl.SAdd(ctx, r.ownerIndexKey(ownerID), key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *FileCacheRepository) InvalidateFileList(ctx context.Context, query *models.QueryFilesDal) error {
	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	key := r.listKey(hash)
	cl := r.redis.GetClient()
	if err := cl.Del(ctx, key).Err(); err != nil {
		return err
	}
	if err := cl.SRem(ctx, fileListIndexKey, key).Err(); err != nil {
		return err
	}
	for _, ownerID := range query.Filter.OwnerIds {
		cl.SRem(ctx, r.ownerIndexKey(ownerID), key)
	}
	return nil
}

func (r *FileCacheRepository) InvalidateFileListsByOwnerID(ctx context.Context, ownerID string) error {
	indexKey := r.ownerIndexKey(ownerID)
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
