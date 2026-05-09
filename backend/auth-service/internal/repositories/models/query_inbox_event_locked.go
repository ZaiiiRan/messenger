package models

import "time"

type QueryInboxEventsLockedDal struct {
	RetryAfter   time.Time  `db:"retry_after" json:"retry_after"`
	CreatedAfter *time.Time `db:"created_after" json:"created_after"`

	Limit int `db:"limit" json:"limit"`
}

func NewQueryInboxEventsLockedDal(
	retryAfter time.Time,
	createdAfter *time.Time,
	pageSize int,
) *QueryInboxEventsLockedDal {
	if pageSize <= 0 {
		pageSize = 50
	}

	return &QueryInboxEventsLockedDal{
		RetryAfter:   retryAfter,
		CreatedAfter: createdAfter,
		Limit:        pageSize,
	}
}
