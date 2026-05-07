package implkafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/producers/models"
)

func (p *Producer) ProduceUserDataDeletionTask(ctx context.Context, task *models.UserDataDeletionTask) error {
	value, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("ProduceUserDataDeletionTask: marshal: %w", err)
	}

	return p.Produce(ctx, Message{
		Key:   task.UserId,
		Value: string(value),
	})
}
