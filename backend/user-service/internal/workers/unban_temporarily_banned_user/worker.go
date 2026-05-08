package unbantemporarilybanneduserworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	userservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UnbanTemporarilyBannedUserWorker struct {
	workerID    string
	cfg         *settings.UnbanTemporarilyBannedUsersWorkerSettings
	userService userservice.UserService
	log         *zap.SugaredLogger
}

func New(
	cfg settings.UnbanTemporarilyBannedUsersWorkerSettings,
	userService userservice.UserService,
	log *zap.SugaredLogger,
) workers.Worker {
	return &UnbanTemporarilyBannedUserWorker{
		workerID:    uuid.New().String(),
		cfg:         &cfg,
		userService: userService,
		log:         log,
	}
}

func (w *UnbanTemporarilyBannedUserWorker) Run(ctx context.Context) {
	w.log.Infow("unban_temporarily_banned_user_worker.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("unban_temporarily_banned_user_worker.stopped", "worker_id", w.workerID)
			return
		default:
		}

		unbannedCount, err := w.runOnce(ctx)

		var timer *time.Timer
		if unbannedCount != int(w.cfg.BatchSize) && err == nil {
			timer = time.NewTimer(time.Millisecond * time.Duration(w.cfg.NoDataIntervalMS))
		} else {
			timer = time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		}
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("unban_temporarily_banned_user_worker.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *UnbanTemporarilyBannedUserWorker) runOnce(ctx context.Context) (int, error) {
	unbannedCount, err := w.userService.UnbanTemporarilyBannedUsers(ctx, int(w.cfg.BatchSize), w.workerID)
	if err != nil {
		if ctx.Err() != nil {
			return 0, nil
		}
		return 0, err
	}
	return unbannedCount, nil
}
