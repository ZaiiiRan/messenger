package updatefriendsinprivacysettingstasksservice

type UpdateFriendsInPrivacySettingsTasksServiceError struct {
	message string
}

func newUpdateFriendsInPrivacySettingsTasksServiceError(message string) *UpdateFriendsInPrivacySettingsTasksServiceError {
	return &UpdateFriendsInPrivacySettingsTasksServiceError{
		message: message,
	}
}

func (e *UpdateFriendsInPrivacySettingsTasksServiceError) Error() string {
	return e.message
}

var (
	ErrCreateUpdateFriendsInPrivacySettingsTasksError  = newUpdateFriendsInPrivacySettingsTasksServiceError("service.update_friends_in_privacy_settings.error.create_error")
	ErrProcessUpdateFriendsInPrivacySettingsTasksError = newUpdateFriendsInPrivacySettingsTasksServiceError("service.update_friends_in_privacy_settings.error.process_error")
)
