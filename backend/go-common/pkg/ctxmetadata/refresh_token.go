package ctxmetadata

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const RefreshTokenKey = "x-refresh-token"

func GetRefreshTokenFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(RefreshTokenKey); len(vals) > 0 && vals[0] != "" {
			return vals[0], nil
		}
	}
	return "", fmt.Errorf("missing metadata")
}

func ForwardRefreshTokenToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(RefreshTokenKey); len(vals) > 0 && vals[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, RefreshTokenKey, vals[0])
		}
	}
	return ctx
}
