package userrelationshipchangestasksconsumerworker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	kafkaconsumer "github.com/ZaiiiRan/messenger/backend/user-service/internal/consumers/impl/kafka"
	consumermodels "github.com/ZaiiiRan/messenger/backend/user-service/internal/consumers/models"
	updatefriendsinprivacysettingstasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/update_friends_in_privacy_settings_tasks"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/kafka"
	prommetrics "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "user_relationships_changes_tasks_consumer"

type UserRelationshipChangesTasksConsumerWorker struct {
	workerId                                   string
	consumer                                   *kafkaconsumer.Consumer
	updateFriendsInPrivacySettingsTasksService updatefriendsinprivacysettingstasksservice.UpdateFriendsInPrivacySettingsTasksService
	log                                        *zap.SugaredLogger
}

func New(
	cfg settings.KafkaConsumerSettings,
	kafkaClient *kafkatransport.KafkaClient,
	updateFriendsInPrivacySettingsTasksService updatefriendsinprivacysettingstasksservice.UpdateFriendsInPrivacySettingsTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) (workers.Worker, error) {
	id := uuid.New().String()

	workerLog := log.With("worker_id", id)

	w := &UserRelationshipChangesTasksConsumerWorker{
		workerId: id,
		log:      log,
		updateFriendsInPrivacySettingsTasksService: updateFriendsInPrivacySettingsTasksService,
	}

	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(0)

	handlerFunc := func(ctx context.Context, messages []kafkaconsumer.Message) error {
		start := time.Now()

		taskMessages := make([]*consumermodels.UserRelationshipChangeTask, 0, len(messages))
		for _, message := range messages {
			var taskMessage *consumermodels.UserRelationshipChangeTask
			if err := json.Unmarshal([]byte(message.Body), &taskMessage); err != nil {
				w.log.Warnw("user_relationship_changes_tasks_consumer_handler.unmarshal_failed", "err", err)
				continue
			}
			if taskMessage == nil {
				w.log.Warnw("user_relationship_changes_tasks_consumer_handler.unmarshal_failed", "err", "null message body")
				continue
			}
			if taskMessage.Id == "" {
				w.log.Warnw("user_relationship_changes_tasks_consumer_handler.handle_failed", "err", "empty id")
				continue
			}

			taskMessages = append(taskMessages, taskMessage)
		}

		err := w.updateFriendsInPrivacySettingsTasksService.CreateUpdateFriendsInPrivacySettingsTasks(ctx, w.workerId, taskMessages, nil)
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
		return nil, fmt.Errorf("user relationship changes tasks consumer worker %s: %w", id, err)
	}
	w.consumer = consumer
	return w, nil
}

func (w *UserRelationshipChangesTasksConsumerWorker) Run(ctx context.Context) {
	w.log.Infow("user_relationship_changes_tasks_consumer.started")
	w.consumer.Run(ctx)
	w.log.Infow("user_relationship_changes_tasks_consumer.stopped")
}
