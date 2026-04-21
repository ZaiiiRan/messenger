package middleware

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LogMiddleware(log *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		reqId := ctxmetadata.GetReqIdFromContext(ctx)

		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)

		log.Infow(
			"grpc.request",
			"req_id", reqId,
			"method", info.FullMethod,
			"code", code.String(),
			"duration_ms", time.Since(start).Microseconds(),
		)

		return resp, err
	}
}