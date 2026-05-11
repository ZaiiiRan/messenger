package inboxevent

var (
	ErrMaxAttemptsReached = NewInboxEventError("domain.outbox_inbox.error.max_attempts_reached")
	ErrCannotUpdateStatus = NewInboxEventError("domain.outbox_inbox.error.cannot_update_status")
)
