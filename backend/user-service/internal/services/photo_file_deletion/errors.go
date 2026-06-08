package photofiledeletionsvc

type photoFileDeletionError struct {
	message string
}

func (e *photoFileDeletionError) Error() string { return e.message }

func newErr(msg string) *photoFileDeletionError { return &photoFileDeletionError{message: msg} }

var (
	ErrCreateTasks    = newErr("service.photo_file_deletion.error.create_tasks_failed")
	ErrSendTasks      = newErr("service.photo_file_deletion.error.send_tasks_failed")
	ErrMarshalPayload = newErr("service.photo_file_deletion.error.marshal_payload_failed")
)
