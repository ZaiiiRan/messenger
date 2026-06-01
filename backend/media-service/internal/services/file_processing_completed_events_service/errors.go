package fileprocessingcompletedeventsservice

type FileProcessingCompletedEventsServiceError struct {
	message string
}

func newFileProcessingCompletedEventsServiceError(message string) *FileProcessingCompletedEventsServiceError {
	return &FileProcessingCompletedEventsServiceError{
		message: message,
	}
}

func (e *FileProcessingCompletedEventsServiceError) Error() string {
	return e.message
}

var (
	ErrCreateFileProcessingCompletedEvents = newFileProcessingCompletedEventsServiceError("service.file_processing_completed_events.create_file_processing_completed_events_error")
	ErrSendFileProcessingCompletedEvents   = newFileProcessingCompletedEventsServiceError("service.file_processing_completed_events.send_file_processing_completed_events_error")
)
