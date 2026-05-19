package userrelationship

type UserRelationshipStatus int16

var (
	None             UserRelationshipStatus = 0
	FriendRequestBy1 UserRelationshipStatus = 1
	FriendRequestBy2 UserRelationshipStatus = 2
	Friends          UserRelationshipStatus = 3
	BlockedBy1       UserRelationshipStatus = 4
	BlockedBy2       UserRelationshipStatus = 5
	BlockedByBoth    UserRelationshipStatus = 6
)

func (s UserRelationshipStatus) String() string {
	switch s {
	case None:
		return "none"
	case FriendRequestBy1:
		return "friend_request_by_1"
	case FriendRequestBy2:
		return "friend_request_by_2"
	case Friends:
		return "friends"
	case BlockedBy1:
		return "blocked_by_1"
	case BlockedBy2:
		return "blocked_by_2"
	case BlockedByBoth:
		return "blocked_by_both"
	default:
		return ""
	}
}

func ToUserRelationshipStatus(value string) (UserRelationshipStatus, error) {
	switch value {
	case "none":
		return None, nil
	case "friend_request_by_1":
		return FriendRequestBy1, nil
	case "friend_request_by_2":
		return FriendRequestBy2, nil
	case "friends":
		return Friends, nil
	case "blocked_by_1":
		return BlockedBy1, nil
	case "blocked_by_2":
		return BlockedBy2, nil
	case "blocked_by_both":
		return BlockedByBoth, nil
	default:
		return FriendRequestBy1, ErrUnknownUserRelationshipStatus
	}
}
