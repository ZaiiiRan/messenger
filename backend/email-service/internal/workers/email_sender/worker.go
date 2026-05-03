package emailsenderworker

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/workers"
	"go.uber.org/zap"
)

type EmailSenderWorker struct {
	id       int
	consumer *kafkatransport.Consumer
	log      *zap.SugaredLogger
}

func New(id int, cfg settings.KafkaConsumerSettings, kafkaClient *kafkatransport.KafkaClient, log *zap.SugaredLogger) (workers.Worker, error) {
	workerLog := log.With("worker_id", id)
	consumer, err := kafkatransport.NewConsumer(cfg, kafkaClient, workerLog)
	if err != nil {
		return nil, fmt.Errorf("email sender worker %d: %w", id, err)
	}
	return &EmailSenderWorker{
		id:       id,
		consumer: consumer,
		log:      workerLog,
	}, nil
}

func (w *EmailSenderWorker) Run(ctx context.Context) {
	w.log.Infow("email_sender.started")
	w.consumer.Run(ctx)
	w.log.Infow("email_sender.stopped")
}
