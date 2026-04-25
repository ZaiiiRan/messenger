package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
)

type PasswordCacheRepository interface {
	GetPasswordById(ctx context.Context, id int64) (*password.Password, error)
	SetPasswordById(ctx context.Context, password *password.Password) error
	DelPasswordById(ctx context.Context, id int64) error
	GetPasswordByUserId(ctx context.Context, userId string) (*password.Password, error)
	SetPasswordByUserId(ctx context.Context, password *password.Password) error
	DelPasswordByUserId(ctx context.Context, userId string) error
}
