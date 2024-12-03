package socialUser

import (
	"backend/internal/models/user"
)

type SocialUser struct {
	User         *user.UserDTO `json:"user"`
	FriendStatus *string       `json:"friend_status"`
}

// Creating SocialUser object
func CreateSocialUser(dto *user.UserDTO, friendStatus *string) *SocialUser {
	return &SocialUser{
		User:         dto,
		FriendStatus: friendStatus,
	}
}

// Get relations between two users
func GetRelations(userID, targetID uint64) (*string, error) {
	return getRelationsFromDB(userID, targetID)
}

// Get Social User
func GetTargetByID(userID, targetID uint64) (*SocialUser, error) {
	target, err := user.GetUserByID(targetID)
	if err != nil {
		return nil, err
	}
	targetDTO := user.CreateUserDTOFromUserObj(target)
	status, err := GetRelations(userID, targetID)
	if err != nil {
		return nil, err
	}
	socialTarget := CreateSocialUser(targetDTO, status)
	return socialTarget, nil
}
