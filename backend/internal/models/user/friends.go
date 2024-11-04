package user

import (
	"backend/internal/dbs/pgDB"
	"backend/internal/errors/appError"
	"database/sql"
)

// Get user friends by userID
func GetUserFriends(userID uint64, limit, offset int) ([]User, error) {
	db := pgDB.GetDB()

	query := `
		SELECT u.* FROM friends f
		JOIN users u ON (u.id = f.friend_1_id OR u.id = f.friend_2_id)
		JOIN friend_statuses fs ON fs.id = f.status_id
		WHERE fs.name = 'accepted'
		AND (f.friend_1_id = $1 OR f.friend_2_id = $1)
		AND u.id != $1
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, userID, limit, offset)
	if err == sql.ErrNoRows {
		return nil, appError.NotFound("users not found")
	} else if err != nil {
		return nil, appError.InternalServerError("internal server error")
	}
	defer rows.Close()

	users, err := createUsersFromSQLRows(rows)
	if err != nil {
		return nil, appError.InternalServerError("internal server error")
	}

	return users, nil
}

// Get incoming friend requests by userID
func GetUserIncomingFriendRequests(userID uint64, limit, offset int) ([]User, error) {
	db := pgDB.GetDB()

	query := `
		SELECT u.* FROM friends f
		JOIN users u ON u.id = f.friend_1_id
		JOIN friend_statuses fs ON fs.id = f.status_id
		WHERE fs.name = 'request'
		AND f.friend_2_id = $1
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, userID, limit, offset)
	if err == sql.ErrNoRows {
		return nil, appError.NotFound("users not found")
	} else if err != nil {
		return nil, appError.InternalServerError("internal server error")
	}
	defer rows.Close()

	users, err := createUsersFromSQLRows(rows)
	if err != nil {
		return nil, appError.InternalServerError("internal server error")
	}

	return users, nil
}

// Get outgoing friend requests by userID
func GetUserOutgoingFriendRequests(userID uint64, limit, offset int) ([]User, error) {
	db := pgDB.GetDB()

	query := `
		SELECT u.* FROM friends f
		JOIN users u ON u.id = f.friend_2_id
		JOIN friend_statuses fs ON fs.id = f.status_id
		WHERE fs.name = 'request'
		AND f.friend_1_id = $1
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, userID, limit, offset)
	if err == sql.ErrNoRows {
		return nil, appError.NotFound("users not found")
	} else if err != nil {
		return nil, appError.InternalServerError("internal server error")
	}
	defer rows.Close()

	users, err := createUsersFromSQLRows(rows)
	if err != nil {
		return nil, appError.InternalServerError("internal server error")
	}

	return users, nil
}


