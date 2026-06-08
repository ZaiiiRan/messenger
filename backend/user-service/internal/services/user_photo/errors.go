package userphotosvc

type userPhotoServiceError struct {
	message string
}

func (e *userPhotoServiceError) Error() string { return e.message }

func newErr(msg string) *userPhotoServiceError { return &userPhotoServiceError{message: msg} }

var (
	ErrPhotoNotFound      = newErr("service.user_photo.error.photo_not_found")
	ErrPhotoNotOwned      = newErr("service.user_photo.error.photo_not_owned")
	ErrFileNotFound       = newErr("service.user_photo.error.file_not_found")
	ErrFileNotOwned       = newErr("service.user_photo.error.file_not_owned")
	ErrInvalidFileBucket  = newErr("service.user_photo.error.invalid_file_bucket")
	ErrFileNotReady       = newErr("service.user_photo.error.file_not_ready")
	ErrInternalError      = newErr("service.user_photo.error.internal")
)
