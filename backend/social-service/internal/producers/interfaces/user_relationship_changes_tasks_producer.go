package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/producers/models"
)

type UserRelationshipChangesTasksProducer interface {
	ProduceUserRelationshipChangesTask(ctx context.Context, task *models.UserRelationshipChangeTask) error
}
