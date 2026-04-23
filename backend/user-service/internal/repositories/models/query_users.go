package models

import (
	"slices"
	"time"
)

type QueryUsersDal struct {
	Ids              []string `db:"ids"`
	Usernames        []string `db:"usernames"`
	PartialUsernames []string `db:"partial_usernames"`
	Emails           []string `db:"emails"`
	PartialEmails    []string `db:"partial_emails"`

	PhoneNumbers []string `db:"phone_numbers"`
	PartialNames []string `db:"partial_names"`

	IsConfirmed         *bool      `db:"is_confirmed"`
	IsDeleted           *bool      `db:"is_deleted"`
	IsPermanentlyBanned *bool      `db:"is_permanently_banned"`
	BannedUntil         *time.Time `db:"banned_until"`

	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

func NewQueryUsersDal(
	ids []string,
	usernames []string,
	partialUsernames []string,
	emails []string,
	partialEmails []string,
	phoneNumbers []string,
	partialNames []string,
	isConfirmed *bool,
	isDeleted *bool,
	isPermanentlyBanned *bool,
	bannedUntil *time.Time,
	page, pageSize int,
) *QueryUsersDal {
	slices.Sort(ids)
	slices.Sort(usernames)
	slices.Sort(partialUsernames)
	slices.Sort(emails)
	slices.Sort(partialEmails)
	slices.Sort(phoneNumbers)
	slices.Sort(partialNames)

	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}

	return &QueryUsersDal{
		Ids:                 ids,
		Usernames:           usernames,
		PartialUsernames:    partialUsernames,
		Emails:              emails,
		PartialEmails:       partialEmails,
		PhoneNumbers:        phoneNumbers,
		PartialNames:        partialNames,
		IsConfirmed:         isConfirmed,
		IsDeleted:           isDeleted,
		IsPermanentlyBanned: isPermanentlyBanned,
		BannedUntil:         bannedUntil,
		Limit:               pageSize,
		Offset:              (page - 1) * pageSize,
	}
}
