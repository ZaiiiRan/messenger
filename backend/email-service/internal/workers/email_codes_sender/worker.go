package emailcodessenderworker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	kafkaconsumer "github.com/ZaiiiRan/messenger/backend/email-service/internal/consumers/impl/kafka"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/consumers/models"
	senderservice "github.com/ZaiiiRan/messenger/backend/email-service/internal/services/sender"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/kafka"
	prommetrics "github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "email_codes_sender"

type EmailCodesSenderWorker struct {
	id            string
	consumer      *kafkaconsumer.Consumer
	senderService senderservice.SenderService
	log           *zap.SugaredLogger
}

func New(
	cfg settings.KafkaConsumerSettings,
	kafkaClient *kafkatransport.KafkaClient,
	senderService senderservice.SenderService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) (workers.Worker, error) {
	id := uuid.New().String()

	workerLog := log.With("worker_id", id)

	w := &EmailCodesSenderWorker{
		id:            id,
		log:           workerLog,
		senderService: senderService,
	}

	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(0)

	handlerFunc := func(ctx context.Context, messages []kafkaconsumer.Message) error {
		start := time.Now()

		sentCount := 0
		for _, message := range messages {
			var codeMessageModel *models.CodeMessage
			if err := json.Unmarshal([]byte(message.Body), &codeMessageModel); err != nil {
				w.log.Warnw("email_sender_handler.unmarshal_failed", "err", err, "body", message.Body)
				continue
			}
			if codeMessageModel == nil {
				w.log.Warnw("email_sender_handler.unmarshal_failed", "err", "null message body")
				continue
			}

			if err := w.senderService.SendCodeMessage(ctx, *codeMessageModel, w.id); err != nil {
				continue
			}
			sentCount++
		}

		metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
		metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
		metrics.ProcessedItemsTotal.WithLabelValues(workerType).Add(float64(sentCount))
		return nil
	}

	consumer, err := kafkaconsumer.New(cfg, kafkaClient, workerLog, handlerFunc)
	if err != nil {
		return nil, fmt.Errorf("email sender worker %s: %w", id, err)
	}
	w.consumer = consumer
	return w, nil
}

func (w *EmailCodesSenderWorker) Run(ctx context.Context) {
	w.log.Infow("email_sender.started")
	w.consumer.Run(ctx)
	w.log.Infow("email_sender.stopped")
}
