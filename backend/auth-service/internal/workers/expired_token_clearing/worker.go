package expiredtokenclearingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	prommetrics "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "expired_token_clearing"

type ExpiredTokenClearingWorker struct {
	workerID     string
	cfg          *settings.ExpiredTokenClearingWorkerSettings
	tokenService tokenservice.TokenService
	log          *zap.SugaredLogger
	metrics      *prommetrics.WorkerMetrics
}

func New(cfg settings.ExpiredTokenClearingWorkerSettings, tokenService tokenservice.TokenService, log *zap.SugaredLogger, metrics *prommetrics.WorkerMetrics) workers.Worker {
	w := &ExpiredTokenClearingWorker{
		cfg:          &cfg,
		tokenService: tokenService,
		log:          log,
		workerID:     uuid.New().String(),
		metrics:      metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
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
	start := time.Now()
	w.tokenService.DeleteExpiredTokens(ctx, w.workerID, w.cfg.BatchSize)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
