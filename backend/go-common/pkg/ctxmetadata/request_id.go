package ctxmetadata

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type CtxKeyReqId struct{}

const RequestIDKey = "x-request-id"

func WithReqId(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, CtxKeyReqId{}, reqId)
}

func GetReqIdFromContext(ctx context.Context) string {
	if v := ctx.Value(CtxKeyReqId{}); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

func GetReqIdFromIncomingContext(ctx context.Context) (string, bool) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(RequestIDKey); len(vals) > 0 && vals[0] != "" {
			return vals[0], true
		}
	}
	return "", false
}

func EnsureReqId(ctx context.Context) (context.Context, string) {
	id := GetReqIdFromContext(ctx)
	if id == "" {
		if v, ok := GetReqIdFromIncomingContext(ctx); ok {
			id = v
		} else {
			id = uuid.NewString()
		}
		ctx = WithReqId(ctx, id)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, RequestIDKey, id)
	return ctx, id
}

func ForwardReqIdToOutgoingContext(ctx context.Context) context.Context {
	id := GetReqIdFromContext(ctx)
	if id == "" {
		if v, ok := GetReqIdFromIncomingContext(ctx); ok {
			id = v
			ctx = WithReqId(ctx, id)
		}
	}
	if id != "" {
		return metadata.AppendToOutgoingContext(ctx, RequestIDKey, id)
	}
	return ctx
}