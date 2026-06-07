package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config/settings"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Conn wraps a single gRPC client connection to an upstream microservice.
type Conn struct {
	conn             *grpc.ClientConn
	waitReadyTimeout time.Duration
}

func Dial(
	ctx context.Context,
	cfg settings.GRPCClientSettings,
	unaryExtra []grpc.UnaryClientInterceptor,
	extra ...grpc.DialOption,
) (*Conn, error) {
	dialOpts := buildDialOptions(cfg, unaryExtra, extra...)

	conn, err := grpc.NewClient(cfg.Address, dialOpts...)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		conn:             conn,
		waitReadyTimeout: time.Duration(cfg.WaitGRPCReadyTimeout) * time.Second,
	}

	if cfg.AutoConnect {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Conn) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Conn) Connect(ctx context.Context) error {
	state := c.conn.GetState()
	if state == connectivity.Ready || state == connectivity.Idle || state == connectivity.Connecting {
		return nil
	}
	c.conn.Connect()

	waitCtx, cancel := context.WithTimeout(ctx, c.waitReadyTimeout)
	defer cancel()
	return c.waitUntilReady(waitCtx)
}

func (c *Conn) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Conn) waitUntilReady(ctx context.Context) error {
	for {
		switch c.conn.GetState() {
		case connectivity.Ready:
			return nil
		case connectivity.Shutdown:
			return fmt.Errorf("gRPC connection is shutdown")
		default:
			if !c.conn.WaitForStateChange(ctx, c.conn.GetState()) {
				return ctx.Err()
			}
		}
	}
}

func buildDialOptions(
	cfg settings.GRPCClientSettings,
	unaryExtra []grpc.UnaryClientInterceptor,
	extra ...grpc.DialOption,
) []grpc.DialOption {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithCodes(codes.Aborted, codes.DeadlineExceeded, codes.Unavailable),
		grpc_retry.WithMax(cfg.RetriesCount),
		grpc_retry.WithPerRetryTimeout(time.Duration(cfg.PerCallTimeout) * time.Second),
	}

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Duration(cfg.BackoffBaseDelay) * time.Millisecond,
				Multiplier: cfg.BackoffMultiplier,
				MaxDelay:   time.Duration(cfg.BackoffMaxDelay) * time.Millisecond,
			},
			MinConnectTimeout: time.Duration(cfg.MinConnectTimeout) * time.Second,
		}),
	}

	if cfg.LBPolicy != "" {
		dialOpts = append(dialOpts,
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, cfg.LBPolicy)),
		)
	}

	if cfg.KeepaliveTime > 0 && cfg.KeepaliveTimeout > 0 {
		dialOpts = append(dialOpts,
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                time.Duration(cfg.KeepaliveTime) * time.Second,
				Timeout:             time.Duration(cfg.KeepaliveTimeout) * time.Second,
				PermitWithoutStream: cfg.KeepalivePermitWithoutStream,
			}),
		)
	}

	unary := append([]grpc.UnaryClientInterceptor{}, unaryExtra...)
	unary = append(unary, grpc_retry.UnaryClientInterceptor(retryOpts...))

	dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(unary...))

	for _, opt := range extra {
		if opt != nil {
			dialOpts = append(dialOpts, opt)
		}
	}
	return dialOpts
}
