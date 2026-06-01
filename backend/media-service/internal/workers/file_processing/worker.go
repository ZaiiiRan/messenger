package fileprocessingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/settings"
	fileprocessingtasksservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_tasks_service"
	prommetrics "github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "file_processing"

type FileProcessingWorker struct {
	workerID                   string
	cfg                        *settings.FileProcessingWorkerSettings
	fileProcessingTasksService fileprocessingtasksservice.FileProcessingTasksService
	log                        *zap.SugaredLogger
	metrics                    *prommetrics.WorkerMetrics
}

func New(
	cfg settings.FileProcessingWorkerSettings,
	fileProcessingTasksService fileprocessingtasksservice.FileProcessingTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &FileProcessingWorker{
		workerID:                   uuid.New().String(),
		cfg:                        &cfg,
		fileProcessingTasksService: fileProcessingTasksService,
		log:                        log,
		metrics:                    metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *FileProcessingWorker) Run(ctx context.Context) {
	w.log.Infow("file_processing.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("file_processing.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("file_processing.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *FileProcessingWorker) runOnce(ctx context.Context) {
	start := time.Now()
	err := w.fileProcessingTasksService.ProcessFileProcessingTasks(ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize), nil)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	if err != nil {
		w.metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
		return
	}
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
