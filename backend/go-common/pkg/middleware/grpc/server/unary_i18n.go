package middleware

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/grpc"
)

func I18nMiddleware(
	getLocalizerFunc func(ctx context.Context) *i18n.Localizer,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		localizer := getLocalizerFunc(ctx)
		ctx = ctxmetadata.WithLocalizer(ctx, localizer)

		return handler(ctx, req)
	}
}
