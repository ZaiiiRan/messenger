package fileservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
)

type fileDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newFileDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *fileDataProvider {
	return &fileDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (fdp *fileDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(fdp.pg)
}

func (fdp *fileDataProvider) getFileByID(ctx context.Context, id string, uow *uow.UnitOfWork) (*file.File, error) {
	if id == "" {
		return nil, nil
	}

	cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
	f, err := cacheRepo.GetFileByID(ctx, id)
	if err == nil && f != nil {
		return f, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	query := models.NewQueryFilesDal(models.FileFilterDal{Ids: []string{id}}, 1, 1)
	list, err := dbRepo.QueryFiles(ctx, query, false, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetFileByID(ctx, list[0])
	return list[0], nil
}

func (fdp *fileDataProvider) getFileByIDLocked(ctx context.Context, id string, uow *uow.UnitOfWork) (*file.File, error) {
	if id == "" {
		return nil, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	query := models.NewQueryFilesDal(models.FileFilterDal{Ids: []string{id}}, 1, 1)
	list, err := dbRepo.QueryFiles(ctx, query, true, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
}

func (fdp *fileDataProvider) getFiles(
	ctx context.Context,
	query *models.QueryFilesDal,
	uow *uow.UnitOfWork,
) ([]*file.File, error) {
	cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
	list, err := cacheRepo.GetFileList(ctx, query)
	if err == nil && list != nil {
		return list, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	list, err = dbRepo.QueryFiles(ctx, query, false, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetFileList(ctx, query, list)
	return list, nil
}

func (fdp *fileDataProvider) getFilesLocked(
	ctx context.Context,
	query *models.QueryFilesDal,
	skipLocked bool,
	uow *uow.UnitOfWork,
) ([]*file.File, error) {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	list, err := dbRepo.QueryFiles(ctx, query, true, skipLocked)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list, nil
}

func (fdp *fileDataProvider) createFiles(
	ctx context.Context,
	files []*file.File,
	uow *uow.UnitOfWork,
) error {
	if len(files) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	if err := dbRepo.CreateFiles(ctx, files); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
		ownerIDs := make(map[string]struct{})
		for _, f := range files {
			cacheRepo.SetFileByID(ctx, f)
			ownerIDs[f.GetOwnerID()] = struct{}{}
		}
		for ownerID := range ownerIDs {
			cacheRepo.InvalidateFileListsByOwnerID(ctx, ownerID)
		}
	})
	return nil
}

func (fdp *fileDataProvider) updateFiles(
	ctx context.Context,
	files []*file.File,
	uow *uow.UnitOfWork,
) error {
	if len(files) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	if err := dbRepo.UpdateFiles(ctx, files); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
		ownerIDs := make(map[string]struct{})
		for _, f := range files {
			cacheRepo.DelFileByID(ctx, f.GetID())
			cacheRepo.SetFileByID(ctx, f)
			ownerIDs[f.GetOwnerID()] = struct{}{}
		}
		for ownerID := range ownerIDs {
			cacheRepo.InvalidateFileListsByOwnerID(ctx, ownerID)
		}
	})
	return nil
}

func (fdp *fileDataProvider) deleteFiles(
	ctx context.Context,
	files []*file.File,
	uow *uow.UnitOfWork,
) error {
	if len(files) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	if err := dbRepo.DeleteFiles(ctx, files); err != nil {
		return err
	}

	ids := make([]string, len(files))
	ownerIDs := make(map[string]struct{})
	for i, f := range files {
		ids[i] = f.GetID()
		ownerIDs[f.GetOwnerID()] = struct{}{}
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
		for _, id := range ids {
			cacheRepo.DelFileByID(ctx, id)
		}
		for ownerID := range ownerIDs {
			cacheRepo.InvalidateFileListsByOwnerID(ctx, ownerID)
		}
	})
	return nil
}

func (fdp *fileDataProvider) save(ctx context.Context, f *file.File, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)

	if f.GetID() == "" {
		if err := dbRepo.CreateFiles(ctx, []*file.File{f}); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateFiles(ctx, []*file.File{f}); err != nil {
			return err
		}
	}

	ownerID := f.GetOwnerID()
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
		cacheRepo.DelFileByID(ctx, f.GetID())
		cacheRepo.SetFileByID(ctx, f)
		cacheRepo.InvalidateFileListsByOwnerID(ctx, ownerID)
	})
	return nil
}

func (fdp *fileDataProvider) delete(ctx context.Context, f *file.File, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewFileRepository(pgConn)
	if err := dbRepo.DeleteFiles(ctx, []*file.File{f}); err != nil {
		return err
	}

	id := f.GetID()
	ownerID := f.GetOwnerID()
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewFileCacheRepository(fdp.redis)
		cacheRepo.DelFileByID(ctx, id)
		cacheRepo.InvalidateFileListsByOwnerID(ctx, ownerID)
	})
	return nil
}
