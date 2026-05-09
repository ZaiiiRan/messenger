package inboxevent

import (
	"encoding/json"
	"time"
)

const (
	MaxAttempts = 3
)

type InboxEvent struct {
	id        string
	payload   json.RawMessage
	status    int16
	attempts  int16
	createdAt time.Time
	updatedAt time.Time
}

func New(
	id string,
	payload json.RawMessage,
) *InboxEvent {
	now := time.Now()
	return &InboxEvent{
		id:        id,
		payload:   payload,
		status:    int16(InboxEventStatusPending),
		attempts:  0,
		createdAt: now,
		updatedAt: now,
	}
}

func FromStorage(
	id string,
	payload json.RawMessage,
	status int16,
	attempts int16,
	createdAt time.Time,
	updatedAt time.Time,
) *InboxEvent {
	return &InboxEvent{
		id:        id,
		payload:   payload,
		status:    status,
		attempts:  attempts,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *InboxEvent) GetID() string               { return e.id }
func (e *InboxEvent) GetPayload() json.RawMessage { return e.payload }
func (e *InboxEvent) GetStatus() InboxEventStatus { return InboxEventStatus(e.status) }
func (e *InboxEvent) GetAttempts() int16          { return e.attempts }
func (e *InboxEvent) GetCreatedAt() time.Time     { return e.createdAt }
func (e *InboxEvent) GetUpdatedAt() time.Time     { return e.updatedAt }

func (e *InboxEvent) IncrementAttempts() error {
	if e.attempts > MaxAttempts {
		return ErrMaxAttemptsReached
	}
	e.attempts++
	return nil
}

func (e *InboxEvent) SetStatus(status InboxEventStatus) error {
	if InboxEventStatus(e.status) == InboxEventStatusSuccess && status != InboxEventStatusSuccess {
		return ErrCannotUpdateStatus
	}
	e.status = int16(status)
	return nil
}

func (e *InboxEvent) ResetAttempts() {
	e.attempts = 0
}

func (e *InboxEvent) SetUpdatedAt(updatedAt *time.Time) {
	if updatedAt == nil {
		e.updatedAt = time.Now()
	} else {
		e.updatedAt = *updatedAt
	}
}
