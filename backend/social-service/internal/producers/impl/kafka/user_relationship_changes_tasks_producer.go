package implkafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/producers/models"
)

func (p *Producer) ProduceUserRelationshipChangesTask(ctx context.Context, task *models.UserRelationshipChangeTask) error {
	value, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("ProduceUserRelationshipChangesTask: marshal: %w", err)
	}

	return p.Produce(ctx, Message{
		Key:   task.Id,
		Value: string(value),
	})
}
