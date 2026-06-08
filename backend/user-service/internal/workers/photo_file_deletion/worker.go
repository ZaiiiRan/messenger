package photofiledeletionworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	photofiledeletionsvc "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/photo_file_deletion"
	prommetrics "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "photo_file_deletion"

type PhotoFileDeletionWorker struct {
	workerID                string
	cfg                     *settings.KafkaSendingWorkerSettings
	photoFileDeletionService photofiledeletionsvc.PhotoFileDeletionService
	log                     *zap.SugaredLogger
	metrics                 *prommetrics.WorkerMetrics
}

func New(
	cfg settings.KafkaSendingWorkerSettings,
	photoFileDeletionService photofiledeletionsvc.PhotoFileDeletionService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &PhotoFileDeletionWorker{
		workerID:                uuid.New().String(),
		cfg:                     &cfg,
		photoFileDeletionService: photoFileDeletionService,
		log:                     log,
		metrics:                 metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *PhotoFileDeletionWorker) Run(ctx context.Context) {
	w.log.Infow("photo_file_deletion.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("photo_file_deletion.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("photo_file_deletion.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *PhotoFileDeletionWorker) runOnce(ctx context.Context) {
	start := time.Now()
	err := w.photoFileDeletionService.SendTasks(ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize))
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		w.metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
		return
	}
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
