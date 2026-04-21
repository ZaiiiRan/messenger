package socialUser

import (
	appErr "backend/internal/errors/appError"
)

// Adding to friend list
func AddFriend(userID, friendID uint64) (*SocialUser, error) {
	friend, err := GetTargetByID(userID, friendID)
	if err != nil {
		return nil, err
	}
	if !friend.User.IsActivated || friend.User.IsDeleted {
		return nil, appErr.NotFound("user not found")
	}
	if friend.User.IsBanned {
		return nil, appErr.BadRequest("user is banned")
	}

	var newStatus string

	if friend.FriendStatus != nil && *friend.FriendStatus == "blocked by target" {
		return nil, appErr.BadRequest("you are blocked by this user")
	} else if friend.FriendStatus != nil && *friend.FriendStatus == "blocked" {
		return nil, appErr.BadRequest("you blocked this user")
	} else if friend.FriendStatus != nil && *friend.FriendStatus == "accepted" {
		return nil, appErr.BadRequest("you are already friends")
	} else if friend.FriendStatus != nil && *friend.FriendStatus == "outgoing request" {
		return nil, appErr.BadRequest("friend request has already been sent")
	} else if friend.FriendStatus != nil && *friend.FriendStatus == "incoming request" {
		newStatus = "accepted"
		newStatusID, err := getFriendStatusIDFromDB(newStatus)
		if err != nil {
			return nil, err
		}
		err = updateFriendStatusInDB(userID, friendID, newStatusID)
		if err != nil {
			return nil, err
		}
	} else {
		newStatus = "request"
		newStatusID, err := getFriendStatusIDFromDB(newStatus)
		if err != nil {
			return nil, err
		}
		err = insertFriendDataToDB(userID, friendID, newStatusID)
		if err != nil {
			return nil, err
		}
		newStatus = "outgoing request"
	}
	friend.FriendStatus = &newStatus

	return friend, nil
}

// Removing friend from friend list
func RemoveFriend(userID, friendID uint64) (*SocialUser, error) {
	friend, err := GetTargetByID(userID, friendID)
	if err != nil {
		return nil, err
	}

	if friend.FriendStatus == nil {
		return nil, appErr.BadRequest("this user is not your friend")
	}

	if friend.FriendStatus != nil && *friend.FriendStatus != "blocked by target" && *friend.FriendStatus != "blocked" {
		err = removeFriendDataFromDB(userID, friendID)
		if err != nil {
			return nil, err
		}
		friend.FriendStatus = nil
	}
	return friend, nil
}

// Adding user to block list
func BlockUser(userID, targetID uint64) (*SocialUser, error) {
	target, err := GetTargetByID(userID, targetID)
	if err != nil {
		return nil, err
	}

	if target.FriendStatus != nil && *target.FriendStatus == "blocked" {
		return nil, appErr.BadRequest("you have already blocked this user")
	} else if target.FriendStatus != nil && *target.FriendStatus != "blocked by target" {
		err = removeFriendDataFromDB(userID, targetID)
		if err != nil {
			return nil, err
		}
	}

	newStatus := "blocked"
	newStatusID, err := getFriendStatusIDFromDB(newStatus)
	if err != nil {
		return nil, err
	}
	err = insertFriendDataToDB(userID, targetID, newStatusID)
	if err != nil {
		return nil, err
	}
	target.FriendStatus = &newStatus
	return target, nil
}

// Removing user from block list
func UnblockUser(userID, targetID uint64) (*SocialUser, error) {
	target, err := GetTargetByID(userID, targetID)
	if err != nil {
		return nil, err
	}

	if target.FriendStatus == nil || (target.FriendStatus != nil && (*target.FriendStatus != "blocked")) {
		return nil, appErr.BadRequest("this user is not blocked")
	}

	if target.FriendStatus != nil && *target.FriendStatus == "blocked" {
		status := "blocked"
		statusID, err := getFriendStatusIDFromDB(status)
		if err != nil {
			return nil, err
		}
		err = removeFriendDataByStatusIDFromDB(userID, targetID, statusID)
		if err != nil {
			return nil, err
		}
	}

	status, err := GetRelations(userID, targetID)
	if err != nil {
		return nil, err
	}
	target.FriendStatus = status

	return target, nil
}
