package ctxmetadata

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const RealIPKey = "x-real-ip"

func GetRealIPFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(RealIPKey); len(values) > 0 && values[0] != "" {
			return values[0], nil
		}
	}
	return "", fmt.Errorf("missing metadata")
}

func ForwardRealIPToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(RealIPKey); len(values) > 0 && values[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, RealIPKey, values[0])
		}
	}
	return ctx
}
