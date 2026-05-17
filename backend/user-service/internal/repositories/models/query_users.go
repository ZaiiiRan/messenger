package models

import (
	"slices"
	"time"
)

type UserFilterDal struct {
	Ids              []string `db:"ids" json:"ids"`
	ExcludeIds       []string `db:"exclude_ids" json:"exclude_ids"`
	Usernames        []string `db:"usernames" json:"usernames"`
	PartialUsernames []string `db:"partial_usernames" json:"partial_usernames"`
	Emails           []string `db:"emails" json:"emails"`
	PartialEmails    []string `db:"partial_emails" json:"partial_emails"`

	PhoneNumbers []string `db:"phone_numbers" json:"phone_numbers"`
	PartialNames []string `db:"partial_names" json:"partial_names"`

	SearchFilter *string `db:"search_filter" json:"search_filter"`

	IsConfirmed          *bool `db:"is_confirmed" json:"is_confirmed"`
	IsDeleted            *bool `db:"is_deleted" json:"is_deleted"`
	IsPermanentlyBanned  *bool `db:"is_permanently_banned" json:"is_permanently_banned"`
	IsTemporarilyBanned  *bool `db:"is_temporarily_banned" json:"is_temporarily_banned"`
	IsPermanentlyDeleted *bool `db:"is_permanently_deleted" json:"is_permanently_deleted"`

	DeletedFrom      *time.Time `db:"deleted_from" json:"deleted_from"`
	DeletedTo        *time.Time `db:"deleted_to" json:"deleted_to"`
	BannedUntilFrom  *time.Time `db:"banned_until_from" json:"banned_until_from"`
	BannedUntilTo    *time.Time `db:"banned_until_to" json:"banned_until_to"`
	CreatedFrom      *time.Time `db:"created_from" json:"created_from"`
	CreatedTo        *time.Time `db:"created_to" json:"created_to"`
	UpdatedFrom      *time.Time `db:"updated_from" json:"updated_from"`
	UpdatedTo        *time.Time `db:"updated_to" json:"updated_to"`
	EmailUpdatedFrom *time.Time `db:"email_updated_from" json:"email_updated_from"`
	EmailUpdatedTo   *time.Time `db:"email_updated_to" json:"email_updated_to"`
}

type UserSortDal struct {
	SortByUsername   bool `db:"sort_by_username" json:"sort_by_username"`
	SortInactiveLast bool `db:"sort_inactive_last" json:"sort_inactive_last"`
	PreserveIDsOrder bool `db:"preserve_ids_order" json:"preserve_ids_order"`
}

type QueryUsersDal struct {
	Filter UserFilterDal `db:"filter" json:"filter"`
	Sort   UserSortDal   `db:"sort" json:"sort"`

	Limit  int `db:"limit" json:"limit"`
	Offset int `db:"offset" json:"offset"`
}

func NewQueryUsersDal(
	filter UserFilterDal,
	sort UserSortDal,
	page, pageSize int,
) *QueryUsersDal {

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
		Sort:   sort,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
}
