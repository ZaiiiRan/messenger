package userrelationshipchangestaskssendingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/config/settings"
	userrelationshipchangestasks "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_relationship_changes_tasks"
	prommetrics "github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "user_relationships_changes_tasks_sending"

type UserRelationshipChangesTasksSendingWorker struct {
	workerID                            string
	cfg                                 *settings.KafkaSendingWorkerSettings
	userRelationshipChangesTasksService userrelationshipchangestasks.UserRelationshipChangesTasksService
	log                                 *zap.SugaredLogger
	metrics                             *prommetrics.WorkerMetrics
}

func New(
	cfg settings.KafkaSendingWorkerSettings,
	userRelationshipChangesTasksService userrelationshipchangestasks.UserRelationshipChangesTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &UserRelationshipChangesTasksSendingWorker{
		workerID:                            uuid.New().String(),
		cfg:                                 &cfg,
		userRelationshipChangesTasksService: userRelationshipChangesTasksService,
		log:                                 log,
		metrics:                             metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *UserRelationshipChangesTasksSendingWorker) Run(ctx context.Context) {
	w.log.Infow("user_relationship_changes_tasks_sending.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("user_relationship_changes_tasks_sending.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("user_relationship_changes_tasks_sending.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *UserRelationshipChangesTasksSendingWorker) runOnce(ctx context.Context) {
	start := time.Now()
	w.userRelationshipChangesTasksService.SendUserRelationshipChangesTasks(
		ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize), nil,
	)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()

}
