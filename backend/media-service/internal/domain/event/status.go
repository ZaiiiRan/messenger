package event

type EventStatus int16

const (
	EventStatusPending EventStatus = 0
	EventStatusSuccess EventStatus = 1
	EventStatusFailed  EventStatus = 2
)
