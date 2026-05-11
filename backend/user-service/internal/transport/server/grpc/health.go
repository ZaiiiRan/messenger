package grpcserver

import (
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func newHealthServer() *health.Server {
	srv := health.NewServer()
	srv.SetServingStatus("user.v1.UserService", healthpb.HealthCheckResponse_SERVING)
	return srv
}
