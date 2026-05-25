package filedeletionfroms3worker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/settings"
	filedeletionfroms3tasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_from_s3_tasks_service"
	prommetrics "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "file_deletion_from_s3"

type FileDeletionFromS3Worker struct {
	workerID                       string
	cfg                            *settings.FileDeletionFromS3WorkerSettings
	fileDeletionFromS3TasksService filedeletionfroms3tasksservice.FileDeletionFromS3TasksService
	log                            *zap.SugaredLogger
	metrics                        *prommetrics.WorkerMetrics
}

func New(
	cfg settings.FileDeletionFromS3WorkerSettings,
	fileDeletionFromS3TasksService filedeletionfroms3tasksservice.FileDeletionFromS3TasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &FileDeletionFromS3Worker{
		workerID:                       uuid.New().String(),
		cfg:                            &cfg,
		fileDeletionFromS3TasksService: fileDeletionFromS3TasksService,
		log:                            log,
		metrics:                        metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *FileDeletionFromS3Worker) Run(ctx context.Context) {
	w.log.Infow("file_deletion_from_s3.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("file_deletion_from_s3.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("file_deletion_from_s3.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *FileDeletionFromS3Worker) runOnce(ctx context.Context) {
	start := time.Now()
	err := w.fileDeletionFromS3TasksService.ProcessFileDeletionFromS3Tasks(ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize), nil)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	if err != nil {
		w.metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
		return
	}
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
