package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type PasswordRepository interface {
	CreatePassword(ctx context.Context, password *password.Password) error
	UpdatePassword(ctx context.Context, password *password.Password) error
	DeletePassword(ctx context.Context, password *password.Password) error
	QueryPassword(ctx context.Context, query *models.QueryPasswordDal) (*password.Password, error)
}
