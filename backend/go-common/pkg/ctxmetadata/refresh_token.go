package ctxmetadata

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"google.golang.org/grpc/metadata"
)

const RefreshTokenKey = "x-refresh-token"

func GetRefreshTokenFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(RefreshTokenKey); len(vals) > 0 && vals[0] != "" {
			return vals[0], nil
		}
	}
	return "", commonerror.ErrMissingMetadata
}

func ForwardRefreshTokenToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(RefreshTokenKey); len(vals) > 0 && vals[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, RefreshTokenKey, vals[0])
		}
	}
	return ctx
}
