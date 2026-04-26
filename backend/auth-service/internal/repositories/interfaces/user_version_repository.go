package interfaces

import (
	"context"

	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type UserVersionRepository interface {
	CreateUserVersion(ctx context.Context, uv *userversion.UserVersion) error
	UpdateUserVersion(ctx context.Context, uv *userversion.UserVersion) error
	DeleteUserVersion(ctx context.Context, uv *userversion.UserVersion) error
	QueryUserVersion(ctx context.Context, query *models.QueryUserVersionDal) (*userversion.UserVersion, error)
}
