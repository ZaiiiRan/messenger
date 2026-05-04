package status

var (
	ErrBannedUntilInPast = NewStatusValidationError("domain.status.validation.error.banned_until_in_past")
)
