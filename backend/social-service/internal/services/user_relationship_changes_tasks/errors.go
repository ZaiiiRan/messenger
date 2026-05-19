package userrelationshipchangestasks

type UserRelationshipChangesTasksServiceError struct {
	message string
}

func newUserRelationshipChangesTasksServiceError(message string) *UserRelationshipChangesTasksServiceError {
	return &UserRelationshipChangesTasksServiceError{
		message: message,
	}
}

func (e *UserRelationshipChangesTasksServiceError) Error() string {
	return e.message
}

var (
	ErrMarshalPayload             = newUserRelationshipChangesTasksServiceError("service.user_relationship_changes_tasks.payload_marshal_error")
	ErrCreateUserDataDeletionTask = newUserRelationshipChangesTasksServiceError("service.user_relationship_changes_tasks.create_user_relationship_changes_task_error")
	ErrSendUserDataDeletionTasks  = newUserRelationshipChangesTasksServiceError("service.user_relationship_changes_tasks.send_user_user_relationship_changes_tasks_error")
)
