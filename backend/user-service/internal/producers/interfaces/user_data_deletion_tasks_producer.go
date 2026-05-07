package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/models"
)

type UserDataDeletionTasksProducer interface {
	ProduceUserDataDeletionTask(ctx context.Context, task *models.UserDataDeletionTask) error
}
