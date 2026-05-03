package expiredtokenclearingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	postgresunitofwork "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ExpiredTokenClearingWorker struct {
	workerID     string
	cfg          *settings.ExpiredTokenClearingWorkerSettings
	tokenService tokenservice.TokenService
	pgClient     *postgres.PostgresClient
	log          *zap.SugaredLogger
}

func New(cfg settings.ExpiredTokenClearingWorkerSettings, tokenService tokenservice.TokenService, pgClient *postgres.PostgresClient, log *zap.SugaredLogger) workers.Worker {
	return &ExpiredTokenClearingWorker{
		cfg:          &cfg,
		tokenService: tokenService,
		pgClient:     pgClient,
		log:          log,
		workerID:     uuid.New().String(),
	}
}

func (w *ExpiredTokenClearingWorker) Run(ctx context.Context) {
	w.log.Infow("expired_token_clearing.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("expired_token_clearing.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			w.log.Infow("expired_token_clearing.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *ExpiredTokenClearingWorker) runOnce(ctx context.Context) {
	uow := postgresunitofwork.New(w.pgClient)
	defer uow.Close()

	if err := w.tokenService.DeleteExpiredTokens(ctx, uow, w.cfg.BatchSize, w.workerID); err != nil {
		if ctx.Err() != nil {
			return
		}
	}
}
