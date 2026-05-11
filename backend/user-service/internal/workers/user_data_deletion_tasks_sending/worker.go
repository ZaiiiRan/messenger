package usersdatadeletiontaskssendingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user_data_deletion_tasks"
	prommetrics "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "user_data_deletion_tasks_sending"

type UserDataDeletionTasksSendingWorker struct {
	workerID                     string
	cfg                          *settings.KafkaSendingWorkerSettings
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService
	log                          *zap.SugaredLogger
	metrics                      *prommetrics.WorkerMetrics
}

func New(
	cfg settings.KafkaSendingWorkerSettings,
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &UserDataDeletionTasksSendingWorker{
		workerID:                     uuid.New().String(),
		cfg:                          &cfg,
		userDataDeletionTasksService: userDataDeletionTasksService,
		log:                          log,
		metrics:                      metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *UserDataDeletionTasksSendingWorker) Run(ctx context.Context) {
	w.log.Infow("user_data_deletion_tasks_sending.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("user_data_deletion_tasks_sending.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("user_data_deletion_tasks_sending.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *UserDataDeletionTasksSendingWorker) runOnce(ctx context.Context) {
	start := time.Now()
	w.userDataDeletionTasksService.SendUserDataDeletionTasks(
		ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize),
	)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
