package models

import (
	"time"
)

type QueryEventsLockedDal struct {
	RetryAfter time.Time `db:"retry_after" json:"retry_after"`

	Limit int `db:"limit" json:"limit"`
}

func NewQueryEventsLockedDal(
	retryAfter time.Time,
	pageSize int,
) *QueryEventsLockedDal {
	if pageSize <= 0 {
		pageSize = 50
	}

	return &QueryEventsLockedDal{
		RetryAfter: retryAfter,
		Limit:   pageSize,
	}
}
