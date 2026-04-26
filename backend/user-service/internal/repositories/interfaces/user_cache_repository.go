package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
)

type UserCacheRepository interface {
	SetUser(ctx context.Context, user *user.User) error
	GetUser(ctx context.Context, id string) (*user.User, error)
	DeleteUser(ctx context.Context, id string) error
	SetUserByUsername(ctx context.Context, user *user.User) error
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	DeleteUserByUsername(ctx context.Context, username string) error
	SetUserByEmail(ctx context.Context, user *user.User) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	DeleteUserByEmail(ctx context.Context, email string) error
	SetUserList(ctx context.Context, query *models.QueryUsersDal, users []*user.User) error
	GetUserList(ctx context.Context, query *models.QueryUsersDal) ([]*user.User, error)
	InvalidateUserList(ctx context.Context, query *models.QueryUsersDal) error
}
