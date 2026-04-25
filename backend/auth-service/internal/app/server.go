package app

import (
	"context"
	"errors"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config"
	usergrpcclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/client/grpc/user_client"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	grpcserver "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/server/grpc"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerApp struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient

	userGrpcClient *usergrpcclient.GRPCClient

	grpcServer *grpcserver.Server

	isGRPCStarted chan bool
}

func NewServerApp() (*ServerApp, error) {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	return &ServerApp{
		cfg:           cfg,
		log:           log,
		isGRPCStarted: make(chan bool),
	}, nil
}

func (a *ServerApp) Run(ctx context.Context) error {
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}
	if err := a.initRedisClient(ctx); err != nil {
		return err
	}
	if err := a.initUserGrpcClient(ctx); err != nil {
		return err
	}

	if err := a.initGrpcServer(); err != nil {
		return err
	}
	a.startGrpcServer()

	<-a.isGRPCStarted
	a.log.Infow("app.started")
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.userGrpcClient.Close()
	a.postgresClient.Close()
	a.redisClient.Close()
	a.grpcServer.Stop(shCtx)

	a.log.Infow("app.stopped")
}

func (a *ServerApp) initPostgresClient(ctx context.Context) error {
	if a.cfg.Migrate.NeedToMigrate {
		err := postgres.Migrate(ctx, a.cfg.DB)
		if err != nil {
			a.log.Errorw("app.postgres_migrate_failed", "err", err)
			return err
		}
	} else {
		a.log.Infow("app.postgres_migration_skipped")
	}

	pgClient, err := postgres.New(ctx, a.cfg.DB)
	if err != nil {
		a.log.Errorw("app.postgres_client_init_failed", "err", err)
		return err
	}
	a.postgresClient = pgClient

	a.log.Infow("app.postgres_connectd")
	return nil
}

func (a *ServerApp) initRedisClient(ctx context.Context) error {
	redisClient, err := redis.New(ctx, a.cfg.Redis)
	if err != nil {
		a.log.Errorw("app.redis_connect_failed", "err", err)
		return err
	}
	a.redisClient = redisClient

	a.log.Infow("app.redis_connected")
	return nil
}

func (a *ServerApp) initUserGrpcClient(ctx context.Context) error {
	userClient, err := usergrpcclient.New(ctx, a.cfg.UserServiceGRPCClient, nil, nil)
	if err != nil {
		a.log.Errorw("app.user_grpc_client_init_failed", "err", err)
		return err
	}
	a.userGrpcClient = userClient

	if a.cfg.UserServiceGRPCClient.AutoConnect {
		a.log.Infow("app.user_grpc_client_connected")
	} else {
		a.log.Infow("app.user_grpc_client_initialized")
	}
	return nil
}

func (a *ServerApp) initGrpcServer() error {
	srv, err := grpcserver.New(a.cfg.GRPCServer, a.log)
	if err != nil {
		a.log.Errorw("app.grpc_server_init_failed", "err", err)
		return err
	}

	a.grpcServer = srv
	return nil
}

func (a *ServerApp) startGrpcServer() {
	go func() {
		a.log.Infow("app.grpc_serve_start", "port", a.cfg.GRPCServer.Port)
		a.isGRPCStarted <- true
		if err := a.grpcServer.Start(); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			a.log.Fatalw("app.grpc_serve_error", "err", err)
		}
	}()
}
