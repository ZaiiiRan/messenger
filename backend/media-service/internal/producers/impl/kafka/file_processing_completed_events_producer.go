package implkafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/models"
)

func (p *Producer) ProduceFileProcessingCompletedEvent(ctx context.Context, event *models.FileProcessingCompletedEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ProduceFileProcessingCompletedEvent: marshal: %w", err)
	}

	return p.Produce(ctx, Message{
		Key:   event.FileId,
		Value: string(value),
	})
}
