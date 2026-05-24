package redisimpl

import (
	"context"
	"fmt"
	"time"

	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	rediscl "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
)

const (
	fileVariantByFilePrefix      = "file_variant:file"
	fileVariantListPrefix        = "file_variant:list"
	fileVariantListIndexKey      = "file_variant:list:index"
	fileVariantListFilePrefix    = "file_variant:list:index:file"
	fileVariantTTL               = 10 * time.Minute
	fileVariantListTTL           = 5 * time.Minute
)

type FileVariantCacheRepository struct {
	redis *rediscl.RedisClient
}

func NewFileVariantCacheRepository(redis *rediscl.RedisClient) interfaces.FileVariantCacheRepository {
	return &FileVariantCacheRepository{redis: redis}
}

func (r *FileVariantCacheRepository) fileKey(fileId string) string {
	return fmt.Sprintf("%s:%s", fileVariantByFilePrefix, fileId)
}

func (r *FileVariantCacheRepository) listKey(hash string) string {
	return fmt.Sprintf("%s:%s", fileVariantListPrefix, hash)
}

func (r *FileVariantCacheRepository) fileIndexKey(fileID string) string {
	return fmt.Sprintf("%s:%s", fileVariantListFilePrefix, fileID)
}

func variantsToDal(variants []*filevariant.FileVariant) ([]models.V1FileVariantDal, error) {
	dals := make([]models.V1FileVariantDal, len(variants))
	for i, fv := range variants {
		dal, err := models.V1FileVariantFromDomain(fv)
		if err != nil {
			return nil, err
		}
		dals[i] = dal
	}
	return dals, nil
}

func dalToVariants(dals []models.V1FileVariantDal) ([]*filevariant.FileVariant, error) {
	variants := make([]*filevariant.FileVariant, len(dals))
	for i, d := range dals {
		fv, err := d.ToDomain()
		if err != nil {
			return nil, err
		}
		variants[i] = fv
	}
	return variants, nil
}

func (r *FileVariantCacheRepository) GetFileVariantsByFileId(ctx context.Context, fileId string) ([]*filevariant.FileVariant, error) {
	cachedList, err := get[[]models.V1FileVariantDal](ctx, r.redis, r.fileKey(fileId))
	if err != nil || cachedList == nil {
		return nil, err
	}
	return dalToVariants(*cachedList)
}

func (r *FileVariantCacheRepository) SetFileVariantsByFileId(ctx context.Context, fileId string, variants []*filevariant.FileVariant) error {
	dals, err := variantsToDal(variants)
	if err != nil {
		return err
	}
	return set(ctx, r.redis, r.fileKey(fileId), dals, fileVariantTTL)
}

func (r *FileVariantCacheRepository) DelFileVariantsByFileId(ctx context.Context, fileId string) error {
	return del(ctx, r.redis, r.fileKey(fileId))
}

func (r *FileVariantCacheRepository) GetFileVariantsList(ctx context.Context, query *models.QueryFileVariantsDal) ([]*filevariant.FileVariant, error) {
	hash, err := queryHash(query)
	if err != nil {
		return nil, err
	}

	cachedList, err := get[[]models.V1FileVariantDal](ctx, r.redis, r.listKey(hash))
	if err != nil || cachedList == nil {
		return nil, err
	}
	return dalToVariants(*cachedList)
}

func (r *FileVariantCacheRepository) SetFileVariantsList(ctx context.Context, query *models.QueryFileVariantsDal, fileVariants []*filevariant.FileVariant) error {
	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	dals, err := variantsToDal(fileVariants)
	if err != nil {
		return err
	}

	key := r.listKey(hash)
	if err := set(ctx, r.redis, key, dals, fileVariantListTTL); err != nil {
		return err
	}

	cl := r.redis.GetClient()
	if err := cl.SAdd(ctx, fileVariantListIndexKey, key).Err(); err != nil {
		return err
	}
	for _, fileID := range query.Filter.FileIds {
		if err := cl.SAdd(ctx, r.fileIndexKey(fileID), key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *FileVariantCacheRepository) InvalidateFileVariantsList(ctx context.Context, query *models.QueryFileVariantsDal) error {
	hash, err := queryHash(query)
	if err != nil {
		return err
	}

	key := r.listKey(hash)
	cl := r.redis.GetClient()
	if err := cl.Del(ctx, key).Err(); err != nil {
		return err
	}
	if err := cl.SRem(ctx, fileVariantListIndexKey, key).Err(); err != nil {
		return err
	}
	for _, fileID := range query.Filter.FileIds {
		cl.SRem(ctx, r.fileIndexKey(fileID), key)
	}
	return nil
}

func (r *FileVariantCacheRepository) InvalidateFileVariantListsByFileID(ctx context.Context, fileID string) error {
	indexKey := r.fileIndexKey(fileID)
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
