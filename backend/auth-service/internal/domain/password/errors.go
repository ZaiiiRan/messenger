package password

var (
	ErrPasswordTooShort      = NewPasswordValidationError("domain.password.validation.error.password_must_be_at_least_8_characters_long")
	ErrPasswordTooLong       = NewPasswordValidationError("domain.password.validation.error.password_must_be_at_most_72_characters_long")
	ErrPasswordNoUppercase   = NewPasswordValidationError("domain.password.validation.error.password_must_contain_at_least_one_uppercase_letter")
	ErrPasswordNoLowercase   = NewPasswordValidationError("domain.password.validation.error.password_must_contain_at_least_one_lowercase_letter")
	ErrPasswordNoDigit       = NewPasswordValidationError("domain.password.validation.error.password_must_contain_at_least_one_digit")
	ErrPasswordNoSpecial     = NewPasswordValidationError("domain.password.validation.error.password_must_contain_at_least_one_special_character")
	ErrOldAndNewPasswordSame = NewPasswordValidationError("domain.password.validation.error.old_and_new_passwords_are_the_same")
	ErrOldPasswordIncorrect  = NewPasswordValidationError("domain.password.validation.error.old_password_is_incorrect")
	ErrPasswordChangeTooSoon = NewPasswordValidationError("domain.password.validation.error.password_can_be_changed_only_once_per_24_hours")
)
