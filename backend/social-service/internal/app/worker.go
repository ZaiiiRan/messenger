package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/config"
	implkafkaproducer "github.com/ZaiiiRan/messenger/backend/social-service/internal/producers/impl/kafka"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_data_deletion_tasks"
	userrelationshipservice "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_relationship"
	userrelationshipchangestasks "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_relationship_changes_tasks"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
	prommetrics "github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/redis"
	userdatadeletionworker "github.com/ZaiiiRan/messenger/backend/social-service/internal/workers/user_data_deletion"
	userdatadeletiontasksconsumerworker "github.com/ZaiiiRan/messenger/backend/social-service/internal/workers/user_data_deletion_tasks_consumer"
	userrelationshipchangestaskssendingworker "github.com/ZaiiiRan/messenger/backend/social-service/internal/workers/user_relationship_changes_tasks_sending"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient

	userRelationshipChangesTasksKafkaClient   *kafka.KafkaClient
	userRelationshipChangesTasksKafkaProducer *implkafkaproducer.Producer

	userDataDeletionTasksKafkaClient *kafka.KafkaClient

	userRelationshipService             userrelationshipservice.UserRelationshipService
	userRelationshipChangesTasksService userrelationshipchangestasks.UserRelationshipChangesTasksService
	userDataDeletionTasksService        userdatadeletiontasksservice.UserDataDeletionTasksService

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
	if err := a.initUserRelationshipChangesTasksKafkaClient(ctx); err != nil {
		return err
	}
	if err := a.initUserRelationshipChangesTasksKafkaProducer(ctx); err != nil {
		return err
	}
	if err := a.initUserDataDeletionTasksKafkaClient(ctx); err != nil {
		return err
	}

	a.initUserRelationshipService()
	a.initUserRelationshipChangesTasksService()
	a.initUserDataDeletionTasksService()

	a.initMetricsServer()
	a.startMetricsServer()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	a.startUserRelationshipChangesTasksWorkers()
	if err := a.startUserDataDeletionTasksConsumerWorkers(); err != nil {
		return err
	}
	a.startUserDataDeletionWorkers()

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

	a.userRelationshipChangesTasksKafkaProducer.Close()
	a.userDataDeletionTasksKafkaClient.Close()

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

	a.log.Infow("app.postgres_connectd")
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

func (a *WorkerApp) initUserRelationshipChangesTasksKafkaClient(ctx context.Context) error {
	kafkaClient, err := kafka.New(a.cfg.UserRelationshipChangesTasksProducer.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.user_relationship_changes_tasks_kafka_client_init_failed", "err", err)
		return err
	}
	a.userRelationshipChangesTasksKafkaClient = kafkaClient

	a.log.Infow("app.user_relationship_changes_tasks_kafka_connected")
	return nil
}

func (a *WorkerApp) initUserRelationshipChangesTasksKafkaProducer(ctx context.Context) error {
	producer, err := implkafkaproducer.New(a.cfg.UserRelationshipChangesTasksProducer, a.userRelationshipChangesTasksKafkaClient)
	if err != nil {
		a.log.Errorw("app.user_relationship_changes_tasks_producer_init_failed", "err", err)
		return err
	}
	a.userRelationshipChangesTasksKafkaProducer = producer
	return nil
}

func (a *WorkerApp) initUserDataDeletionTasksKafkaClient(ctx context.Context) error {
	kafkaClient, err := kafka.New(a.cfg.UserDataDeletionTasksConsumer.KafkaConsumerSettings.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.user_data_deletion_tasks_kafka_client_init_failed", "err", err)
		return err
	}
	a.userDataDeletionTasksKafkaClient = kafkaClient
	return nil
}

func (a *WorkerApp) initUserRelationshipService() {
	a.userRelationshipService = userrelationshipservice.New(a.postgresClient, a.redisClient, a.log)
}

func (a *WorkerApp) initUserRelationshipChangesTasksService() {
	a.userRelationshipChangesTasksService = userrelationshipchangestasks.New(a.postgresClient, a.userRelationshipChangesTasksKafkaProducer, a.log)
}

func (a *WorkerApp) initUserDataDeletionTasksService() {
	a.userDataDeletionTasksService = userdatadeletiontasksservice.New(a.postgresClient, a.userRelationshipService, a.log)
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

func (a *WorkerApp) startUserRelationshipChangesTasksWorkers() {
	for i := 0; i < int(a.cfg.UserRelationshipChangesTasksSendindWorker.Count); i++ {
		w := userrelationshipchangestaskssendingworker.New(
			a.cfg.UserRelationshipChangesTasksSendindWorker,
			a.userRelationshipChangesTasksService,
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

func (a *WorkerApp) startUserDataDeletionTasksConsumerWorkers() error {
	for i := 0; i < int(a.cfg.UserDataDeletionTasksConsumer.Count); i++ {
		w, err := userdatadeletiontasksconsumerworker.New(
			a.cfg.UserDataDeletionTasksConsumer.KafkaConsumerSettings,
			a.userDataDeletionTasksKafkaClient,
			a.userDataDeletionTasksService,
			a.log,
			a.workerMetrics,
		)
		if err != nil {
			a.log.Errorw("app.user_data_deletion_tasks_consumer_worker_init_failed", "err", err)
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

func (a *WorkerApp) startUserDataDeletionWorkers() {
	for i := 0; i < int(a.cfg.UserDataDeletionWorker.Count); i++ {
		w := userdatadeletionworker.New(
			a.cfg.UserDataDeletionWorker,
			a.userDataDeletionTasksService,
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
