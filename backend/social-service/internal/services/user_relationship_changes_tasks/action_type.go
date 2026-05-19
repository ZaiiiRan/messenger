package userrelationshipchangestasks

type ActionType string

const (
	AddToFriends      ActionType = "add_to_friends"
	RemoveFromFriends ActionType = "remove_from_friends"
	Block             ActionType = "block"
	Unblock           ActionType = "unblock"
)
