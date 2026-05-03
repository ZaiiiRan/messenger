package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
)

type CodeCacheRepository interface {
	GetCodeById(ctx context.Context, id int64) (*code.Code, error)
	SetCodeById(ctx context.Context, code *code.Code) error
	DelCodeById(ctx context.Context, id int64) error
	GetCodeByUserId(ctx context.Context, userId string) (*code.Code, error)
	SetCodeByUserId(ctx context.Context, code *code.Code) error
	DelCodeByUserId(ctx context.Context, userId string) error
}
