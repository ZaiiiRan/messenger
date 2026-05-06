package outboxevent

var (
	ErrMaxAttemptsReached = NewOutboxEventError("domain.outbox_event.error.max_attempts_reached")
	ErrCannotUpdateStatus = NewOutboxEventError("domain.outbox_event.error.cannot_update_status")
)
