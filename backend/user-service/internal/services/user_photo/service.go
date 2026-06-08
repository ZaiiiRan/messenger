package userphotosvc

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	filepb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/file/v1"
	userphoto "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user_photo"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/postgres"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	filegrpcclient "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/client/grpc/file_client"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"go.uber.org/zap"
)

const expectedBucket = "user_photos"

type UserPhotoService interface {
	AddPhoto(ctx context.Context, userID, fileID string) (*userphoto.UserPhoto, error)
	GetPhotos(ctx context.Context, userID string, limit, offset int32) ([]*userphoto.UserPhoto, int64, error)
	GetLatestPhoto(ctx context.Context, userID string) (*userphoto.UserPhoto, error)
	GetLatestPhotosByUserIDs(ctx context.Context, userIDs []string) (map[string]*userphoto.UserPhoto, error)
	DeletePhoto(ctx context.Context, photoID, userID string) error
}

type service struct {
	pg         *postgres.PostgresClient
	fileClient *filegrpcclient.Client
	log        *zap.SugaredLogger
}

func New(pg *postgres.PostgresClient, fileClient *filegrpcclient.Client, log *zap.SugaredLogger) UserPhotoService {
	return &service{pg: pg, fileClient: fileClient, log: log}
}

func (s *service) AddPhoto(ctx context.Context, userID, fileID string) (*userphoto.UserPhoto, error) {
	l := s.log.With("op", "user_photo.add", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "user_id", userID)

	if err := s.validateFile(ctx, userID, fileID); err != nil {
		return nil, err
	}

	u := uow.New(s.pg)
	defer u.Close()

	conn, err := u.GetConn(ctx)
	if err != nil {
		l.Errorw("user_photo.add_failed.get_conn", "err", err)
		return nil, ErrInternalError
	}

	if _, err := u.BeginTransaction(ctx); err != nil {
		l.Errorw("user_photo.add_failed.begin_tx", "err", err)
		return nil, ErrInternalError
	}

	photo, err := postgresimpl.NewUserPhotoRepository(conn).Create(ctx, userphoto.New(userID, fileID))
	if err != nil {
		l.Errorw("user_photo.add_failed.create", "err", err)
		return nil, ErrInternalError
	}

	if err := u.Commit(ctx); err != nil {
		l.Errorw("user_photo.add_failed.commit", "err", err)
		return nil, ErrInternalError
	}

	l.Infow("user_photo.add_success", "photo_id", photo.GetID())
	return photo, nil
}

func (s *service) GetPhotos(ctx context.Context, userID string, limit, offset int32) ([]*userphoto.UserPhoto, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	u := uow.New(s.pg)
	defer u.Close()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return nil, 0, ErrInternalError
	}

	photos, total, err := postgresimpl.NewUserPhotoRepository(conn).GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, ErrInternalError
	}
	return photos, total, nil
}

func (s *service) GetLatestPhoto(ctx context.Context, userID string) (*userphoto.UserPhoto, error) {
	u := uow.New(s.pg)
	defer u.Close()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return nil, nil
	}

	return postgresimpl.NewUserPhotoRepository(conn).GetLatestByUserID(ctx, userID)
}

func (s *service) GetLatestPhotosByUserIDs(ctx context.Context, userIDs []string) (map[string]*userphoto.UserPhoto, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	u := uow.New(s.pg)
	defer u.Close()

	conn, err := u.GetConn(ctx)
	if err != nil {
		return nil, nil
	}

	return postgresimpl.NewUserPhotoRepository(conn).GetLatestByUserIDs(ctx, userIDs)
}

func (s *service) DeletePhoto(ctx context.Context, photoID, userID string) error {
	l := s.log.With("op", "user_photo.delete", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "photo_id", photoID)

	u := uow.New(s.pg)
	defer u.Close()

	conn, err := u.GetConn(ctx)
	if err != nil {
		l.Errorw("user_photo.delete_failed.get_conn", "err", err)
		return ErrInternalError
	}

	repo := postgresimpl.NewUserPhotoRepository(conn)

	photo, err := repo.GetByID(ctx, photoID)
	if err != nil {
		return ErrInternalError
	}
	if photo == nil {
		return ErrPhotoNotFound
	}
	if photo.GetUserID() != userID {
		return ErrPhotoNotOwned
	}

	if _, err := u.BeginTransaction(ctx); err != nil {
		l.Errorw("user_photo.delete_failed.begin_tx", "err", err)
		return ErrInternalError
	}

	if err := repo.DeleteByID(ctx, photoID); err != nil {
		l.Errorw("user_photo.delete_failed.delete", "err", err)
		return ErrInternalError
	}

	if err := u.Commit(ctx); err != nil {
		l.Errorw("user_photo.delete_failed.commit", "err", err)
		return ErrInternalError
	}

	l.Infow("user_photo.delete_success")
	return nil
}

func (s *service) validateFile(ctx context.Context, userID, fileID string) error {
	resp, err := s.fileClient.FileClient().GetFiles(ctx, &filepb.GetFilesRequest{
		FileIds:                []string{fileID},
		IncludeProprietaryInfo: true,
	})
	if err != nil || len(resp.GetFiles()) == 0 {
		return ErrFileNotFound
	}

	f := resp.Files[0]

	if f.OwnerId != userID {
		return ErrFileNotOwned
	}
	if f.Status == "waiting_for_upload" {
		return ErrFileNotReady
	}
	for _, v := range f.Variants {
		if v.Bucket != nil && *v.Bucket != expectedBucket {
			return ErrInvalidFileBucket
		}
	}
	return nil
}
