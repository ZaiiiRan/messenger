package event

import (
	"encoding/json"
	"time"
)

const (
	MaxAttempts = 3
)

type Event struct {
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
) *Event {
	now := time.Now()
	return &Event{
		id:        id,
		payload:   payload,
		status:    int16(EventStatusPending),
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
) *Event {
	return &Event{
		id:        id,
		payload:   payload,
		status:    status,
		attempts:  attempts,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *Event) GetID() string               { return e.id }
func (e *Event) GetPayload() json.RawMessage { return e.payload }
func (e *Event) GetStatus() EventStatus      { return EventStatus(e.status) }
func (e *Event) GetAttempts() int16          { return e.attempts }
func (e *Event) GetCreatedAt() time.Time     { return e.createdAt }
func (e *Event) GetUpdatedAt() time.Time     { return e.updatedAt }

func (e *Event) SetID(id string) {
	if e.id == "" {
		e.id = id
	}
}

func (e *Event) IncrementAttempts() error {
	if e.attempts > MaxAttempts {
		return ErrMaxAttemptsReached
	}
	e.attempts++
	return nil
}

func (e *Event) SetPayload(payload json.RawMessage) {
	e.payload = payload
}

func (e *Event) ResetAttempts() {
	e.attempts = 0
}

func (e *Event) SetStatus(status EventStatus) error {
	if EventStatus(e.status) == EventStatusSuccess && status != EventStatusSuccess {
		return ErrCannotUpdateStatus
	}
	e.status = int16(status)
	return nil
}

func (e *Event) SetUpdatedAt(updatedAt *time.Time) {
	if updatedAt == nil {
		e.updatedAt = time.Now()
	} else {
		e.updatedAt = *updatedAt
	}
}
