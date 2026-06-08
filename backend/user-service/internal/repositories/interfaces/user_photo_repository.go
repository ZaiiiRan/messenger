package interfaces

import (
	"context"

	userphoto "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user_photo"
)

type UserPhotoRepository interface {
	Create(ctx context.Context, photo *userphoto.UserPhoto) (*userphoto.UserPhoto, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int32) ([]*userphoto.UserPhoto, int64, error)
	GetLatestByUserID(ctx context.Context, userID string) (*userphoto.UserPhoto, error)
	GetLatestByUserIDs(ctx context.Context, userIDs []string) (map[string]*userphoto.UserPhoto, error)
	GetFileIDsByUserIDs(ctx context.Context, userIDs []string) (map[string][]string, error)
	GetByID(ctx context.Context, id string) (*userphoto.UserPhoto, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteByUserIDs(ctx context.Context, userIDs []string) error
}
