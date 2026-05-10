package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	codeservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_data_deletion_tasks"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	expiredtokenclearingworker "github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers/expired_token_clearing"
	userdatadeletiontasksworker "github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers/user_data_deletion_tasks"
	userdatadeletiontasksconsumerworker "github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers/user_data_deletion_tasks_consumer"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	postgresClient                   *postgres.PostgresClient
	redisClient                      *redis.RedisClient
	userDataDeletionTasksKafkaClient *kafkatransport.KafkaClient

	codeService                  codeservice.CodeService
	passwordService              passwordservice.PasswordService
	tokenService                 tokenservice.TokenService
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
	if err := a.initUserDataDeletionTasksKafkaClient(); err != nil {
		return err
	}

	a.initPasswordService()
	a.initCodeService()
	a.initTokenService()
	a.initUserDataDeletionTasksSerivce()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	a.startExpiredTokenClearingWorkers()
	if err := a.startUserDataDeletionTasksConsumerWorkers(); err != nil {
		return err
	}
	a.startUserDataDeletionTasksWorkers()

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

	a.postgresClient.Close()
	a.redisClient.Close()
	a.userDataDeletionTasksKafkaClient.Close()

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

func (a *WorkerApp) initUserDataDeletionTasksKafkaClient() error {
	kafkaClient, err := kafkatransport.New(a.cfg.UserDataDeletionTasksConsumer.KafkaConsumerSettings.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.kafka_connect_failed", "err", err)
		return err
	}
	a.userDataDeletionTasksKafkaClient = kafkaClient
	a.log.Infow("app.kafka_connected")
	return nil
}

func (a *WorkerApp) initCodeService() {
	a.codeService = codeservice.New(a.postgresClient, a.redisClient, a.log)
}

func (a *WorkerApp) initPasswordService() {
	a.passwordService = passwordservice.New(a.postgresClient, a.redisClient, a.log)
}

func (a *WorkerApp) initTokenService() {
	a.tokenService = tokenservice.New(settings.JWTSettings{}, a.postgresClient, a.redisClient, a.log)
}

func (a *WorkerApp) initUserDataDeletionTasksSerivce() {
	a.userDataDeletionTasksService = userdatadeletiontasksservice.New(a.postgresClient, a.passwordService, a.codeService, a.tokenService, a.log)
}

func (a *WorkerApp) startExpiredTokenClearingWorkers() {
	for i := 0; i < int(a.cfg.ExpiredTokenClearingWorker.Count); i++ {
		w := expiredtokenclearingworker.New(a.cfg.ExpiredTokenClearingWorker, a.tokenService, a.postgresClient, a.log)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
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
		)
		if err != nil {
			a.log.Errorw("app.user_data_deletion_tasks_consumer_worker_init_failed", "err", err, "worker_id", i)
			return err
		}
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
	return nil
}

func (a *WorkerApp) startUserDataDeletionTasksWorkers() {
	for i := 0; i < int(a.cfg.UserDataDeletionTasksWorker.Count); i++ {
		w := userdatadeletiontasksworker.New(
			a.cfg.UserDataDeletionTasksWorker,
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
