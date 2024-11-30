package socialUser

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"backend/internal/models/user"
	"database/sql"
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

	db := pgDB.GetDB()
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
		_, err = db.Exec(`UPDATE friends SET status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted')
		WHERE (friend_1_id = $1 AND friend_2_id = $2) OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, friendID)
		newStatus = "accepted"
	} else {
		_, err = db.Exec(`INSERT INTO friends (friend_1_id, friend_2_id, status_id)
		VALUES ($1, $2, (SELECT id FROM friend_statuses WHERE name = 'request'))`, userID, friendID)
		newStatus = "outgoing request"
	}
	if err != nil {
		logger.GetInstance().Error(err.Error(), "add friend", map[string]interface{}{"userID": userID, "friendID": friendID}, err)
		return nil, appErr.InternalServerError("internal server error")
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
	db := pgDB.GetDB()
	if friend.FriendStatus != nil && *friend.FriendStatus != "blocked by target" && *friend.FriendStatus != "blocked" {
		_, err = db.Exec(`DELETE FROM friends WHERE (friend_1_id = $1 AND friend_2_id = $2)
		OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, friendID)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "remove friend", map[string]interface{}{"userID": userID, "friendID": friendID}, err)
			return nil, appErr.InternalServerError("internal server error")
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

	db := pgDB.GetDB()
	var newStatus string

	if target.FriendStatus != nil && *target.FriendStatus == "blocked" {
		return nil, appErr.BadRequest("you have already blocked this user")
	} else if target.FriendStatus != nil && *target.FriendStatus != "blocked by target" {
		_, err = db.Exec(`DELETE FROM friends WHERE (friend_1_id = $1 AND friend_2_id = $2)
		OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, targetID)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "block user (deleting)", map[string]interface{}{"userID": userID, "targetID": targetID}, err)
			return nil, appErr.InternalServerError("inernal server error")
		}
	}

	_, err = db.Exec(`INSERT INTO friends (friend_1_id, friend_2_id, status_id)
    VALUES ($1, $2, (SELECT id FROM friend_statuses WHERE name = 'blocked'))`, userID, targetID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "block user (blocking)", map[string]interface{}{"userID": userID, "targetID": targetID}, err)
		return nil, appErr.InternalServerError("inernal server error")
	}
	newStatus = "blocked"
	target.FriendStatus = &newStatus
	return target, nil
}

// Removing user from block list
func UnblockUser(userID, targetID uint64) (*SocialUser, error) {
	target, err := GetTargetByID(userID, targetID)
	if err != nil {
		return nil, err
	}
	db := pgDB.GetDB()

	if target.FriendStatus != nil && *target.FriendStatus == "blocked" {
		_, err = db.Exec(`DELETE FROM friends WHERE friend_1_id = $1 AND friend_2_id = $2 
        AND status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked')`, userID, targetID)
		if err != nil {
			logger.GetInstance().Error(err.Error(), "unblock user", map[string]interface{}{"userID": userID, "targetID": targetID}, err)
			return nil, appErr.InternalServerError("inernal server error")
		}
	}

	status, err := GetRelations(userID, targetID)
	if err != nil {
		return nil, err
	}
	target.FriendStatus = status

	return target, nil
}

// Get relations between two users
func GetRelations(userID, targetID uint64) (*string, error) {
	db := pgDB.GetDB()
	query := `
		SELECT 
			CASE 
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted') THEN 'accepted'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') 
					AND f.friend_1_id = $1 THEN 'blocked'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') 
					AND f.friend_2_id = $1 THEN 'blocked by target'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') 
					AND f.friend_1_id = $1 THEN 'outgoing request'
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') 
					AND f.friend_2_id = $1 THEN 'incoming request'
				ELSE NULL
			END AS friendship_status
		FROM friends f
		WHERE 
			(f.friend_1_id = $1 AND f.friend_2_id = $2)
			OR (f.friend_1_id = $2 AND f.friend_2_id = $1)
		ORDER BY 
			CASE 
				WHEN f.friend_1_id = $1 AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') THEN 1
				WHEN f.friend_2_id = $1 AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') THEN 2
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted') THEN 3
				WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') THEN 4
				ELSE 5
			END
		LIMIT 1
	`

	var friendshipStatus string
	err := db.QueryRow(query, userID, targetID).Scan(&friendshipStatus)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get relations between users", map[string]interface{}{"userID": userID, "targetID": targetID}, err)
		return nil, appErr.InternalServerError("internal server error")
	}

	return &friendshipStatus, nil
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
