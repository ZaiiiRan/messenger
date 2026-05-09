package interfaces

import (
	"context"

	inboxevent "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/inbox_event"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
)

type InboxEventRepository interface {
	CreateInboxEvents(ctx context.Context, events []*inboxevent.InboxEvent) error
	UpdateInboxEvents(ctx context.Context, events []*inboxevent.InboxEvent) error
	DeleteInboxEvents(ctx context.Context, events []*inboxevent.InboxEvent) error
	QueryInboxEvents(ctx context.Context, query *models.QueryInboxEventsDal) ([]*inboxevent.InboxEvent, error)
	QueryInboxEventsLocked(ctx context.Context, query *models.QueryInboxEventsLockedDal) ([]*inboxevent.InboxEvent, error)
}
