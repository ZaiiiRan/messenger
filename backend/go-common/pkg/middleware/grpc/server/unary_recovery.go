package middleware

import (
	"context"
	"runtime/debug"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryMiddleware(log *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorw(
					"grpc.panic",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)

				err = status.Errorf(codes.Internal, "%s", commonerror.ErrInternal.Error())
			}
		}()
		return handler(ctx, req)
	}
}
