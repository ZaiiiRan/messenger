package event

var (
	ErrMaxAttemptsReached = NewEventError("domain.event.error.max_attempts_reached")
	ErrCannotUpdateStatus = NewEventError("domain.event.error.cannot_update_status")
)
