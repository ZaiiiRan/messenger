package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GRPCClient struct {
	conn             *grpc.ClientConn
	waitReadyTimeout time.Duration
	opts             []grpc.DialOption
}

func New(
	ctx context.Context,
	cfg settings.GRPCClientSettings,
	unaryExtra []grpc.UnaryClientInterceptor,
	streamExtra []grpc.StreamClientInterceptor,
	extra ...grpc.DialOption,
) (*GRPCClient, error) {
	dialOpts := buildClientDialOptions(cfg, unaryExtra, streamExtra, extra...)
	conn, err := grpc.NewClient(cfg.Address, dialOpts...)
	if err != nil {
		return nil, err
	}

	cl := &GRPCClient{
		conn:             conn,
		waitReadyTimeout: time.Duration(cfg.WaitGRPCReadyTimeout) * time.Second,
		opts:             dialOpts,
	}

	if cfg.AutoConnect {
		if err := cl.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return cl, nil
}

func (c *GRPCClient) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *GRPCClient) Connect(ctx context.Context) error {
	if c.conn.GetState() == connectivity.Ready || c.conn.GetState() == connectivity.Idle || c.conn.GetState() == connectivity.Connecting {
		return nil
	}
	c.conn.Connect()
	waitCtx, cancel := context.WithTimeout(ctx, c.waitReadyTimeout)
	defer cancel()
	return c.waitUntilReady(waitCtx)
}

func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *GRPCClient) waitUntilReady(ctx context.Context) error {
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

func buildClientDialOptions(
	cfg settings.GRPCClientSettings,
	unaryExtra []grpc.UnaryClientInterceptor,
	streamExtra []grpc.StreamClientInterceptor,
	extra ...grpc.DialOption,
) []grpc.DialOption {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithCodes(codes.Aborted, codes.DeadlineExceeded, codes.Internal),
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

	dialOpts = append(dialOpts,
		grpc.WithChainUnaryInterceptor(
			append(unaryExtra, grpc_retry.UnaryClientInterceptor(retryOpts...))...,
		),
		grpc.WithChainStreamInterceptor(
			append(streamExtra, grpc_retry.StreamClientInterceptor(retryOpts...))...,
		),
	)

	for _, opt := range extra {
        if opt != nil {
            dialOpts = append(dialOpts, opt)
        }
    }
	return dialOpts
}