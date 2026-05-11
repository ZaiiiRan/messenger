package status

var (
	ErrBannedUntilInPast              = NewStatusValidationError("domain.status.validation.error.banned_until_in_past")
	ErrPermanentlyDeletedIfNotDeleted = NewStatusValidationError("domain.status.validation.error.permanently_delete_if_not_deleted")
)
