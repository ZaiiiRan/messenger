package models

import (
	"time"
)

type QueryOutboxEventsLockedDal struct {
	RetryAfter time.Time `db:"retry_after" json:"retry_after"`

	Limit int `db:"limit" json:"limit"`
}

func NewQueryOutboxEventsLockedDal(
	retryAfter time.Time,
	pageSize int,
) *QueryOutboxEventsLockedDal {
	if pageSize <= 0 {
		pageSize = 50
	}

	return &QueryOutboxEventsLockedDal{
		RetryAfter: retryAfter,
		Limit:   pageSize,
	}
}
