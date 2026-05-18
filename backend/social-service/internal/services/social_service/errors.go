package socialservice

type SocialServiceError struct {
	message string
}

func newSocialServiceError(message string) *SocialServiceError {
	return &SocialServiceError{
		message: message,
	}
}

func (e *SocialServiceError) Error() string {
	return e.message
}

var (
	ErrUserIdIsRequired                = newSocialServiceError("service.social.error.user_id_is_required")
	ErrUsernameIsRequired              = newSocialServiceError("service.social.error.username_is_required")
	ErrUserIdsIsRequired               = newSocialServiceError("service.social.error.user_ids_is_required")
	ErrUsernamesIsRequired             = newSocialServiceError("service.social.error.usernames_is_required")
	ErrTooManyUserIds                  = newSocialServiceError("service.social.error.too_many_user_ids")
	ErrTooManyUsernames                = newSocialServiceError("service.social.error.too_many_usernames")
	ErrUserNotFound                    = newSocialServiceError("service.social.error.user_not_found")
	ErrCannotAddYourselfToFriends      = newSocialServiceError("service.social.error.cannot_add_yourself_to_friends")
	ErrCannotAddDeletedUserToFriends   = newSocialServiceError("service.social.error.cannot_deleted_user_to_friends")
	ErrCannotRemoveYourselfFromFriends = newSocialServiceError("service.social.error.cannot_remove_yourself_from_friends")
	ErrCannotBlockYourself             = newSocialServiceError("service.social.error.cannot_block_yourself")
	ErrCannotUnblockYourself           = newSocialServiceError("service.social.error.cannot_unblock_yourself")
	ErrUsersNotFound                   = newSocialServiceError("service.social.error.users_not_found")
	ErrPagesizeTooLarge                = newSocialServiceError("service.social.error.pagesize_too_large")
	ErrSearchFilterIsRequired          = newSocialServiceError("service.social.error.search_filter_is_required")
	ErrSearchFilterTooShort            = newSocialServiceError("service.social.error.search_filter_too_short")
	ErrSearchFilterTooLong             = newSocialServiceError("service.social.error.search_filter_too_long")
	ErrFieldsAreRequired               = newSocialServiceError("service.social.error.fields_are_required")
	ErrPrivacySettingsAreRequired      = newSocialServiceError("service.social.error.privacy_settings_are_required")
	ErrInvalidPrivacySettingValue      = newSocialServiceError("service.social.error.invalid_privacy_setting_value")
	ErrPrivacyListTooLong              = newSocialServiceError("service.social.error.privacy_list_too_long")
	ErrUserNotFriend                   = newSocialServiceError("service.social.error.user_not_friend")
	ErrNothingToUpdate                 = newSocialServiceError("service.social.error.nothing_to_update")
)
