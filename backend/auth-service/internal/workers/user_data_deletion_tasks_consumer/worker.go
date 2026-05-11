package userdatadeletiontasksconsumerworker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	kafkaconsumer "github.com/ZaiiiRan/messenger/backend/auth-service/internal/consumers/impl/kafka"
	consumersmodels "github.com/ZaiiiRan/messenger/backend/auth-service/internal/consumers/models"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_data_deletion_tasks"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/kafka"
	prommetrics "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "user_data_deletion_tasks_consumer"

type UserDataDeletionTasksConsumerWorker struct {
	workerID                     string
	consumer                     *kafkaconsumer.Consumer
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService
	log                          *zap.SugaredLogger
}

func New(
	cfg settings.KafkaConsumerSettings,
	kafkaClient *kafkatransport.KafkaClient,
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) (workers.Worker, error) {
	id := uuid.New().String()

	workerLog := log.With("worker_id", id)

	w := &UserDataDeletionTasksConsumerWorker{
		workerID:                     id,
		log:                          log,
		userDataDeletionTasksService: userDataDeletionTasksService,
	}

	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(0)

	handlerFunc := func(ctx context.Context, messages []kafkaconsumer.Message) error {
		start := time.Now()

		taskMessages := make([]*consumersmodels.UserDataDeletionTask, 0, len(messages))
		for _, message := range messages {
			var taskMessage *consumersmodels.UserDataDeletionTask
			if err := json.Unmarshal([]byte(message.Body), &taskMessage); err != nil {
				w.log.Warnw("user_data_deletion_tasks_consumer_handler.unmarshal_failed", "err", err, "body", message.Body)
				continue
			}
			if taskMessage == nil {
				w.log.Warnw("user_data_deletion_tasks_consumer_handler.unmarshal_failed", "err", "null message body")
				continue
			}
			if taskMessage.Id == "" {
				w.log.Warnw("user_data_deletion_tasks_consumer_handler.handle_failed", "err", "empty id")
				continue
			}

			taskMessages = append(taskMessages, taskMessage)
		}

		err := w.userDataDeletionTasksService.CreateUserDataDeletionTasks(ctx, w.workerID, taskMessages)
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
		return nil, fmt.Errorf("user data deletion tasks consumer worker %s: %w", id, err)
	}
	w.consumer = consumer
	return w, nil
}

func (w *UserDataDeletionTasksConsumerWorker) Run(ctx context.Context) {
	w.log.Infow("user_data_deletion_conumser.started")
	w.consumer.Run(ctx)
	w.log.Infow("user_data_deletion_conumser.stopped")
}
