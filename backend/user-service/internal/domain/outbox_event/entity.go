package outboxevent

import (
	"encoding/json"
	"time"
)

const (
	MaxAttempts = 3
)

type OutboxEvent struct {
	id        string
	payload   json.RawMessage
	status    int16
	attempts  int16
	createdAt time.Time
	updatedAt time.Time
}

func New(payload json.RawMessage) *OutboxEvent {
	now := time.Now()
	return &OutboxEvent{
		payload:   payload,
		status:    int16(OutboxEventStatusPending),
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
) *OutboxEvent {
	return &OutboxEvent{
		id:        id,
		payload:   payload,
		status:    status,
		attempts:  attempts,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *OutboxEvent) GetID() string                { return e.id }
func (e *OutboxEvent) GetPayload() json.RawMessage  { return e.payload }
func (e *OutboxEvent) GetStatus() OutboxEventStatus { return OutboxEventStatus(e.status) }
func (e *OutboxEvent) GetAttempts() int16           { return e.attempts }
func (e *OutboxEvent) GetCreatedAt() time.Time      { return e.createdAt }
func (e *OutboxEvent) GetUpdatedAt() time.Time      { return e.updatedAt }

func (e *OutboxEvent) SetID(id string) {
	if e.id == "" {
		e.id = id
	}
}

func (e *OutboxEvent) IncrementAttempts() error {
	if e.attempts > MaxAttempts {
		return ErrMaxAttemptsReached
	}
	e.attempts++
	return nil
}

func (e *OutboxEvent) SetPayload(payload json.RawMessage) {
	e.payload = payload
}

func (e *OutboxEvent) ResetAttempts() {
	e.attempts = 0
}

func (e *OutboxEvent) SetStatus(status OutboxEventStatus) error {
	if OutboxEventStatus(e.status) == OutboxEventStatusSuccess && status != OutboxEventStatusSuccess {
		return ErrCannotUpdateStatus
	}
	e.status = int16(status)
	return nil
}

func (e *OutboxEvent) SetUpdatedAt(updatedAt *time.Time) {
	if updatedAt == nil {
		e.updatedAt = time.Now()
	} else {
		e.updatedAt = *updatedAt
	}
}
