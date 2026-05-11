package profile

var (
	ErrProfileIsEmpty         = NewProfileValidationError("domain.profile.validation.error.profile_is_empty")
	ErrInvalidPhoneFormat     = NewProfileValidationError("domain.profile.validation.error.invalid_phone_format")
	ErrInvalidPhoneLength     = NewProfileValidationError("domain.profile.validation.error.invalid_phone_length")
	ErrBirthdateInFuture      = NewProfileValidationError("domain.profile.validation.error.birthdate_in_future")
	ErrInvalidBirthdateFormat = NewProfileValidationError("domain.profile.validation.error.invalid_birthdate_format")
	ErrFirstNameIsEmpty       = NewProfileValidationError("domain.profile.validation.error.firstname_is_empty")
	ErrFirstNameTooShort      = NewProfileValidationError("domain.profile.validation.error.firstname_too_short")
	ErrFirstNameTooLong       = NewProfileValidationError("domain.profile.validation.error.firstname_too_long")
	ErrInvalidFirstNameFormat = NewProfileValidationError("domain.profile.validation.error.invalid_firstname_format")
	ErrLastNameIsEmpty        = NewProfileValidationError("domain.profile.validation.error.lastname_is_empty")
	ErrLastNameTooShort       = NewProfileValidationError("domain.profile.validation.error.lastname_too_short")
	ErrLastNameTooLong        = NewProfileValidationError("domain.profile.validation.error.lastname_too_long")
	ErrInvalidLastNameFormat  = NewProfileValidationError("domain.profile.validation.error.invalid_lastname_format")
	ErrBioTooLong             = NewProfileValidationError("domain.profile.validation.error.bio_too_long")
)

func getNameIsEmptyError(prefix string) *ProfileValidationError {
	switch prefix {
	case "first":
		return ErrFirstNameIsEmpty
	case "last":
		return ErrLastNameIsEmpty
	default:
		return nil
	}
}

func getNameTooShortError(prefix string) *ProfileValidationError {
	switch prefix {
	case "first":
		return ErrFirstNameTooShort
	case "last":
		return ErrLastNameTooShort
	default:
		return nil
	}
}

func getNameTooLongError(prefix string) *ProfileValidationError {
	switch prefix {
	case "first":
		return ErrFirstNameTooLong
	case "last":
		return ErrLastNameTooLong
	default:
		return nil
	}
}

func getInvalidNameFormatError(prefix string) *ProfileValidationError {
	switch prefix {
	case "first":
		return ErrInvalidFirstNameFormat
	case "last":
		return ErrInvalidLastNameFormat
	default:
		return nil
	}
}
