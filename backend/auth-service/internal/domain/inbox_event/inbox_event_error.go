package inboxevent

type InboxEventError struct {
	message string
}

func NewInboxEventError(message string) *InboxEventError {
	return &InboxEventError{
		message: message,
	}
}

func (e *InboxEventError) Error() string {
	return e.message
}
