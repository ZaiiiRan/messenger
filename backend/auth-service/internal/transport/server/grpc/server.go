package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	authservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/auth"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/utils"
	commonmiddleware "github.com/ZaiiiRan/messenger/backend/go-common/pkg/middleware/grpc/server"
	grpc_prom "github.com/grpc-ecosystem/go-grpc-prometheus"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func New(
	srvSettings settings.GRPCServerSettings,
	jwtSettings settings.JWTSettings,
	authService authservice.AuthService,
	log *zap.SugaredLogger,
	reg *prometheus.Registry,
) (*Server, error) {
	grpcMetrics := grpc_prom.NewServerMetrics()
	reg.MustRegister(grpcMetrics)

	s := grpc.NewServer(
		newChainUnaryInterceptor(&jwtSettings, grpcMetrics, log),
		grpc.KeepaliveParams(getGRPCKeepAliveServerParams(&srvSettings)),
		grpc.KeepaliveEnforcementPolicy(getGRPCKeepAliveEnforcement(&srvSettings)),
	)

	pb.RegisterAuthServiceServer(s, newAuthHandler(authService))

	lis, err := net.Listen("tcp", srvSettings.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcMetrics.InitializeMetrics(s)

	hs := newHealthServer()
	healthpb.RegisterHealthServer(s, hs)

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

func newChainUnaryInterceptor(jwtSettings *settings.JWTSettings, grpcMetrics *grpc_prom.ServerMetrics, log *zap.SugaredLogger) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		grpcMetrics.UnaryServerInterceptor(),
		commonmiddleware.RequestIdMiddleware(),
		commonmiddleware.LogMiddleware(log),
		commonmiddleware.RecoveryMiddleware(log),

		commonmiddleware.I18nMiddleware(utils.CreateLocalizer),

		commonmiddleware.UserAuthMiddleware(
			[]byte(jwtSettings.AccessTokenSecret),
			commonmiddleware.MiddlewareOnly(
				"/auth.v1.AuthService/GetNewConfirmationCode",
				"/auth.v1.AuthService/Confirm",
				"/auth.v1.AuthService/ChangePassword",
				"/auth.v1.AuthService/GetActiveSessions",
				"/auth.v1.AuthService/InvalidateSessions",
			),
		),

		commonmiddleware.UserPermissionMiddleware(
			commonmiddleware.MiddlewareOnly(
				"/auth.v1.AuthService/ChangePassword",
				"/auth.v1.AuthService/GetActiveSessions",
				"/auth.v1.AuthService/InvalidateSessions",
			),
		),
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
		MinTime:             10 * time.Second,
		PermitWithoutStream: c.PermitWithoutStream,
	}
}
