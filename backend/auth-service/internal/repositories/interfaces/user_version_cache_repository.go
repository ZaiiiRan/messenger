package interfaces

import (
	"context"

	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
)

type UserVersionCacheRepository interface {
	GetByUserId(ctx context.Context, userId string) (*userversion.UserVersion, error)
	SetByUserId(ctx context.Context, uv *userversion.UserVersion) error
	DelByUserId(ctx context.Context, userId string) error
}
