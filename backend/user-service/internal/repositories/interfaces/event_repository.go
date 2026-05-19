package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
)

type EventRepository interface {
	Create(ctx context.Context, events []*event.Event) error
	Update(ctx context.Context, events []*event.Event) error
	Delete(ctx context.Context, events []*event.Event) error
	Query(ctx context.Context, query *models.QueryEventsDal) ([]*event.Event, error)
	QueryLocked(ctx context.Context, query *models.QueryEventsLockedDal) ([]*event.Event, error)
}
