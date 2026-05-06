package outboxevent

type OutboxEventError struct {
	message string
}

func NewOutboxEventError(message string) *OutboxEventError {
	return &OutboxEventError{message: message}
}

func (e *OutboxEventError) Error() string {
	return e.message
}
