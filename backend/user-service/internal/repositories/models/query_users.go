package models

import (
	"slices"
	"time"
)

type UserFilterDal struct {
	Ids              []string `db:"ids" json:"ids"`
	Usernames        []string `db:"usernames" json:"usernames"`
	PartialUsernames []string `db:"partial_usernames" json:"partial_usernames"`
	Emails           []string `db:"emails" json:"emails"`
	PartialEmails    []string `db:"partial_emails" json:"partial_emails"`

	PhoneNumbers []string `db:"phone_numbers" json:"phone_numbers"`
	PartialNames []string `db:"partial_names" json:"partial_names"`

	IsConfirmed         *bool `db:"is_confirmed" json:"is_confirmed"`
	IsDeleted           *bool `db:"is_deleted" json:"is_deleted"`
	IsPermanentlyBanned *bool `db:"is_permanently_banned" json:"is_permanently_banned"`
	IsTemporarilyBanned *bool `db:"is_temporarily_banned" json:"is_temporarily_banned"`

	DeletedFrom *time.Time `db:"deleted_from" json:"deleted_from"`
	DeletedTo   *time.Time `db:"deleted_to" json:"deleted_to"`
	CreatedFrom *time.Time `db:"created_from" json:"created_from"`
	CreatedTo   *time.Time `db:"created_to" json:"created_to"`
	UpdatedFrom *time.Time `db:"updated_from" json:"updated_from"`
	UpdatedTo   *time.Time `db:"updated_to" json:"updated_to"`
}

type QueryUsersDal struct {
	Filter UserFilterDal `db:"filter" json:"filter"`

	Limit  int `db:"limit" json:"limit"`
	Offset int `db:"offset" json:"offset"`
}

func NewQueryUsersDal(
	filter UserFilterDal,
	page, pageSize int,
) *QueryUsersDal {

	slices.Sort(filter.Ids)
	slices.Sort(filter.Usernames)
	slices.Sort(filter.PartialUsernames)
	slices.Sort(filter.Emails)
	slices.Sort(filter.PartialEmails)
	slices.Sort(filter.PhoneNumbers)
	slices.Sort(filter.PartialNames)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryUsersDal{
		Filter: filter,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
}
