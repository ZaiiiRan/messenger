package ctxmetadata

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

const CountryNameKey = "x-country-name"

func GetCountryNameFromIncomingContext(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(CountryNameKey); len(values) > 0 && values[0] != "" {
			return values[0], nil
		}
	}
	return "", fmt.Errorf("missing metadata")
}

func ForwardCountryNameToOutgoingContext(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(CountryNameKey); len(values) > 0 && values[0] != "" {
			return metadata.AppendToOutgoingContext(ctx, CountryNameKey, values[0])
		}
	}
	return ctx
}
