package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/event"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type EventRepository interface {
	CreateEvents(ctx context.Context, events []*event.Event) error
	UpdateEvents(ctx context.Context, events []*event.Event) error
	DeleteEvents(ctx context.Context, events []*event.Event) error
	QueryEvents(ctx context.Context, query *models.QueryEventsDal) ([]*event.Event, error)
	QueryEventsLocked(ctx context.Context, query *models.QueryEventsLockedDal) ([]*event.Event, error)
}
