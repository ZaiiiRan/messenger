package interfaces

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/event"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type InboxEventRepository interface {
	CreateInboxEvents(ctx context.Context, events []*event.Event) error
	UpdateInboxEvents(ctx context.Context, events []*event.Event) error
	DeleteInboxEvents(ctx context.Context, events []*event.Event) error
	QueryInboxEvents(ctx context.Context, query *models.QueryEventsDal) ([]*event.Event, error)
	QueryInboxEventsLocked(ctx context.Context, query *models.QueryEventsLockedDal) ([]*event.Event, error)
}
