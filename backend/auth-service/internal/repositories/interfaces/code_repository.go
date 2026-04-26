package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type CodeRepository interface {
	CreateCode(ctx context.Context, code *code.Code) error
	UpdateCode(ctx context.Context, code *code.Code) error
	DeleteCode(ctx context.Context, code *code.Code) error
	QueryCode(ctx context.Context, query *models.QueryCodeDal) (*code.Code, error)
}
