package userdatadeletiontasksworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_data_deletion_tasks"
	prommetrics "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "user_data_deletion_tasks"

type UserDataDeletionTasksWorker struct {
	workerID                     string
	cfg                          *settings.UserDataDeletionTasksWorkerSettings
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService
	log                          *zap.SugaredLogger
	metrics                      *prommetrics.WorkerMetrics
}

func New(
	cfg settings.UserDataDeletionTasksWorkerSettings,
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &UserDataDeletionTasksWorker{
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

func (w *UserDataDeletionTasksWorker) Run(ctx context.Context) {
	w.log.Infow("user_data_deletion_tasks_worker.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("user_data_deletion_tasks_worker.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("user_data_deletion_tasks_worker.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *UserDataDeletionTasksWorker) runOnce(ctx context.Context) {
	start := time.Now()
	w.userDataDeletionTasksService.ProcessUserDataDeletionTasks(
		ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize),
	)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
