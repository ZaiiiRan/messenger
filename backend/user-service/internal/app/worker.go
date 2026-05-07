package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config"
	implkafkaproducer "github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/impl/kafka"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user_data_deletion_tasks"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/redis"
	unconfirmeduserdataclearingworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/unconfirmed_user_data_clearing_worker"
	usersdatadeletiontaskssendingworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/user_data_deletion_tasks_sending_worker"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient

	userDataDeletionTasksKafkaClient   *kafka.KafkaClient
	userDataDeletionTasksKafkaProducer *implkafkaproducer.Producer

	userService                  userservice.UserService
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService

	workersCtx    context.Context
	workersCancel context.CancelFunc
	workersWG     sync.WaitGroup
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
		cfg: cfg,
		log: log,
	}, nil
}

func (a *WorkerApp) Run(ctx context.Context) error {
	if err := a.initPostgresClient(ctx); err != nil {
		return err
	}
	if err := a.initRedisClient(ctx); err != nil {
		return err
	}
	if err := a.initUserDataDeletionTasksKafkaClient(ctx); err != nil {
		return err
	}
	if err := a.initUserDataDeletionTasksKafkaProducer(); err != nil {
		return err
	}

	a.initUserDataDeletionTasksService()
	a.initUserService()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	a.startUnconfirmedUsersDataClearingWorkers()
	a.startUserDataDeletionTasksSendingWorkers()

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
		a.workersWG.Wait()
		close(workersStopped)
	}()

	select {
	case <-workersStopped:
	case <-shCtx.Done():
		a.log.Warnw("app.workers_shutdown_timeout")
	}

	a.userDataDeletionTasksKafkaProducer.Close()

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

func (a *WorkerApp) initUserDataDeletionTasksKafkaClient(ctx context.Context) error {
	kafkaClient, err := kafka.New(a.cfg.UserDataDeletionTasksProducer.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.user_data_deletion_tasks_kafka_client_init_failed", "err", err)
		return err
	}
	a.userDataDeletionTasksKafkaClient = kafkaClient

	a.log.Infow("app.user_data_deletion_tasks_kafka_connected")
	return nil
}

func (a *WorkerApp) initUserDataDeletionTasksKafkaProducer() error {
	producer, err := implkafkaproducer.New(a.cfg.UserDataDeletionTasksProducer, a.userDataDeletionTasksKafkaClient)
	if err != nil {
		a.log.Errorw("app.user_data_deletion_tasks_kafka_producer_init_failed", "err", err)
		return err
	}
	a.userDataDeletionTasksKafkaProducer = producer
	return nil
}

func (a *WorkerApp) initUserDataDeletionTasksService() {
	a.userDataDeletionTasksService = userdatadeletiontasksservice.New(a.postgresClient, a.userDataDeletionTasksKafkaProducer, a.log)
}

func (a *WorkerApp) initUserService() {
	a.userService = userservice.New(a.postgresClient, a.redisClient, a.log, a.userDataDeletionTasksService)
}

func (a *WorkerApp) startUnconfirmedUsersDataClearingWorkers() {
	for i := 0; i < int(a.cfg.UnconfirmedUsersDataClearingWorker.Count); i++ {
		w := unconfirmeduserdataclearingworker.New(
			a.cfg.UnconfirmedUsersDataClearingWorker,
			a.userService,
			a.log,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startUserDataDeletionTasksSendingWorkers() {
	for i := 0; i < int(a.cfg.UserDataDeletionTasksSendingWorker.Count); i++ {
		w := usersdatadeletiontaskssendingworker.New(
			a.cfg.UserDataDeletionTasksSendingWorker,
			a.userDataDeletionTasksService,
			a.log,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}
