package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
)

type TokenCacheRepository interface {
	GetToken(ctx context.Context, token string) (*token.Token, error)
	SetToken(ctx context.Context, token *token.Token) error
	DelToken(ctx context.Context, token string) error
}
