package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/mediaprocessor"
	implkafkaproducer "github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/impl/kafka"
	producersinterfaces "github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/interfaces"
	fileprocessingcompletedeventsservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_completed_events_service"
	fileprocessingtasksservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_tasks_service"
	fileservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_service"
	filegrpcclient "github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/client/grpc/file_client"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/postgres"
	prommetrics "github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/s3"
	fileprocessingworker "github.com/ZaiiiRan/messenger/backend/media-service/internal/workers/file_processing"
	fileprocessingcompletedsendingworker "github.com/ZaiiiRan/messenger/backend/media-service/internal/workers/file_processing_completed_sending"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient

	userPhotosS3Client *s3.S3Client
	chatFilesS3Client  *s3.S3Client

	fileGrpcClient *filegrpcclient.Client

	processingCompletedKafkaClient *kafka.KafkaClient
	completedEventsProducer        producersinterfaces.FileProcessingCompletedEventsProducer

	processor                            mediaprocessor.Processor
	fileService                          fileservice.FileService
	fileProcessingCompletedEventsService fileprocessingcompletedeventsservice.FileProcessingCompletedEventsService
	fileProcessingTasksService           fileprocessingtasksservice.FileProcessingTasksService

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
	if err := a.initFileGrpcClient(ctx); err != nil {
		return err
	}
	if err := a.initUserPhotosS3Client(ctx); err != nil {
		return err
	}
	if err := a.initChatFilesS3Client(ctx); err != nil {
		return err
	}
	if err := a.initProcessingCompletedKafkaClient(); err != nil {
		return err
	}
	if err := a.initCompletedEventsProducer(); err != nil {
		return err
	}

	a.initProcessor()
	a.initFileService()
	a.initFileProcessingCompletedEventsService()
	a.initFileProcessingTasksService()

	a.initMetricsServer()
	a.startMetricsServer()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	a.startFileProcessingWorkers()
	a.startFileProcessingCompletedSendingWorkers()

	<-a.isMetricsStarted
	a.log.Infow("app.started")
	return nil
}

func (a *WorkerApp) Stop(ctx context.Context) error {
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

	if a.completedEventsProducer != nil {
		a.completedEventsProducer.Close()
	}
	if a.processingCompletedKafkaClient != nil {
		a.processingCompletedKafkaClient.Close()
	}

	a.fileGrpcClient.Close()
	a.postgresClient.Close()

	a.log.Infow("app.stopped")
	return nil
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

func (a *WorkerApp) initFileGrpcClient(ctx context.Context) error {
	fileClient, err := filegrpcclient.New(ctx, a.cfg.FileServiceGRPCClient, nil, nil)
	if err != nil {
		a.log.Errorw("app.file_grpc_client_init_failed", "err", err)
		return err
	}
	a.fileGrpcClient = fileClient

	if a.cfg.FileServiceGRPCClient.AutoConnect {
		a.log.Infow("app.file_grpc_client_connected")
	} else {
		a.log.Infow("app.file_grpc_client_initialized")
	}
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

func (a *WorkerApp) initProcessingCompletedKafkaClient() error {
	kafkaClient, err := kafka.New(a.cfg.FileProcessingCompletedEventsProducer.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.processing_completed_kafka_client_init_failed", "err", err)
		return err
	}
	a.processingCompletedKafkaClient = kafkaClient

	a.log.Infow("app.processing_completed_kafka_connected")
	return nil
}

func (a *WorkerApp) initCompletedEventsProducer() error {
	producer, err := implkafkaproducer.New(a.cfg.FileProcessingCompletedEventsProducer, a.processingCompletedKafkaClient)
	if err != nil {
		a.log.Errorw("app.completed_events_producer_init_failed", "err", err)
		return err
	}
	a.completedEventsProducer = producer
	return nil
}

func (a *WorkerApp) initProcessor() {
	a.processor = mediaprocessor.New(
		a.cfg.Processor.FfmpegBin,
		a.cfg.Processor.FfprobeBin,
		a.cfg.Processor.TmpDir,
	)
}

func (a *WorkerApp) initFileService() {
	a.fileService = fileservice.New(a.fileGrpcClient, a.log)
}

func (a *WorkerApp) initFileProcessingCompletedEventsService() {
	a.fileProcessingCompletedEventsService = fileprocessingcompletedeventsservice.New(
		a.postgresClient,
		a.completedEventsProducer,
		a.log,
	)
}

func (a *WorkerApp) initFileProcessingTasksService() {
	a.fileProcessingTasksService = fileprocessingtasksservice.New(
		a.postgresClient,
		a.fileService,
		a.processor,
		a.userPhotosS3Client,
		a.chatFilesS3Client,
		a.fileProcessingCompletedEventsService,
		a.log,
	)
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

func (a *WorkerApp) startFileProcessingWorkers() {
	for i := 0; i < int(a.cfg.FileProcessingWorker.Count); i++ {
		w := fileprocessingworker.New(
			a.cfg.FileProcessingWorker,
			a.fileProcessingTasksService,
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

func (a *WorkerApp) startFileProcessingCompletedSendingWorkers() {
	for i := 0; i < int(a.cfg.FileProcessingCompletedSendingWorker.Count); i++ {
		w := fileprocessingcompletedsendingworker.New(
			a.cfg.FileProcessingCompletedSendingWorker,
			a.fileProcessingCompletedEventsService,
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
