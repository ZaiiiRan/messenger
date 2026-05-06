package models

import (
	"slices"
	"time"
)

type OutboxEventFilterDal struct {
	Ids      []string `db:"ids" json:"ids"`
	Statuses []int16  `db:"statuses" json:"statuses"`

	AttemptsFrom *int16     `db:"attempts_from" json:"attempts_from"`
	AttemptsTo   *int16     `db:"attempts_to" json:"attempts_to"`
	CreatedFrom  *time.Time `db:"created_from" json:"created_from"`
	CreatedTo    *time.Time `db:"created_to" json:"created_to"`
	UpdatedFrom  *time.Time `db:"updated_from" json:"updated_from"`
	UpdatedTo    *time.Time `db:"updated_to" json:"updated_to"`
}

type QueryOutboxEventsDal struct {
	Filter OutboxEventFilterDal `db:"filter" json:"filter"`

	Limit  int `db:"limit" json:"limit"`
	Offset int `db:"offset" json:"offset"`
}

func NewQueryOutboxEventsDal(
	filter OutboxEventFilterDal,
	page, pageSize int,
) *QueryOutboxEventsDal {
	slices.Sort(filter.Ids)
	slices.Sort(filter.Statuses)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryOutboxEventsDal{
		Filter: filter,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
}
