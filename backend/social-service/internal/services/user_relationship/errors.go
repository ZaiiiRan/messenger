package userrelationshipservice

type UserRelationshipServiceError struct {
	message string
}

func newUserRelationshipServiceError(message string) *UserRelationshipServiceError {
	return &UserRelationshipServiceError{
		message: message,
	}
}

func (e *UserRelationshipServiceError) Error() string {
	return e.message
}

var (
	ErrAddUserToFriends         = newUserRelationshipServiceError("service.user_relationship.error.add_user_to_friends_error")
	ErrAlreadyFriends           = newUserRelationshipServiceError("service.user_relationship.error.already_friends")
	ErrFriendRequestAlreadySent = newUserRelationshipServiceError("service.user_relationship.error.friend_request_already_sent")
	ErrBlockedByFriendCandidate = newUserRelationshipServiceError("service.user_relationship.error.blocked_by_friend_candidate")
	ErrRemoveFromFriends        = newUserRelationshipServiceError("service.user_relationship.error.remove_from_friends_error")
	ErrBlockUser                = newUserRelationshipServiceError("service.user_relationship.error.block_user_error")
	ErrAlreadyBlocked           = newUserRelationshipServiceError("service.user_relationship.error.already_blocked")
	ErrUnblockUser              = newUserRelationshipServiceError("service.user_relationship.error.unblock_user_error")
	ErrGetUserRelationship      = newUserRelationshipServiceError("service.user_relationship.error.get_user_relationship_error")
)
