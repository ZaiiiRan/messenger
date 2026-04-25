package ctxmetadata

import (
	"context"
	"fmt"

	claims "github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	"google.golang.org/grpc/metadata"
)

type CtxKeyClaims struct{}

const AuthorizationKey = "authorization"

func WithUserClaims(ctx context.Context, claims *claims.UserClaims) context.Context {
	return withClaims(ctx, claims)
}

func GetUserClaimsFromContext(ctx context.Context) (*claims.UserClaims, bool) {
	if c, ok := getClaimsFromContext[*claims.UserClaims](ctx); ok {
		return c, true
	}
	return nil, false
}

func GetAuthMetadataFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(AuthorizationKey); len(values) > 0 && values[0] != "" {
			return values[0], nil
		}
	}
	return "", fmt.Errorf("missing metadata")
}

func ForwardAuthToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(AuthorizationKey); len(values) > 0 && values[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, AuthorizationKey, values[0])
		}
	}
	return ctx
}

func withClaims[T any](ctx context.Context, claims T) context.Context {
	return context.WithValue(ctx, CtxKeyClaims{}, claims)
}

func getClaimsFromContext[T any](ctx context.Context) (T, bool) {
	claims, ok := ctx.Value(CtxKeyClaims{}).(T)
	return claims, ok
}
