package models

import "time"

type FileFilterDal struct {
	Ids      []string `json:"ids" db:"ids"`
	OwnerIds []string `json:"owner_ids" db:"owner_ids"`
	Statuses []int16  `json:"statuses" db:"statuses"`

	CreatedFrom *time.Time `db:"created_from" json:"created_from"`
	CreatedTo   *time.Time `db:"created_to" json:"created_to"`
	UpdatedFrom *time.Time `db:"updated_from" json:"updated_from"`
	UpdatedTo   *time.Time `db:"updated_to" json:"updated_to"`
}

type QueryFilesDal struct {
	Filter FileFilterDal `json:"filter" db:"filter"`

	Limit  int `json:"limit" db:"limit"`
	Offset int `json:"offset" db:"offset"`
}

func NewQueryFilesDal(
	filter FileFilterDal,
	page, pageSize int,
) *QueryFilesDal {
	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryFilesDal{
		Filter: filter,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
}
