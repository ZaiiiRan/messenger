package privacysettings

type PrivacySettingsError struct {
	message string
}

func NewPrivacySettingsError(message string) *PrivacySettingsError {
	return &PrivacySettingsError{
		message: message,
	}
}

func (e *PrivacySettingsError) Error() string {
	return e.message
}
