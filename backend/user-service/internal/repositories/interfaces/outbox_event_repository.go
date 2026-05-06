package interfaces

import (
	"context"

	outboxevent "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/outbox_event"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
)

type OutboxEventRepository interface {
	Create(ctx context.Context, events []*outboxevent.OutboxEvent) error
	Update(ctx context.Context, events []*outboxevent.OutboxEvent) error
	Delete(ctx context.Context, events []*outboxevent.OutboxEvent) error
	Query(ctx context.Context, query *models.QueryOutboxEventsDal) ([]*outboxevent.OutboxEvent, error)
	QueryLocked(ctx context.Context, query *models.QueryOutboxEventsLockedDal) ([]*outboxevent.OutboxEvent, error)
}
