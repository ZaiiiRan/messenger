package emailsenderworker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	kafkaconsumer "github.com/ZaiiiRan/messenger/backend/email-service/internal/consumers/impl/kafka"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/consumers/models"
	senderservice "github.com/ZaiiiRan/messenger/backend/email-service/internal/services/sender"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EmailSenderWorker struct {
	id            string
	consumer      *kafkaconsumer.Consumer
	senderService senderservice.SenderService
	log           *zap.SugaredLogger
}

func New(cfg settings.KafkaConsumerSettings, kafkaClient *kafkatransport.KafkaClient, senderService senderservice.SenderService, log *zap.SugaredLogger) (workers.Worker, error) {
	id := uuid.New().String()

	workerLog := log.With("worker_id", id)

	w := &EmailSenderWorker{
		id:            id,
		log:           workerLog,
		senderService: senderService,
	}

	handlerFunc := func(ctx context.Context, messages []kafkaconsumer.Message) error {
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
		}
		return nil
	}

	consumer, err := kafkaconsumer.NewConsumer(cfg, kafkaClient, workerLog, handlerFunc)
	if err != nil {
		return nil, fmt.Errorf("email sender worker %s: %w", id, err)
	}
	w.consumer = consumer
	return w, nil
}

func (w *EmailSenderWorker) Run(ctx context.Context) {
	w.log.Infow("email_sender.started")
	w.consumer.Run(ctx)
	w.log.Infow("email_sender.stopped")
}
