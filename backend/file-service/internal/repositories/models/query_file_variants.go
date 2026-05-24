package models

import (
	"time"
)

type FileVariantFilterDal struct {
	Ids         []string   `db:"ids" json:"ids"`
	FileIds     []string   `db:"file_ids" json:"file_ids"`
	Variants    []int16    `db:"variants" json:"variants"`
	S3Keys      []string   `db:"s3_keys" json:"s3_keys"`
	CreatedFrom *time.Time `db:"created_from" json:"created_from"`
	CreatedTo   *time.Time `db:"created_to" json:"created_to"`
	UpdatedFrom *time.Time `db:"updated_from" json:"updated_from"`
	UpdatedTo   *time.Time `db:"updated_to" json:"updated_to"`
}

type QueryFileVariantsDal struct {
	Filter FileVariantFilterDal `db:"filter" json:"filter"`

	Limit  *int `db:"limit" json:"limit"`
	Offset *int `db:"offset" json:"offset"`
}

func NewQueryFileVariantsDal(
	filter FileVariantFilterDal,
	page, pageSize *int,
) *QueryFileVariantsDal {
	var limit, offset *int

	if pageSize != nil {
		if *pageSize > 0 {
			limit = pageSize
		} else {
			lim := 50
			limit = &lim
		}
	}

	if page != nil && limit != nil {
		if *page > 1 {
			off := (*page - 1) * *limit
			offset = &off
		} else {
			off := 0
			offset = &off
		}
	}

	return &QueryFileVariantsDal{
		Filter: filter,
		Limit:  limit,
		Offset: offset,
	}
}
