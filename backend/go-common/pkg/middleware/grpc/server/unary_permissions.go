package middleware

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserPermissionMiddleware(shouldProtect MethodMatcher) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if shouldProtect == nil || !shouldProtect(info.FullMethod) {
			return handler(ctx, req)
		}

		claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
		if claims == nil || claims.IsDeleted {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
		}

		if claims.IsPermanentlyBanned || claims.IsTemporarilyBanned || !claims.IsConfirmed {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied")
		}

		return handler(ctx, req)
	}
}
