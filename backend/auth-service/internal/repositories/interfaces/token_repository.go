package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, token *token.Token) error
	UpdateToken(ctx context.Context, token *token.Token) error
	DeleteToken(ctx context.Context, tokenStr string) error
	DeleteTokensByUserId(ctx context.Context, userId string) error
	QueryToken(ctx context.Context, query *models.QueryTokenDal) (*token.Token, error)
}
