package filedeletionworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/settings"
	filedeletiontasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_tasks_service"
	prommetrics "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "file_deletion"

type FileDeletionWorker struct {
	workerID                 string
	cfg                      *settings.FileDeletionWorkerSettings
	fileDeletionTasksService filedeletiontasksservice.FileDeletionTasksService
	log                      *zap.SugaredLogger
	metrics                  *prommetrics.WorkerMetrics
}

func New(
	cfg settings.FileDeletionWorkerSettings,
	fileDeletionTasksService filedeletiontasksservice.FileDeletionTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &FileDeletionWorker{
		workerID:                 uuid.New().String(),
		cfg:                      &cfg,
		fileDeletionTasksService: fileDeletionTasksService,
		log:                      log,
		metrics:                  metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *FileDeletionWorker) Run(ctx context.Context) {
	w.log.Infow("file_deletion.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("file_deletion.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("file_deletion.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *FileDeletionWorker) runOnce(ctx context.Context) {
	start := time.Now()
	err := w.fileDeletionTasksService.ProcessFileDeletionTasks(ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize), nil)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	if err != nil {
		w.metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
		return
	}
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
