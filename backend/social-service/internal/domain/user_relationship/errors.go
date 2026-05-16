package userrelationship

var (
	ErrUnknownUserRelationshipStatus = NewUserRelationshipError("domain.user_relationship.error.unknown_user_relationship_status")
	ErrCannotBecomeFriends           = NewUserRelationshipError("domain.user_relationship.error.cannot_become_friends")
	ErrCannotBeMutualBlock           = NewUserRelationshipError("domain.user_relationship.error.cannot_be_mutual_block")
)
