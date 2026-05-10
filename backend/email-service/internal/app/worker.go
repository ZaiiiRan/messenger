package app

import (
	"context"
	"sync"
	"time"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config"
	senderservice "github.com/ZaiiiRan/messenger/backend/email-service/internal/services/sender"
	appi18n "github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/i18n"
	kafkatransport "github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/kafka"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/smtp"
	emailcodessenderworker "github.com/ZaiiiRan/messenger/backend/email-service/internal/workers/email_codes_sender"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/logger"
	"go.uber.org/zap"
)

type WorkerApp struct {
	cfg *config.WorkerConfig
	log *zap.SugaredLogger

	emailCodesKafkaClient *kafkatransport.KafkaClient
	smtpClient            *smtp.SMTPClient

	senderService senderservice.SenderService

	workersCtx    context.Context
	workersCancel context.CancelFunc
	workersWG     sync.WaitGroup
}

func NewWorkerApp() (*WorkerApp, error) {
	cfg, err := config.LoadWorkerConfig()
	if err != nil {
		return nil, err
	}

	log, err := logger.New()
	if err != nil {
		return nil, err
	}

	return &WorkerApp{
		cfg: cfg,
		log: log,
	}, nil
}

func (a *WorkerApp) Run(ctx context.Context) error {
	appi18n.Init()

	if err := a.initEmailCodesKafkaClient(); err != nil {
		return err
	}
	a.initSMTPClient()
	a.initSenderService()

	a.workersCtx, a.workersCancel = context.WithCancel(ctx)

	if err := a.startEmailCodesSenderWorkers(); err != nil {
		return err
	}

	a.log.Infow("app.started")
	return nil
}

func (a *WorkerApp) Stop(ctx context.Context) {
	a.log.Infow("app.stopping")

	shCtx, cancel := context.WithTimeout(ctx, time.Duration(a.cfg.Shutdown.ShutdownTimeout)*time.Second)
	defer cancel()

	if a.workersCancel != nil {
		a.workersCancel()
	}

	workersStopped := make(chan struct{})
	go func() {
		a.workersWG.Wait()
		close(workersStopped)
	}()

	select {
	case <-workersStopped:
	case <-shCtx.Done():
		a.log.Warnw("app.workers_shutdown_timeout")
	}

	a.emailCodesKafkaClient.Close()
	a.smtpClient.Close()

	a.log.Infow("app.stopped")
}

func (a *WorkerApp) initEmailCodesKafkaClient() error {
	kafkaClient, err := kafkatransport.New(a.cfg.EmailCodesSenderWorker.KafkaConsumerSettings.KafkaSettings)
	if err != nil {
		a.log.Errorw("app.email_codes_kafka_connect_failed", "err", err)
		return err
	}
	a.emailCodesKafkaClient = kafkaClient
	a.log.Infow("app.email_codes_kafka_connected")
	return nil
}

func (a *WorkerApp) initSMTPClient() {
	a.smtpClient = smtp.New(a.cfg.SMTPClient)
}

func (a *WorkerApp) initSenderService() {
	a.senderService = senderservice.NewSenderService(a.cfg.HTMLGenerator, a.smtpClient, a.log)
}

func (a *WorkerApp) startEmailCodesSenderWorkers() error {
	for i := 0; i < int(a.cfg.EmailCodesSenderWorker.Count); i++ {
		w, err := emailcodessenderworker.New(a.cfg.EmailCodesSenderWorker.KafkaConsumerSettings, a.emailCodesKafkaClient, a.senderService, a.log)
		if err != nil {
			a.log.Errorw("app.email_sender_worker_init_failed", "err", err, "worker_id", i)
			return err
		}
		a.workersWG.Add(1)
		go func() {
			defer a.workersWG.Done()
			w.Run(a.workersCtx)
		}()
	}
	return nil
}
