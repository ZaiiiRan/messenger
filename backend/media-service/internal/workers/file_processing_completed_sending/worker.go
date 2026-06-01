package fileprocessingcompletedsendingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/settings"
	fileprocessingcompletedeventsservice "github.com/ZaiiiRan/messenger/backend/media-service/internal/services/file_processing_completed_events_service"
	prommetrics "github.com/ZaiiiRan/messenger/backend/media-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "file_processing_completed_sending"

type FileProcessingCompletedSendingWorker struct {
	workerID                             string
	cfg                                  *settings.KafkaSendingWorkerSettings
	fileProcessingCompletedEventsService fileprocessingcompletedeventsservice.FileProcessingCompletedEventsService
	log                                  *zap.SugaredLogger
	metrics                              *prommetrics.WorkerMetrics
}

func New(
	cfg settings.KafkaSendingWorkerSettings,
	fileProcessingCompletedEventsService fileprocessingcompletedeventsservice.FileProcessingCompletedEventsService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &FileProcessingCompletedSendingWorker{
		workerID:                             uuid.New().String(),
		cfg:                                  &cfg,
		fileProcessingCompletedEventsService: fileProcessingCompletedEventsService,
		log:                                  log,
		metrics:                              metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *FileProcessingCompletedSendingWorker) Run(ctx context.Context) {
	w.log.Infow("file_processing_completed_sending.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("file_processing_completed_sending.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("file_processing_completed_sending.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *FileProcessingCompletedSendingWorker) runOnce(ctx context.Context) {
	start := time.Now()
	err := w.fileProcessingCompletedEventsService.SendFileProcessingCompletedEvents(ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize))
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	if err != nil {
		w.metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
		return
	}
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
