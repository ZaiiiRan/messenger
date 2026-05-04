package middleware

import (
	"context"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserAuthMiddleware(secretKey []byte, shouldProtect MethodMatcher) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if shouldProtect == nil || !shouldProtect(info.FullMethod) {
			return handler(ctx, req)
		}

		tokenStr, err := extractBearerToken(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%s", commonerror.ErrUnauthorized.Error())
		}

		claims, err := jwt.ParseUserToken(tokenStr, secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "%s", commonerror.ErrUnauthorized.Error())
		}

		ctx = ctxmetadata.WithUserClaims(ctx, claims)

		return handler(ctx, req)
	}
}

func extractBearerToken(ctx context.Context) (string, error) {
	authHeader, err := ctxmetadata.GetAuthMetadataFromIncomingContext(ctx)
	if err != nil || len(authHeader) == 0 {
		return "", err
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", commonerror.ErrUnauthorized
	}

	return parts[1], nil
}
