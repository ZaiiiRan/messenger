package event

type EventError struct {
	message string
}

func NewEventError(message string) *EventError {
	return &EventError{message: message}
}

func (e *EventError) Error() string {
	return e.message
}
