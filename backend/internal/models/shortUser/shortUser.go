package shortUser

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"database/sql"
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

// All users searching by username or email
func SearchAll(userID uint64, search string, limit, offset int) ([]ShortUser, error) {
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname,
			u.is_deleted, u.is_banned, u.is_activated
		FROM users u
		WHERE
			u.id != $1
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE $2 OR u.email ILIKE $2)
		ORDER BY (u.username = $3 OR u.email = $3) DESC, u.username
		LIMIT $4 OFFSET $5
	`

	return queryUsers(query, userID, "%"+search+"%", search, limit, offset)
}

// Friends searching by username or email
func SearchFriends(userID uint64, search string, limit, offset int) ([]ShortUser, error) {
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname,
			u.is_deleted, u.is_banned, u.is_activated
		FROM users u
		JOIN friends f ON (f.friend_1_id = u.id OR f.friend_2_id = u.id)
		WHERE
			u.id != $1
			AND (f.friend_1_id = $1 OR f.friend_2_id = $1)
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted')
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC, u.username
		LIMIT $3 OFFSET $4
	`

	return queryUsers(query, userID, search, limit, offset)
}

// Incoming friend requests searching by username or email
func SearchIncomingFriendRequests(userID uint64, search string, limit, offset int) ([]ShortUser, error) {
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname,
			u.is_deleted, u.is_banned, u.is_activated
		FROM users u
		JOIN friends f ON f.friend_1_id = u.id
		WHERE 
			u.id != $1
			AND f.friend_2_id = $1
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request')
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC, u.username
		LIMIT $3 OFFSET $4
	`

	return queryUsers(query, userID, search, limit, offset)
}

// Outgoing friend requests searching by username or email
func SearchOutgoingFriendRequests(userID uint64, search string, limit, offset int) ([]ShortUser, error) {
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname,
			u.is_deleted, u.is_banned, u.is_activated
		FROM users u
		JOIN friends f ON f.friend_2_id = u.id
		WHERE 
			u.id != $1
			AND f.friend_1_id = $1
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request')
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC, u.username
		LIMIT $3 OFFSET $4
	`

	return queryUsers(query, userID, search, limit, offset)
}

// Block list searching by username or email
func SearchBlockList(userID uint64, search string, limit, offset int) ([]ShortUser, error) {
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname,
			u.is_deleted, u.is_banned, u.is_activated
		FROM users u
		JOIN friends f ON f.friend_2_id = u.id
		WHERE 
			u.id != $1
			AND f.friend_1_id = $1
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked')
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC, u.username
		LIMIT $3 OFFSET $4
	`

	return queryUsers(query, userID, search, limit, offset)
}

// parsing users from sql rows
func createUsersFromSQLRows(rows *sql.Rows) ([]ShortUser, error) {
	var users []ShortUser

	for rows.Next() {
		var user ShortUser
		err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname,
			&user.IsDeleted, &user.IsBanned, &user.IsActivated)
		if err != nil {
			return nil, appErr.InternalServerError("internal server error")
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	return users, nil
}

// execute query to db
func queryUsers(query string, params ...interface{}) ([]ShortUser, error) {
	db := pgDB.GetDB()

	rows, err := db.Query(query, params...)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("users not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	defer rows.Close()

	users, err := createUsersFromSQLRows(rows)
	if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}

	if len(users) == 0 {
		return nil, appErr.NotFound("users not found")
	}

	return users, nil
}