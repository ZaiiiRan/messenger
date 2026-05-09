package inboxevent

type InboxEventStatus int16

const (
	InboxEventStatusPending InboxEventStatus = 0
	InboxEventStatusSuccess InboxEventStatus = 1
	InboxEventStatusFailed  InboxEventStatus = 2
)
