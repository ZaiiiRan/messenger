package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config"
	implkafkaproducer "github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/impl/kafka"
	photofiledeletionsvc "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/photo_file_deletion"
	updatefriendsinprivacysettingstasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/update_friends_in_privacy_settings_tasks"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user_data_deletion_tasks"
	filegrpcclient "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/client/grpc/file_client"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	prommetrics "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/redis"
	deleteduserdataclearingworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/deleted_user_data_clearing"
	photofiledeletionworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/photo_file_deletion"
	unbantemporarilybanneduserworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/unban_temporarily_banned_user"
	unconfirmeduserdataclearingworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/unconfirmed_user_data_clearing"
	updatingfriendsinprivacysettingsworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/updating_friends_in_privacy_settings"
	usersdatadeletiontaskssendingworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/user_data_deletion_tasks_sending"
	userrelationshipchangestasksconsumerworker "github.com/ZaiiiRan/messenger/backend/user-service/internal/workers/user_relationship_changes_tasks_consumer"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	postgresClient *postgres.PostgresClient
	redisClient    *redis.RedisClient
	fileGrpcClient *filegrpcclient.Client

	userDataDeletionTasksKafkaClient   *kafka.KafkaClient
	userDataDeletionTasksKafkaProducer *implkafkaproducer.Producer

	userRelationshipChangesTasksKafkaClient *kafka.KafkaClient

	userService                           userservice.UserService
	userDataDeletionTasksService          userdatadeletiontasksservice.UserDataDeletionTasksService
	photoFileDeletionService              photofiledeletionsvc.PhotoFileDeletionService
	updateFriendsInPrivacySettingsService updatefriendsinprivacysettingstasksservice.UpdateFriendsInPrivacySettingsTasksService

	metricsServer    *prommetrics.Server
	workerMetrics    *prommetrics.WorkerMetrics
	isMetricsStarted chan bool

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
	// File client is optional — if it fails we skip photo-related features.
	if err := a.initFileGrpcClient(ctx); err != nil {
		a.log.Warnw("app.file_grpc_client_init_failed_skipping", "err", err)
	}
	if err := a.initUserDataDeletionTasksKafkaClient(ctx); err != nil {
		return err
	}
	if err := a.initUserRelationshipChangesTasksKafkaClient(ctx); err != nil {
		return err
	}
	if err := a.initUserDataDeletionTasksKafkaProducer(); err != nil {
		return err
	}

	a.initUserDataDeletionTasksService()
	a.initPhotoFileDeletionService()
	a.initUserService()
	a.initUpdateFriendsInPrivacySettingsService()

	a.initMetricsServer()
	a.startMetricsServer()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	a.startUnconfirmedUsersDataClearingWorkers()
	a.startDeletedUsersDataClearingWorkers()
	a.startUnbanTemporarilyBannedUserWorkers()
	a.startUserDataDeletionTasksSendingWorkers()
	a.startPhotoFileDeletionWorkers()
	a.startUserRelationshipChangesTasksConsumerWorkers()
	a.startUpdatingFriendsInPrivacySettingsWorkers()

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
		a.workersWG.Wait()
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

	a.userDataDeletionTasksKafkaProducer.Close()
	a.userRelationshipChangesTasksKafkaClient.Close()
	if a.fileGrpcClient != nil {
		a.fileGrpcClient.Close()
	}

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

func (a *WorkerApp) initUserRelationshipChangesTasksKafkaClient(ctx context.Context) error {
	kafkaClient, err := kafka.New(a.cfg.UserRelationshipChangesTasksConsumer.KafkaConsumerSettings.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.user_relationship_changes_tasks_kafka_client_init_failed", "err", err)
		return err
	}
	a.userRelationshipChangesTasksKafkaClient = kafkaClient

	a.log.Infow("app.user_relationship_changes_tasks_kafka_connected")
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

func (a *WorkerApp) initUpdateFriendsInPrivacySettingsService() {
	a.updateFriendsInPrivacySettingsService = updatefriendsinprivacysettingstasksservice.New(a.postgresClient, a.userService, a.log)
}

func (a *WorkerApp) initPhotoFileDeletionService() {
	if a.fileGrpcClient != nil {
		a.photoFileDeletionService = photofiledeletionsvc.New(a.postgresClient, a.fileGrpcClient, a.log)
	}
}

func (a *WorkerApp) initUserService() {
	a.userService = userservice.New(a.postgresClient, a.redisClient, a.log, a.userDataDeletionTasksService, a.photoFileDeletionService)
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

func (a *WorkerApp) startUnconfirmedUsersDataClearingWorkers() {
	for i := 0; i < int(a.cfg.UnconfirmedUsersDataClearingWorker.Count); i++ {
		w := unconfirmeduserdataclearingworker.New(
			a.cfg.UnconfirmedUsersDataClearingWorker,
			a.userService,
			a.log,
			a.workerMetrics,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startDeletedUsersDataClearingWorkers() {
	for i := 0; i < int(a.cfg.DeletedUsersDataClearingWorker.Count); i++ {
		w := deleteduserdataclearingworker.New(
			a.cfg.DeletedUsersDataClearingWorker,
			a.userService,
			a.log,
			a.workerMetrics,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startUnbanTemporarilyBannedUserWorkers() {
	for i := 0; i < int(a.cfg.UnbanTemporarilyBannedUsersWorker.Count); i++ {
		w := unbantemporarilybanneduserworker.New(
			a.cfg.UnbanTemporarilyBannedUsersWorker,
			a.userService,
			a.log,
			a.workerMetrics,
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
			a.workerMetrics,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startUserRelationshipChangesTasksConsumerWorkers() error {
	for i := 0; i < int(a.cfg.UserRelationshipChangesTasksConsumer.Count); i++ {
		w, err := userrelationshipchangestasksconsumerworker.New(
			a.cfg.UserRelationshipChangesTasksConsumer.KafkaConsumerSettings,
			a.userDataDeletionTasksKafkaClient,
			a.updateFriendsInPrivacySettingsService,
			a.log,
			a.workerMetrics,
		)
		if err != nil {
			a.log.Errorw("app.user_relationship_changes_tasks_consumer_worker_init_failed", "err", err, "worker_id", i)
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

func (a *WorkerApp) initFileGrpcClient(ctx context.Context) error {
	cl, err := filegrpcclient.New(ctx, a.cfg.FileServiceGRPCClient, nil, nil)
	if err != nil {
		return err
	}
	a.fileGrpcClient = cl
	a.log.Infow("app.file_grpc_client_initialized")
	return nil
}

func (a *WorkerApp) startPhotoFileDeletionWorkers() {
	if a.photoFileDeletionService == nil {
		return
	}
	for i := 0; i < int(a.cfg.PhotoFileDeletionWorker.Count); i++ {
		w := photofiledeletionworker.New(
			a.cfg.PhotoFileDeletionWorker,
			a.photoFileDeletionService,
			a.log,
			a.workerMetrics,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}

func (a *WorkerApp) startUpdatingFriendsInPrivacySettingsWorkers() {
	for i := 0; i < int(a.cfg.UpdateFriendsInPrivacySettingsWorker.Count); i++ {
		w := updatingfriendsinprivacysettingsworker.New(
			a.cfg.UpdateFriendsInPrivacySettingsWorker,
			a.updateFriendsInPrivacySettingsService,
			a.log,
			a.workerMetrics,
		)
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
}
