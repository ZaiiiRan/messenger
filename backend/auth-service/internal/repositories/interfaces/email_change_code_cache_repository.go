package interfaces

import (
	"context"

	emailchangecode "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code/email_change_code"
)

type EmailChangeCodeCacheRepository interface {
	GetCodeById(ctx context.Context, id int64) (*emailchangecode.EmailChangeCode, error)
	SetCodeById(ctx context.Context, c *emailchangecode.EmailChangeCode) error
	DelCodeById(ctx context.Context, id int64) error
	GetCodeByUserId(ctx context.Context, userId string) (*emailchangecode.EmailChangeCode, error)
	SetCodeByUserId(ctx context.Context, c *emailchangecode.EmailChangeCode) error
	DelCodeByUserId(ctx context.Context, userId string) error
}
