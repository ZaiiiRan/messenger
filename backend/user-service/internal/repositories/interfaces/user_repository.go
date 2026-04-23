package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	Query(ctx context.Context, query *models.QueryUsersDal) ([]*user.User, error)
}
