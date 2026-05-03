package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	commonmiddleware "github.com/ZaiiiRan/messenger/backend/go-common/pkg/middleware/grpc/server"
	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func New(srvSettings settings.GRPCServerSettings, userService userservice.UserService, log *zap.SugaredLogger) (*Server, error) {
	s := grpc.NewServer(
		newChainUnaryInterceptor(log),
		grpc.KeepaliveParams(getGRPCKeepAliveServerParams(&srvSettings)),
		grpc.KeepaliveEnforcementPolicy(getGRPCKeepAliveEnforcement(&srvSettings)),
	)

	pb.RegisterUserServiceServer(s, newUserHandler(userService))

	lis, err := net.Listen("tcp", srvSettings.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	return &Server{
		srv: s,
		lis: lis,
	}, nil
}

func (s *Server) Start() error {
	return s.srv.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

func (s *Server) Addr() string {
	if s.lis != nil {
		return s.lis.Addr().String()
	}
	return ""
}

func newChainUnaryInterceptor(log *zap.SugaredLogger) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		commonmiddleware.RequestIdMiddleware(),
		commonmiddleware.LogMiddleware(log),
		commonmiddleware.RecoveryMiddleware(log),
	)
}

func getGRPCKeepAliveServerParams(c *settings.GRPCServerSettings) keepalive.ServerParameters {
	if c == nil {
		return keepalive.ServerParameters{}
	}
	return keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(c.MaxConnectionIdle) * time.Second,
		MaxConnectionAge:      time.Duration(c.MaxConnectionAge) * time.Second,
		MaxConnectionAgeGrace: time.Duration(c.MaxConnectionAgeGrace) * time.Second,
		Time:                  time.Duration(c.KeepaliveTime) * time.Second,
		Timeout:               time.Duration(c.KeepaliveTimeout) * time.Second,
	}
}

func getGRPCKeepAliveEnforcement(c *settings.GRPCServerSettings) keepalive.EnforcementPolicy {
	if c == nil {
		return keepalive.EnforcementPolicy{}
	}
	return keepalive.EnforcementPolicy{
		MinTime:             0,
		PermitWithoutStream: c.PermitWithoutStream,
	}
}
