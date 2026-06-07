package grpcclient

import (
	"context"

	authpb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/auth/v1"
	filepb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/file/v1"
	mediapb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/media/v1"
	socialpb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/social/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/api-gateway/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config"
	"google.golang.org/grpc"
)

// Clients holds the typed gRPC clients for every upstream microservice the
// gateway talks to, together with the underlying connections so they can be
// closed gracefully on shutdown.
type Clients struct {
	Auth   authpb.AuthServiceClient
	User   userpb.UserServiceClient
	Social socialpb.SocialServiceClient
	File   filepb.FileServiceClient
	Media  mediapb.MediaServiceClient

	conns []*Conn
}

// New dials every upstream service described in the config. The same set of
// extra unary interceptors (e.g. Prometheus client metrics) is applied to all
// connections.
func New(ctx context.Context, cfg *config.ServerConfig, unaryExtra []grpc.UnaryClientInterceptor) (*Clients, error) {
	c := &Clients{}

	authConn, err := Dial(ctx, cfg.AuthServiceGRPCClient, unaryExtra)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.conns = append(c.conns, authConn)
	c.Auth = authpb.NewAuthServiceClient(authConn.Conn())

	userConn, err := Dial(ctx, cfg.UserServiceGRPCClient, unaryExtra)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.conns = append(c.conns, userConn)
	c.User = userpb.NewUserServiceClient(userConn.Conn())

	socialConn, err := Dial(ctx, cfg.SocialServiceGRPCClient, unaryExtra)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.conns = append(c.conns, socialConn)
	c.Social = socialpb.NewSocialServiceClient(socialConn.Conn())

	fileConn, err := Dial(ctx, cfg.FileServiceGRPCClient, unaryExtra)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.conns = append(c.conns, fileConn)
	c.File = filepb.NewFileServiceClient(fileConn.Conn())

	mediaConn, err := Dial(ctx, cfg.MediaServiceGRPCClient, unaryExtra)
	if err != nil {
		c.Close()
		return nil, err
	}
	c.conns = append(c.conns, mediaConn)
	c.Media = mediapb.NewMediaServiceClient(mediaConn.Conn())

	return c, nil
}

func (c *Clients) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}
