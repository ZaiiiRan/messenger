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

	isBlocked, err := checkBlock(userID, friendID)
	if err != nil {
		return err
	}
	if isBlocked {
		return appErr.BadRequest("you are blocked by this user")
	}

	db := pgDB.GetDB()
	var status string
	query := `
		SELECT fs.name FROM friends f 
        JOIN friend_statuses fs ON f.status_id = fs.id
        WHERE 
			(f.friend_1_id = $1 AND f.friend_2_id = $2) 
            OR (f.friend_1_id = $2 AND f.friend_2_id = $1)
	`
	err = db.QueryRow(query, userID, friend).Scan(&status)
	if err == sql.ErrNoRows {
		_, err = db.Exec(`INSERT INTO friends (friend_1_id, friend_2_id, status_id)
		VALUES ($1, $2, (SELECT id FROM friend_statuses WHERE name = 'request'))`, userID, friendID)
	} else if err == nil && status == "request" {
		_, err = db.Exec(`UPDATE friends SET status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted')
		WHERE (friend_1_id = $1 AND friend_2_id = $2) OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, friendID)
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
	_, err = db.Exec(`DELETE FROM friends WHERE (friend_1_id = $1 AND friend_2_id = $2)
	OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, friendID)
	if err != nil {
		return appErr.InternalServerError("internal server error")
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
	_, err = db.Exec(`DELETE FROM friends WHERE (friend_1_id = $1 AND friend_2_id = $2)
	OR (friend_1_id = $2 AND friend_2_id = $1)`, userID, targetID)
	if err != nil {
		return appErr.InternalServerError("inernal server error")
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
	_, err = db.Exec(`DELETE FROM friends WHERE friend_friend_1_id = $1 AND friend_2_id = $2 
    AND status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked'`, userID, targetID)
	if err != nil {
		return appErr.InternalServerError("inernal server error")
	}
	return nil
}

// checking whether the user is in the block list of the target
func checkBlock(userID, targetID uint64) (bool, error) {
	db := pgDB.GetDB()

	var blockingStatus string
	checkBlockingQuery := `
		SELECT fs.name FROM friends f 
		JOIN friend_statuses fs ON f.status_id = fs.id
		WHERE f.friend_1_id = $1 AND f.friend_2_id = $2 AND fs.name = 'blocked'
	`
	err := db.QueryRow(checkBlockingQuery, targetID, userID).Scan(&blockingStatus)
	if err != nil && err == sql.ErrNoRows {
		return false, nil
	} else if err == nil && blockingStatus == "blocked" {
		return true, nil
	} else if err != nil {
		return false, appErr.InternalServerError("internal server error")
	}
	return false, nil
}