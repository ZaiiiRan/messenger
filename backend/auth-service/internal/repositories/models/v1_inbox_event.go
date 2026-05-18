package models

import (
	"encoding/json"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/event"
)

type V1InboxEventDal struct {
	Id        string          `db:"id" json:"id"`
	Payload   json.RawMessage `db:"payload" json:"payload"`
	Status    int16           `db:"status" json:"status"`
	Attempts  int16           `db:"attempts" json:"attempts"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

func V1InboxEventFromDomain(event *event.Event) V1InboxEventDal {
	if event == nil {
		return V1InboxEventDal{}
	}
	return V1InboxEventDal{
		Id:        event.GetID(),
		Payload:   event.GetPayload(),
		Status:    int16(event.GetStatus()),
		Attempts:  event.GetAttempts(),
		CreatedAt: event.GetCreatedAt(),
		UpdatedAt: event.GetUpdatedAt(),
	}
}

func (e V1InboxEventDal) IsNull() bool { return false }
func (e V1InboxEventDal) Index(i int) any {
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

func (e V1InboxEventDal) ToDomain() *event.Event {
	return event.FromStorage(
		e.Id,
		e.Payload,
		e.Status,
		e.Attempts,
		e.CreatedAt,
		e.UpdatedAt,
	)
}
