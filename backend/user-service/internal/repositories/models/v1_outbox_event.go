package models

import (
	"encoding/json"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/event"
)

type V1OutboxEventDal struct {
	Id        string          `db:"id" json:"id"`
	Payload   json.RawMessage `db:"payload" json:"payload"`
	Status    int16           `db:"status" json:"status"`
	Attempts  int16           `db:"attempts" json:"attempts"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

func V1OutboxEventFromDomain(event *event.Event) V1OutboxEventDal {
	if event == nil {
		return V1OutboxEventDal{}
	}
	return V1OutboxEventDal{
		Id:        event.GetID(),
		Payload:   event.GetPayload(),
		Status:    int16(event.GetStatus()),
		Attempts:  event.GetAttempts(),
		CreatedAt: event.GetCreatedAt(),
		UpdatedAt: event.GetUpdatedAt(),
	}
}

func (e V1OutboxEventDal) IsNull() bool { return false }
func (e V1OutboxEventDal) Index(i int) any {
	switch i {
	case 0:
		return e.Id
	case 1:
		return e.Payload
	case 2:
		return e.Status
	case 3:
		return e.Attempts
	case 4:
		return e.CreatedAt
	case 5:
		return e.UpdatedAt
	default:
		return nil
	}
}

func (e V1OutboxEventDal) ToDomain() *event.Event {
	return event.FromStorage(
		e.Id,
		e.Payload,
		e.Status,
		e.Attempts,
		e.CreatedAt,
		e.UpdatedAt,
	)
}
