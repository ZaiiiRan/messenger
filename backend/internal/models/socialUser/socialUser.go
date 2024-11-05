package socialUser

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
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
func AddFriend(userID, friendID uint64) error {
	friend, err := user.GetUserByID(friendID)
	if err != nil {
		return err
	}
	if !friend.IsActivated || friend.IsDeleted {
		return appErr.BadRequest("user not found")
	}
	if friend.IsBanned {
		return appErr.BadRequest("user is banned")
	}

	db := pgDB.GetDB()
	status, err := GetRelations(userID, friendID)
	if err != nil {
		return err
	}
	if status != nil && *status == "blocked by target" {
		return appErr.BadRequest("you are blocked by this user")
	} else if status != nil && *status == "blocked" {
		return appErr.BadRequest("you blocked this user")
	} else if status != nil && *status == "accepted" {
		return appErr.BadRequest("you are already friends")
	} else if status != nil && *status == "outgoing request" {
		return appErr.BadRequest("friend request has already been sent")
	} else if status != nil && *status == "incoming request" {
		_, err = db.Exec(`UPDATE friends SET status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted')
		WHERE (friend_1_id = $1 AND friend_2_id = $2) OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, friendID)
	} else {
		_, err = db.Exec(`INSERT INTO friends (friend_1_id, friend_2_id, status_id)
		VALUES ($1, $2, (SELECT id FROM friend_statuses WHERE name = 'request'))`, userID, friendID)
	}
	if err != nil {
		return appErr.InternalServerError("internal server error")
	}

	return nil
}

// Removing friend from friend list
func RemoveFriend(userID, friendID uint64) error {
	_, err := user.GetUserByID(friendID)
	if err != nil {
		return err
	}
	db := pgDB.GetDB()
	status, err := GetRelations(userID, friendID)
	if err != nil {
		return err
	}
	if status != nil {
		_, err = db.Exec(`DELETE FROM friends WHERE (friend_1_id = $1 AND friend_2_id = $2)
		OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, friendID)
		if err != nil {
			return appErr.InternalServerError("internal server error")
		}
	}
	return nil
}

// Adding user to block list
func BlockUser(userID, targetID uint64) error {
	_, err := user.GetUserByID(targetID)
	if err != nil {
		return err
	}
	db := pgDB.GetDB()
	status, err := GetRelations(userID, targetID)
	if err != nil {
		return err
	}
	if status != nil && *status == "blocked" {
		return appErr.BadRequest("you have already blocked this user")
	} else if status != nil && *status != "blocked by target" {
		_, err = db.Exec(`DELETE FROM friends WHERE (friend_1_id = $1 AND friend_2_id = $2)
		OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, targetID)
		if err != nil {
			return appErr.InternalServerError("inernal server error")
		}
	}

	_, err = db.Exec(`INSERT INTO friends (friend_1_id, friend_2_id, status_id)
    VALUES ($1, $2, (SELECT id FROM friend_statuses WHERE name = 'blocked'))`, userID, targetID)
	if err != nil {
		return appErr.InternalServerError("inernal server error")
	}
	return nil
}

// Removing user from block list
func UnblockUser(userID, targetID uint64) error {
	_, err := user.GetUserByID(targetID)
	if err != nil {
		return err
	}
	db := pgDB.GetDB()
	status, err := GetRelations(userID, targetID)
	if err != nil {
		return err
	}
	if status != nil && *status == "blocked" {
		_, err = db.Exec(`DELETE FROM friends WHERE friend_1_id = $1 AND friend_2_id = $2 
        AND status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked')`, userID, targetID)
		if err != nil {
			return appErr.InternalServerError("inernal server error")
		}
	}

	return nil
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
		LIMIT 1
	`

	var friendshipStatus string
	err := db.QueryRow(query, userID, targetID).Scan(&friendshipStatus)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
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
