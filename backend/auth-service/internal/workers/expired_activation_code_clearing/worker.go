package expiredactivationcodeclearingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	codeservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/code"
	prommetrics "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "expired_activation_code_clearing"

type ExpiredActivationCodeClearingWorker struct {
	workerID    string
	cfg         *settings.ExpiredCodesClearingWorkerSettings
	codeService codeservice.CodeService
	log         *zap.SugaredLogger
	metrics     *prommetrics.WorkerMetrics
}

func New(cfg settings.ExpiredCodesClearingWorkerSettings, codeService codeservice.CodeService, log *zap.SugaredLogger, metrics *prommetrics.WorkerMetrics) workers.Worker {
	w := &ExpiredActivationCodeClearingWorker{
		cfg:         &cfg,
		codeService: codeService,
		log:         log,
		workerID:    uuid.New().String(),
		metrics:     metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *ExpiredActivationCodeClearingWorker) Run(ctx context.Context) {
	w.log.Infow("expired_activation_code_clearing.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("expired_activation_code_clearing.stopped", "worker_id", w.workerID)
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
			w.log.Infow("expired_activation_code_clearing.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *ExpiredActivationCodeClearingWorker) runOnce(ctx context.Context) {
	start := time.Now()
	w.codeService.DeleteExpiredCodes(ctx, w.workerID, w.cfg.BatchSize, code.CodeTypeActivation)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
