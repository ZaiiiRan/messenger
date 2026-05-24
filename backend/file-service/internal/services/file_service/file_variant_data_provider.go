package fileservice

import (
	"context"

	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
)

type fileVariantDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newFileVariantDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *fileVariantDataProvider {
	return &fileVariantDataProvider{pg: pg, redis: redis}
}

func (fvdp *fileVariantDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(fvdp.pg)
}

func (fvdp *fileVariantDataProvider) getFileVariantsByFileID(
	ctx context.Context,
	fileID string,
	uow *uow.UnitOfWork,
) ([]*filevariant.FileVariant, error) {
	if fileID == "" {
		return nil, nil
	}

	cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
	list, err := cacheRepo.GetFileVariantsByFileId(ctx, fileID)
	if err == nil && list != nil {
		return list, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	query := models.NewQueryFileVariantsDal(models.FileVariantFilterDal{FileIds: []string{fileID}}, nil, nil)
	list, err = dbRepo.QueryFileVariants(ctx, query, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetFileVariantsByFileId(ctx, fileID, list)
	return list, nil
}

func (fvdp *fileVariantDataProvider) getFileVariantsByFileIDLocked(
	ctx context.Context,
	fileID string,
	uow *uow.UnitOfWork,
) ([]*filevariant.FileVariant, error) {
	if fileID == "" {
		return nil, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	query := models.NewQueryFileVariantsDal(models.FileVariantFilterDal{FileIds: []string{fileID}}, nil, nil)
	list, err := dbRepo.QueryFileVariants(ctx, query, true)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list, nil
}

func (fvdp *fileVariantDataProvider) getFileVariants(
	ctx context.Context,
	query *models.QueryFileVariantsDal,
	uow *uow.UnitOfWork,
) ([]*filevariant.FileVariant, error) {
	cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
	list, err := cacheRepo.GetFileVariantsList(ctx, query)
	if err == nil && list != nil {
		return list, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	list, err = dbRepo.QueryFileVariants(ctx, query, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetFileVariantsList(ctx, query, list)
	return list, nil
}

func (fvdp *fileVariantDataProvider) getFileVariantsLocked(
	ctx context.Context,
	query *models.QueryFileVariantsDal,
	uow *uow.UnitOfWork,
) ([]*filevariant.FileVariant, error) {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	list, err := dbRepo.QueryFileVariants(ctx, query, true)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list, nil
}

func (fvdp *fileVariantDataProvider) createFileVariants(
	ctx context.Context,
	variants []*filevariant.FileVariant,
	uow *uow.UnitOfWork,
) error {
	if len(variants) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	if err := dbRepo.CreateFileVariants(ctx, variants); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
		fileIDs := make(map[string]struct{})
		for _, fv := range variants {
			fileIDs[fv.GetFileID()] = struct{}{}
		}
		for fileID := range fileIDs {
			cacheRepo.DelFileVariantsByFileId(ctx, fileID)
			cacheRepo.InvalidateFileVariantListsByFileID(ctx, fileID)
		}
	})
	return nil
}

func (fvdp *fileVariantDataProvider) updateFileVariants(
	ctx context.Context,
	variants []*filevariant.FileVariant,
	uow *uow.UnitOfWork,
) error {
	if len(variants) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	if err := dbRepo.UpdateFileVariants(ctx, variants); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
		fileIDs := make(map[string]struct{})
		for _, fv := range variants {
			fileIDs[fv.GetFileID()] = struct{}{}
		}
		for fileID := range fileIDs {
			cacheRepo.DelFileVariantsByFileId(ctx, fileID)
			cacheRepo.InvalidateFileVariantListsByFileID(ctx, fileID)
		}
	})
	return nil
}

func (fvdp *fileVariantDataProvider) deleteFileVariants(
	ctx context.Context,
	variants []*filevariant.FileVariant,
	uow *uow.UnitOfWork,
) error {
	if len(variants) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	if err := dbRepo.DeleteFileVariants(ctx, variants); err != nil {
		return err
	}

	fileIDs := make(map[string]struct{})
	for _, fv := range variants {
		fileIDs[fv.GetFileID()] = struct{}{}
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
		for fileID := range fileIDs {
			cacheRepo.DelFileVariantsByFileId(ctx, fileID)
			cacheRepo.InvalidateFileVariantListsByFileID(ctx, fileID)
		}
	})
	return nil
}

func (fvdp *fileVariantDataProvider) save(
	ctx context.Context,
	fv *filevariant.FileVariant,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)

	if fv.GetID() == "" {
		if err := dbRepo.CreateFileVariants(ctx, []*filevariant.FileVariant{fv}); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateFileVariants(ctx, []*filevariant.FileVariant{fv}); err != nil {
			return err
		}
	}

	fileID := fv.GetFileID()
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
		cacheRepo.DelFileVariantsByFileId(ctx, fileID)
		cacheRepo.InvalidateFileVariantListsByFileID(ctx, fileID)
	})
	return nil
}

func (fvdp *fileVariantDataProvider) delete(
	ctx context.Context,
	fv *filevariant.FileVariant,
	uow *uow.UnitOfWork,
) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileVariantRepository(pgConn)
	if err := dbRepo.DeleteFileVariants(ctx, []*filevariant.FileVariant{fv}); err != nil {
		return err
	}

	fileID := fv.GetFileID()
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileVariantCacheRepository(fvdp.redis)
		cacheRepo.DelFileVariantsByFileId(ctx, fileID)
		cacheRepo.InvalidateFileVariantListsByFileID(ctx, fileID)
	})
	return nil
}
