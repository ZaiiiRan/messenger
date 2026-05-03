package app

import (
	"context"
	"errors"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config"
	implkafkaproducer "github.com/ZaiiiRan/messenger/backend/auth-service/internal/producers/impl/kafka"
	authservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/auth"
	codeservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	userservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_service"
	usergrpcclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/client/grpc/user_client"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/i18n"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	prommetrics "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	grpcserver "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/server/grpc"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerApp struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	postgresClient            *postgres.PostgresClient
	redisClient               *redis.RedisClient
	emailCodeTasksKafkaClient *kafka.KafkaClient

	userGrpcClient *usergrpcclient.Client

	emailCodeTasksProducer *implkafkaproducer.Producer

	userService     userservice.UserService
	codeService     codeservice.CodeService
	passwordService passwordservice.PasswordService
	tokenService    tokenservice.TokenService
	authService     authservice.AuthService

	grpcServer    *grpcserver.Server
	metricsServer *prommetrics.Server

	isGRPCStarted    chan bool
	isMetricsStarted chan bool
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
		cfg:              cfg,
		log:              log,
		isGRPCStarted:    make(chan bool),
		isMetricsStarted: make(chan bool),
	}, nil
}

func (a *ServerApp) Run(ctx context.Context) error {
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}
	if err := a.initRedisClient(ctx); err != nil {
		return err
	}
	if err := a.initEmailCodesTasksKafkaClient(ctx); err != nil {
		return err
	}
	if err := a.initEmailCodeTasksProducer(); err != nil {
		return err
	}
	if err := a.initUserGrpcClient(ctx); err != nil {
		return err
	}

	a.initUserService()
	a.initCodeService()
	a.initPasswordService()
	a.initTokenService()
	a.initAuthService()

	a.initI18n()
	a.initMetricsServer()
	a.startMetricsServer()

	if err := a.initGrpcServer(); err != nil {
		return err
	}
	a.startGrpcServer()

	<-a.isMetricsStarted
	<-a.isGRPCStarted
	a.log.Infow("app.started")
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.userGrpcClient.Close()
	a.emailCodeTasksProducer.Close()
	a.postgresClient.Close()
	a.redisClient.Close()
	a.emailCodeTasksKafkaClient.Close()
	a.grpcServer.Stop(shCtx)
	a.metricsServer.Stop(shCtx)

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

func (a *ServerApp) initEmailCodesTasksKafkaClient(ctx context.Context) error {
	kafkaClient, err := kafka.New(a.cfg.EmailCodesTasksProducer.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.email_codes_tasks_kafka_client_init_failed", "err", err)
		return err
	}
	a.emailCodeTasksKafkaClient = kafkaClient

	a.log.Infow("app.email_codes_tasks_kafka_connected")
	return nil
}

func (a *ServerApp) initEmailCodeTasksProducer() error {
	producer, err := implkafkaproducer.New(a.cfg.EmailCodesTasksProducer, a.emailCodeTasksKafkaClient, a.log)
	if err != nil {
		a.log.Errorw("app.email_code_tasks_producer_init_failed", "err", err)
		return err
	}
	a.emailCodeTasksProducer = producer
	a.log.Infow("app.email_code_tasks_producer_started")
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

func (a *ServerApp) initUserService() {
	a.userService = userservice.New(a.userGrpcClient, a.log)
}

func (a *ServerApp) initCodeService() {
	a.codeService = codeservice.New(a.postgresClient, a.redisClient, a.log)
}

func (a *ServerApp) initPasswordService() {
	a.passwordService = passwordservice.New(a.postgresClient, a.redisClient, a.log)
}

func (a *ServerApp) initTokenService() {
	a.tokenService = tokenservice.New(a.cfg.JWT, a.postgresClient, a.redisClient, a.log)
}

func (a *ServerApp) initAuthService() {
	a.authService = authservice.New(a.codeService, a.passwordService, a.tokenService, a.userService, a.emailCodeTasksProducer, a.postgresClient, a.log)
}

func (a *ServerApp) initMetricsServer() {
	a.metricsServer = prommetrics.New(a.cfg.MetricsServer)
}

func (a *ServerApp) startMetricsServer() {
	go func() {
		a.log.Infow("app.metrics_serve_start", "port", a.cfg.MetricsServer.Port)
		a.isMetricsStarted <- true
		if err := a.metricsServer.Start(); err != nil {
			a.log.Fatalw("app.metrics_serve_error", "err", err)
		}
	}()
}

func (a *ServerApp) initGrpcServer() error {
	srv, err := grpcserver.New(a.cfg.GRPCServer, a.cfg.JWT, a.authService, a.log, a.metricsServer.Registry())
	if err != nil {
		a.log.Errorw("app.grpc_server_init_failed", "err", err)
		return err
	}

	a.grpcServer = srv
	return nil
}

func (a *ServerApp) initI18n() {
	i18n.Init()
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
