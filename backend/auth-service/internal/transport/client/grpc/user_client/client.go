package usergrpcclient

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	grpcclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/client/grpc"
	middleware "github.com/ZaiiiRan/messenger/backend/go-common/pkg/middleware/grpc/client"
	"google.golang.org/grpc"
)

type Client struct {
	client     *grpcclient.GRPCClient
	userClient pb.UserServiceClient
}

func New(
	ctx context.Context,
	cfg settings.GRPCClientSettings,
	unaryExtra []grpc.UnaryClientInterceptor,
	streamExtra []grpc.StreamClientInterceptor,
	extra ...grpc.DialOption,
) (*Client, error) {
	unaryExtra = append(
		[]grpc.UnaryClientInterceptor{
			middleware.PropagateClientMetaUnary(),
		},
		unaryExtra...,
	)

	cl, err := grpcclient.New(ctx, cfg, unaryExtra, streamExtra, extra...)
	if err != nil {
		return nil, err
	}

	return &Client{
		client:     cl,
		userClient: pb.NewUserServiceClient(cl.Conn()),
	}, nil
}

func (c *Client) UserClient() pb.UserServiceClient {
	return c.userClient
}

func (c *Client) Connect(ctx context.Context) error {
	if c.client != nil {
		return c.client.Connect(ctx)
	}
	return nil
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
