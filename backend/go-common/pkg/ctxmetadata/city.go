package ctxmetadata

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"google.golang.org/grpc/metadata"
)

const CityKey = "x-city"

func GetCityFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(CityKey); len(values) > 0 && values[0] != "" {
			return values[0], nil
		}
	}
	return "", commonerror.ErrMissingMetadata
}

func ForwardCityToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(CityKey); len(values) > 0 && values[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, CityKey, values[0])
		}
	}
	return ctx
}
