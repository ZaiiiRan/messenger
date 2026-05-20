package updatingfriendsinprivacysettingsworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	updatefriendsinprivacysettingstasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/update_friends_in_privacy_settings_tasks"
	prommetrics "github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/prom_metrics"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const workerType = "update_friends_in_privacy_settings"

type UpdateFriendsInPrivacySettingsWorker struct {
	workerID                                   string
	cfg                                        *settings.UpdateFriendsInPrivacySettingsWorker
	updateFriendsInPrivacySettingsTasksService updatefriendsinprivacysettingstasksservice.UpdateFriendsInPrivacySettingsTasksService
	log                                        *zap.SugaredLogger
	metrics                                    *prommetrics.WorkerMetrics
}

func New(
	cfg settings.UpdateFriendsInPrivacySettingsWorker,
	updateFriendsInPrivacySettingsTasksService updatefriendsinprivacysettingstasksservice.UpdateFriendsInPrivacySettingsTasksService,
	log *zap.SugaredLogger,
	metrics *prommetrics.WorkerMetrics,
) workers.Worker {
	w := &UpdateFriendsInPrivacySettingsWorker{
		workerID: uuid.New().String(),
		cfg:      &cfg,
		updateFriendsInPrivacySettingsTasksService: updateFriendsInPrivacySettingsTasksService,
		log:     log,
		metrics: metrics,
	}
	metrics.CyclesTotal.WithLabelValues(workerType, "success").Add(0)
	metrics.CyclesTotal.WithLabelValues(workerType, "error").Add(0)
	return w
}

func (w *UpdateFriendsInPrivacySettingsWorker) Run(ctx context.Context) {
	w.log.Infow("updating_friends_in_privacy_settings.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("updating_friends_in_privacy_settings.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("updating_friends_in_privacy_settings.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *UpdateFriendsInPrivacySettingsWorker) runOnce(ctx context.Context) {
	start := time.Now()
	w.updateFriendsInPrivacySettingsTasksService.ProcessUpdateFriendsInPrivacySettingsTasks(ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize), nil)
	w.metrics.CycleDuration.WithLabelValues(workerType).Observe(time.Since(start).Seconds())
	w.metrics.CyclesTotal.WithLabelValues(workerType, "success").Inc()
}
