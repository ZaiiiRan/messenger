package notinusefiledeletionworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/settings"
	fileservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_service"
	prommetrics "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "not_in_use_file_deletion"

type NotInUseFileDeletionWorker struct {
	workerID    string
	cfg         *settings.NotInUseFileDeletionWorkerSettings
	fileService fileservice.FileService
	log         *zap.SugaredLogger
	metrics     *prommetrics.WorkerMetrics
}

func New(
	cfg settings.NotInUseFileDeletionWorkerSettings,
	fileService fileservice.FileService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &NotInUseFileDeletionWorker{
		workerID:    uuid.New().String(),
		cfg:         &cfg,
		fileService: fileService,
		log:         log,
		metrics:     metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *NotInUseFileDeletionWorker) Run(ctx context.Context) {
	w.log.Infow("not_in_use_file_deletion.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("not_in_use_file_deletion.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("not_in_use_file_deletion.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *NotInUseFileDeletionWorker) runOnce(ctx context.Context) {
	start := time.Now()
	w.fileService.DeleteNotInUseFiles(ctx, w.workerID, int(w.cfg.BatchSize), nil)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
