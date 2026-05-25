package filedeletiontasksconsumerworker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/settings"
	kafkaconsumer "github.com/ZaiiiRan/messenger/backend/file-service/internal/consumers/impl/kafka"
	consumermodels "github.com/ZaiiiRan/messenger/backend/file-service/internal/consumers/models"
	filedeletiontasksservice "github.com/ZaiiiRan/messenger/backend/file-service/internal/services/file_deletion_tasks_service"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/kafka"
	prommetrics "github.com/ZaiiiRan/messenger/backend/file-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "file_deletion_tasks_consumer"

type FileDeletionTasksConsumerWorker struct {
	workerID                 string
	consumer                 *kafkaconsumer.Consumer
	fileDeletionTasksService filedeletiontasksservice.FileDeletionTasksService
	log                      *zap.SugaredLogger
}

func New(
	cfg settings.KafkaConsumerSettings,
	kafkaClient *kafkatransport.KafkaClient,
	fileDeletionTasksService filedeletiontasksservice.FileDeletionTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) (workers.Worker, error) {
	id := uuid.New().String()

	workerLog := log.With("worker_id", id)

	w := &FileDeletionTasksConsumerWorker{
		workerID:                 id,
		log:                      workerLog,
		fileDeletionTasksService: fileDeletionTasksService,
	}

	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(0)

	handlerFunc := func(ctx context.Context, messages []kafkaconsumer.Message) error {
		start := time.Now()

		taskMessages := make([]*consumermodels.FileDeletionTask, 0, len(messages))
		for _, msg := range messages {
			var taskMessage *consumermodels.FileDeletionTask
			if err := json.Unmarshal([]byte(msg.Body), &taskMessage); err != nil {
				w.log.Warnw("file_deletion_tasks_consumer_handler.unmarshal_failed", "err", err)
				continue
			}
			if taskMessage == nil {
				w.log.Warnw("file_deletion_tasks_consumer_handler.unmarshal_failed", "err", "nil task message")
				continue
			}
			if taskMessage.Id == "" {
				w.log.Warnw("file_deletion_tasks_consumer_handler.unmarshal_failed", "err", "empty task message id")
				continue
			}

			taskMessages = append(taskMessages, taskMessage)
		}

		err := w.fileDeletionTasksService.CreateFileDeletionTasks(ctx, w.workerID, taskMessages, nil)
		metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
		if err != nil {
			metrics.CyclesTotal.WithLabelValues(workerType, "error").Inc()
			return err
		}
		metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
		metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(float64(len(taskMessages)))
		return nil
	}

	consumer, err := kafkaconsumer.New(cfg, kafkaClient, workerLog, handlerFunc)
	if err != nil {
		return nil, fmt.Errorf("file deletion tasks consumer worker %s: %w", id, err)
	}
	w.consumer = consumer
	return w, nil
}

func (w *FileDeletionTasksConsumerWorker) Run(ctx context.Context) {
	w.log.Infow("file_deletion_tasks_consumer_worker.started")
	w.consumer.Run(ctx)
	w.log.Infow("file_deletion_tasks_consumer_worker.stopped")
}
