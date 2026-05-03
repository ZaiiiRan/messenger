package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type TokenCacheRepository interface {
	GetToken(ctx context.Context, token string) (*token.Token, error)
	SetToken(ctx context.Context, token *token.Token) error
	DelToken(ctx context.Context, token string) error
	GetTokenList(ctx context.Context, query *models.QueryTokensDal) ([]*token.Token, error)
	SetTokenList(ctx context.Context, query *models.QueryTokensDal, tokens []*token.Token) error
	DelTokenListsByUserId(ctx context.Context, userId string) error
}
