package ctxmetadata

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const UserAgentKey = "x-user-agent"

func GetUAFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(UserAgentKey); len(values) > 0 && values[0] != "" {
			return values[0], nil
		}
	}
	return "", fmt.Errorf("missing metadata")
}

func ForwardUAToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(UserAgentKey); len(values) > 0 && values[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, UserAgentKey, values[0])
		}
	}
	return ctx
}
