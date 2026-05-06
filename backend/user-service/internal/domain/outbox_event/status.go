package outboxevent

type OutboxEventStatus int16

const (
	OutboxEventStatusPending OutboxEventStatus = 0
	OutboxEventStatusSuccess OutboxEventStatus = 1
	OutboxEventStatusFailed  OutboxEventStatus = 2
)