package interfaces

import (
	"context"

	emailchangecode "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code/email_change_code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type EmailChangeCodeRepository interface {
	CreateCode(ctx context.Context, code *emailchangecode.EmailChangeCode) error
	UpdateCode(ctx context.Context, code *emailchangecode.EmailChangeCode) error
	DeleteCode(ctx context.Context, code *emailchangecode.EmailChangeCode) error
	QueryCode(ctx context.Context, query *models.QueryCodeDal) (*emailchangecode.EmailChangeCode, error)
	DeleteExpiredCodes(ctx context.Context, query *models.QueryExpiredCodesDal) ([]*emailchangecode.EmailChangeCode, error)
}
