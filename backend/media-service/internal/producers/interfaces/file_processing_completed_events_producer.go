package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/producers/models"
)

type FileProcessingCompletedEventsProducer interface {
	ProduceFileProcessingCompletedEvent(ctx context.Context, event *models.FileProcessingCompletedEvent) error
	Close()
}
