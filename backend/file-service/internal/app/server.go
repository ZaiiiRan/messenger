package app

import (
	"context"
	"errors"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config"
	filedeletionfroms3tasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_from_s3_tasks_service"
	fileservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_service"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/i18n"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	prommetrics "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/s3"
	grpcserver "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/server/grpc"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ServerApp struct {
	cfg *config.ServerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient

	userPhotosS3Client *s3.S3Client
	chatFilesS3Client  *s3.S3Client

	fileDeletionFromS3TasksService filedeletionfroms3tasksservice.FileDeletionFromS3TasksService
	fileService                    fileservice.FileService

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
	if err := a.initUserPhotosS3Client(ctx); err != nil {
		return err
	}
	if err := a.initChatFilesS3Client(ctx); err != nil {
		return err
	}

	a.initFileDeletionFromS3TasksService()
	a.initFileService()

	a.initI18n()
	a.initMetricsServer()
	a.startMetricsServer()

	if err := a.initGrpcServer(); err != nil {
		return err
	}
	a.startGrpcServer()

	<-a.isGRPCStarted
	<-a.isMetricsStarted
	a.log.Infow("app.started")
	return nil
}

func (a *ServerApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	a.postgresClient.Close()
	a.redisClient.Close()
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

	a.log.Infow("app.postgres_connected")
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

func (a *ServerApp) initUserPhotosS3Client(ctx context.Context) error {
	s3Client, err := s3.New(ctx, a.cfg.UserPhotosS3)
	if err != nil {
		a.log.Errorw("app.user_photos_s3_client_init_failed", "err", err)
		return err
	}
	a.userPhotosS3Client = s3Client

	a.log.Infow("app.user_photos_s3_connected")
	return nil
}

func (a *ServerApp) initChatFilesS3Client(ctx context.Context) error {
	s3Client, err := s3.New(ctx, a.cfg.ChatFilesS3)
	if err != nil {
		a.log.Errorw("app.chat_files_s3_client_init_failed", "err", err)
		return err
	}
	a.chatFilesS3Client = s3Client

	a.log.Infow("app.chat_files_s3_connected")
	return nil
}

func (a *ServerApp) initI18n() {
	i18n.Init()
}

func (a *ServerApp) initFileDeletionFromS3TasksService() {
	a.fileDeletionFromS3TasksService = filedeletionfroms3tasksservice.New(a.postgresClient, a.userPhotosS3Client, a.chatFilesS3Client, a.log)
}

func (a *ServerApp) initFileService() {
	a.fileService = fileservice.New(a.postgresClient, a.redisClient, a.userPhotosS3Client, a.chatFilesS3Client, a.fileDeletionFromS3TasksService, a.log)
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
	srv, err := grpcserver.New(a.cfg.GRPCServer, a.cfg.JWT, a.fileService, a.log, a.metricsServer.Registry())
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
