package deleteduserdataclearingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	prommetrics "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "deleted_user_data_clearing"

type DeletedUserDataClearingWorker struct {
	workerID    string
	cfg         *settings.DeletedUsersDataClearingWorkerSettings
	userService userservice.UserService
	log         *zap.SugaredLogger
	metrics     *prommetrics.WorkerMetrics
}

func New(
	cfg settings.DeletedUsersDataClearingWorkerSettings,
	userService userservice.UserService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &DeletedUserDataClearingWorker{
		workerID:    uuid.New().String(),
		cfg:         &cfg,
		userService: userService,
		log:         log,
		metrics:     metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(0)
	return w
}

func (w *DeletedUserDataClearingWorker) Run(ctx context.Context) {
	w.log.Infow("deleted_user_data_clearing.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("deleted_user_data_clearing.stopped", "worker_id", w.workerID)
			return
		default:
		}

		deletedCount, err := w.runOnce(ctx)

		var timer *time.Timer
		if deletedCount != int(w.cfg.BatchSize) && err == nil {
			timer = time.NewTimer(time.Millisecond * time.Duration(w.cfg.NoDataIntervalMS))
		} else {
			timer = time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		}
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("deleted_user_data_clearing.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *DeletedUserDataClearingWorker) runOnce(ctx context.Context) (int, error) {
	start := time.Now()
	deletedCount, err := w.userService.ClearDeletedUsers(ctx, int(w.cfg.BatchSize), w.workerID)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	if err != nil {
		if ctx.Err() != nil {
			return 0, nil
		}
		w.metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
		return 0, err
	}
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
	w.metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(float64(deletedCount))
	return deletedCount, nil
}
