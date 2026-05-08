package deleteduserdataclearingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DeletedUserDataClearingWorker struct {
	workerID    string
	cfg         *settings.DeletedUsersDataClearingWorkerSettings
	userService userservice.UserService
	log         *zap.SugaredLogger
}

func New(
	cfg settings.DeletedUsersDataClearingWorkerSettings,
	userService userservice.UserService,
	log *zap.SugaredLogger,
) workers.Worker {
	return &DeletedUserDataClearingWorker{
		workerID:    uuid.New().String(),
		cfg:         &cfg,
		userService: userService,
		log:         log,
	}
}

func (w *DeletedUserDataClearingWorker) Run(ctx context.Context) {
	w.log.Infow("deleted_user_data_clearing.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("deleted_user_data_clearing.stopped", "worker_id", w.workerID)
			return
		default:
		}

		deletedCount, err := w.runOnce(ctx)

		var timer *time.Timer
		if deletedCount != int(w.cfg.BatchSize) && err == nil {
			timer = time.NewTimer(time.Millisecond * time.Duration(w.cfg.NoDataIntervalMS))
		} else {
			timer = time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		}
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("deleted_user_data_clearing.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *DeletedUserDataClearingWorker) runOnce(ctx context.Context) (int, error) {
	deletedCount, err := w.userService.ClearDeletedUsers(ctx, int(w.cfg.BatchSize), w.workerID)
	if err != nil {
		if ctx.Err() != nil {
			return 0, nil
		}
		return 0, err
	}
	return deletedCount, nil
}
