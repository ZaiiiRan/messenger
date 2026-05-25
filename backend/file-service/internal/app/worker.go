package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config"
	filedeletionfroms3tasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_from_s3_tasks_service"
	filedeletiontasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_tasks_service"
	fileservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_service"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/postgres"
	prommetrics "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/s3"
	filedeletionworker "github.com/ZaiiiRan/messenger/backend/file-service/internal/workers/file_deletion"
	filedeletionfroms3worker "github.com/ZaiiiRan/messenger/backend/file-service/internal/workers/file_deletion_from_s3"
	filedeletiontasksconsumerworker "github.com/ZaiiiRan/messenger/backend/file-service/internal/workers/file_deletion_tasks_consumer"
	notinusefiledeletionworker "github.com/ZaiiiRan/messenger/backend/file-service/internal/workers/not_in_use_file_deletion"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient

	userPhotosS3Client *s3.S3Client
	chatFilesS3Client  *s3.S3Client

	fileDeletionTasksKafkaClient *kafka.KafkaClient

	fileDeletionFromS3TasksService filedeletionfroms3tasksservice.FileDeletionFromS3TasksService
	fileService                    fileservice.FileService
	fileDeletionTasksService       filedeletiontasksservice.FileDeletionTasksService

	metricsServer    *prommetrics.Server
	workerMetrics    *prommetrics.WorkerMetrics
	isMetricsStarted chan bool

	workersCtx    context.Context
	workersCancel context.CancelFunc
	workersWg     sync.WaitGroup
}

func NewWorkerApp() (*WorkerApp, error) {
	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	return &WorkerApp{
		cfg:              cfg,
		log:              log,
		isMetricsStarted: make(chan bool),
	}, nil
}

func (a *WorkerApp) Run(ctx context.Context) error {
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
	if err := a.initFileDeletionTasksKafkaClient(ctx); err != nil {
		return err
	}

	a.initFileDeletionFromS3TasksService()
	a.initFileService()
	a.initFileDeletionTasksService()

	a.initMetricsServer()
	a.startMetricsServer()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	a.startFileDeletionFromS3Workers()
	a.startNotInUseFileDeletionWorkers()
	if err := a.startFileDeletionTasksConsumerWorkers(); err != nil {
		return err
	}
	a.startFileDeletionWorkers()

	<-a.isMetricsStarted
	a.log.Infow("app.started")
	return nil
}

func (a *WorkerApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	if a.workersCancel != nil {
		a.workersCancel()
	}

	workersStopped := make(chan struct{})
	go func() {
		a.workersWg.Wait()
		close(workersStopped)
	}()

	select {
	case <-workersStopped:
	case <-shCtx.Done():
		a.log.Warnw("app.workers_shutdown_timeout")
	}

	if a.metricsServer != nil {
		a.metricsServer.Stop(shCtx)
	}

	a.fileDeletionTasksKafkaClient.Close()

	a.postgresClient.Close()
	a.redisClient.Close()

	a.log.Infow("app.stopped")
}

func (a *WorkerApp) initPostgresClient(ctx context.Context) error {
	pgClient, err := postgres.New(ctx, a.cfg.DB)
	if err != nil {
		a.log.Errorw("app.postgres_client_init_failed", "err", err)
		return err
	}
	a.postgresClient = pgClient

	a.log.Infow("app.postgres_connected")
	return nil
}

func (a *WorkerApp) initRedisClient(ctx context.Context) error {
	redisClient, err := redis.New(ctx, a.cfg.Redis)
	if err != nil {
		a.log.Errorw("app.redis_connect_failed", "err", err)
		return err
	}
	a.redisClient = redisClient

	a.log.Infow("app.redis_connected")
	return nil
}

func (a *WorkerApp) initUserPhotosS3Client(ctx context.Context) error {
	s3Client, err := s3.New(ctx, a.cfg.UserPhotosS3)
	if err != nil {
		a.log.Errorw("app.user_photos_s3_client_init_failed", "err", err)
		return err
	}
	a.userPhotosS3Client = s3Client

	a.log.Infow("app.user_photos_s3_connected")
	return nil
}

func (a *WorkerApp) initChatFilesS3Client(ctx context.Context) error {
	s3Client, err := s3.New(ctx, a.cfg.ChatFilesS3)
	if err != nil {
		a.log.Errorw("app.chat_files_s3_client_init_failed", "err", err)
		return err
	}
	a.chatFilesS3Client = s3Client

	a.log.Infow("app.chat_files_s3_connected")
	return nil
}

func (a *WorkerApp) initFileDeletionTasksKafkaClient(ctx context.Context) error {
	kafkaClient, err := kafka.New(a.cfg.FileDeletionTasksConsumer.KafkaConsumerSettings.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.file_deletion_tasks_kafka_client_init_failed", "err", err)
		return err
	}
	a.fileDeletionTasksKafkaClient = kafkaClient
	return nil
}

func (a *WorkerApp) initFileDeletionFromS3TasksService() {
	a.fileDeletionFromS3TasksService = filedeletionfroms3tasksservice.New(a.postgresClient, a.userPhotosS3Client, a.chatFilesS3Client, a.log)
}

func (a *WorkerApp) initFileService() {
	a.fileService = fileservice.New(a.postgresClient, a.redisClient, a.userPhotosS3Client, a.chatFilesS3Client, a.fileDeletionFromS3TasksService, a.log)
}

func (a *WorkerApp) initFileDeletionTasksService() {
	a.fileDeletionTasksService = filedeletiontasksservice.New(a.postgresClient, a.fileService, a.log)
}

func (a *WorkerApp) initMetricsServer() {
	a.metricsServer = prommetrics.New(a.cfg.MetricsServer)
	a.workerMetrics = prommetrics.NewWorkerMetrics(a.metricsServer.Registry())
}

func (a *WorkerApp) startMetricsServer() {
	go func() {
		a.log.Infow("app.metrics_serve_start", "port", a.cfg.MetricsServer.Port)
		a.isMetricsStarted <- true
		if err := a.metricsServer.Start(); err != nil {
			a.log.Fatalw("app.metrics_serve_error", "err", err)
		}
	}()
}

func (a *WorkerApp) startFileDeletionFromS3Workers() {
	for i := 0; i < int(a.cfg.FileDeletionFromS3Worker.Count); i++ {
		w := filedeletionfroms3worker.New(
			a.cfg.FileDeletionFromS3Worker,
			a.fileDeletionFromS3TasksService,
			a.log,
			a.workerMetrics,
		)
		a.workersWg.Add(1)
		go func() {
			defer a.workersWg.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startNotInUseFileDeletionWorkers() {
	for i := 0; i < int(a.cfg.NotInUseFileDeletionWorker.Count); i++ {
		w := notinusefiledeletionworker.New(
			a.cfg.NotInUseFileDeletionWorker,
			a.fileService,
			a.log,
			a.workerMetrics,
		)
		a.workersWg.Add(1)
		go func() {
			defer a.workersWg.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startFileDeletionTasksConsumerWorkers() error {
	for i := 0; i < int(a.cfg.FileDeletionTasksConsumer.Count); i++ {
		w, err := filedeletiontasksconsumerworker.New(
			a.cfg.FileDeletionTasksConsumer.KafkaConsumerSettings,
			a.fileDeletionTasksKafkaClient,
			a.fileDeletionTasksService,
			a.log,
			a.workerMetrics,
		)
		if err != nil {
			a.log.Errorw("app.file_deletion_tasks_consumer_worker_init_failed", "err", err)
			return err
		}
		a.workersWg.Add(1)
		go func() {
			defer a.workersWg.Done()
			w.Run(a.workersCtx)
		}()
	}
	return nil
}

func (a *WorkerApp) startFileDeletionWorkers() {
	for i := 0; i < int(a.cfg.FileDeletionWorker.Count); i++ {
		w := filedeletionworker.New(
			a.cfg.FileDeletionWorker,
			a.fileDeletionTasksService,
			a.log,
			a.workerMetrics,
		)
		a.workersWg.Add(1)
		go func() {
			defer a.workersWg.Done()
			w.Run(a.workersCtx)
		}()
	}
}
