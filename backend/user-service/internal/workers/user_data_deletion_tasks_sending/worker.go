package usersdatadeletiontaskssendingworker

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	userdatadeletiontasksservice "github.com/ZaiiiRan/messenger/backend/user-service/internal/services/user_data_deletion_tasks"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/workers"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserDataDeletionTasksSendingWorker struct {
	workerID                     string
	cfg                          *settings.KafkaSendingWorkerSettings
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService
	log                          *zap.SugaredLogger
}

func New(
	cfg settings.KafkaSendingWorkerSettings,
	userDataDeletionTasksService userdatadeletiontasksservice.UserDataDeletionTasksService,
	log *zap.SugaredLogger,
) workers.Worker {
	return &UserDataDeletionTasksSendingWorker{
		workerID:                     uuid.New().String(),
		cfg:                          &cfg,
		userDataDeletionTasksService: userDataDeletionTasksService,
		log:                          log,
	}
}

func (w *UserDataDeletionTasksSendingWorker) Run(ctx context.Context) {
	w.log.Infow("user_data_deletion_tasks_sending.started", "worker_id", w.workerID)
	for {
		select {
		case <-ctx.Done():
			w.log.Infow("user_data_deletion_tasks_sending.stopped", "worker_id", w.workerID)
			return
		default:
		}

		w.runOnce(ctx)

		timer := time.NewTimer(time.Millisecond * time.Duration(w.cfg.IntervalMS))
		select {
		case <-ctx.Done():
			<-timer.C
			w.log.Infow("user_data_deletion_tasks_sending.stopped", "worker_id", w.workerID)
			return
		case <-timer.C:
		}
	}
}

func (w *UserDataDeletionTasksSendingWorker) runOnce(ctx context.Context) {
	w.userDataDeletionTasksService.SendUserDataDeletionTasks(
		ctx, w.workerID, w.cfg.RetryIntervalMS, int(w.cfg.BatchSize),
	)
}
