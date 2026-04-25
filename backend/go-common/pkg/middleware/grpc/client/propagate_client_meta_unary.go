package client

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"google.golang.org/grpc"
)

func PropagateClientMetaUnary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = ctxmetadata.ForwardReqIdToOutgoingContext(ctx)
		ctx = ctxmetadata.ForwardAuthToOutgoingContext(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
