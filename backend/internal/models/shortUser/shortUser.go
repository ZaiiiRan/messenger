package shortUser

import (
	"backend/internal/models/user"
	"backend/internal/models/user/userDTO"
)

type ShortUser struct {
	ID          uint64 `json:"user_id"`
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	IsDeleted   bool   `json:"is_deleted"`
	IsBanned    bool   `json:"is_banned"`
	IsActivated bool   `json:"is_activated"`
}

// Creating Short user object from user dto object
func CreateShortUserFromUserDTO(dto *userDTO.UserDTO) *ShortUser {
	return &ShortUser{
		ID:          dto.ID,
		Username:    dto.Username,
		Firstname:   dto.Firstname,
		Lastname:    dto.Lastname,
		IsDeleted:   dto.IsDeleted,
		IsBanned:    dto.IsBanned,
		IsActivated: dto.IsActivated,
	}
}

// Creating Short user object from user object
func CreateShortUserFromUser(user *user.User) *ShortUser {
	return &ShortUser{
		ID:          user.ID,
		Username:    user.Username,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		IsDeleted:   user.IsDeleted,
		IsBanned:    user.IsBanned,
		IsActivated: user.IsActivated,
	}
}

// All users searching by username or email
func SearchAll(actorID uint64, search string, limit, offset int) ([]ShortUser, error) {
	return getAllUsersFromDB(actorID, search, limit, offset)
}

// Friends searching by username or email
func SearchFriends(actorID uint64, search string, limit, offset int) ([]ShortUser, error) {
	return getFriendsFromDB(actorID, search, limit, offset)
}

// Finding friends who are not in chat
func SearchFriendsAreNotChatting(actorID, chatID uint64, search string, limit, offset int) ([]ShortUser, error) {
	return getFriendsAreNotChattingFromDB(actorID, chatID, search, limit, offset)
}

// Incoming friend requests searching by username or email
func SearchIncomingFriendRequests(actorID uint64, search string, limit, offset int) ([]ShortUser, error) {
	return getIncomingFriendRequestsFromDB(actorID, search, limit, offset)
}

// Outgoing friend requests searching by username or email
func SearchOutgoingFriendRequests(actorID uint64, search string, limit, offset int) ([]ShortUser, error) {
	return getOutgoingFriendRequests(actorID, search, limit, offset)
}

// Block list searching by username or email
func SearchBlockList(actorID uint64, search string, limit, offset int) ([]ShortUser, error) {
	return getBlockedUsersFromDB(actorID, search, limit, offset)
}
