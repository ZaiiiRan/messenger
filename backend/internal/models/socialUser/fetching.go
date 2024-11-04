package socialUser

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/models/user"
	"database/sql"
)

// General user search by username or email
func GetUsersByUsernameOrEmail(userID uint64, search string, limit, offset int) ([]SocialUser, error) {
	query := `
		SELECT u.id, u.username, u.email, u.firstname, u.lastname, u.phone,
		u.birthdate, u.is_deleted, u.is_banned, u.is_activated,
		CASE
			WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted') THEN 'accepted'
			WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked') THEN 'blocked'
			WHEN f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request') THEN 'request'
			ELSE NULL
		END AS friend_status
		FROM users u
		LEFT JOIN friends f ON (f.friend_1_id = $1 AND f.friend_2_id = u.id)
			OR (f.friend_1_id = u.id AND f.friend_2_id = $1)
		WHERE 
			u.id != $1
			AND ($2 = '' OR u.username ILIKE $2 OR u.email ILIKE $2)
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
		ORDER BY
			(u.username = $3 OR u.email = $3) DESC, u.username
		LIMIT $4 OFFSET $5
	`
	return queryUsers(query, userID, "%"+search+"%", search, limit, offset)
}

// Get user friends
func GetUserFriendsByUsernameOrEmail(userID uint64, search string, limit, offset int) ([]SocialUser, error) {
	query := `
		SELECT u.id, u.username, u.email, u.firstname, u.lastname, u.phone,
		u.birthdate, u.is_deleted, u.is_banned, u.is_activated,
		'accepted' AS friend_status
		FROM users u
		JOIN friends f ON (f.friend_1_id = u.id OR f.friend_2_id = u.id)
		WHERE 
			(f.friend_1_id = $1 OR f.friend_2_id = $1)
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'accepted')
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC
		LIMIT $3 OFFSET $4
	`
	return queryUsers(query, userID, search, limit, offset)
}

// Get incoming friend requests
func GetUserIncomingFriendRequestsByUsernameOrEmail(userID uint64, search string, limit, offset int) ([]SocialUser, error) {
	query := `
		SELECT u.id, u.username, u.email, u.firstname, u.lastname, u.phone, u.birthdate, 
			'request' AS friend_status
		FROM users u
		JOIN friends f ON f.friend_1_id = u.id
		WHERE 
			f.friend_2_id = $1
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request')
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC
		LIMIT $3 OFFSET $4
	`
	return queryUsers(query, userID, search, limit, offset)
}

// Get outgoing friend requests
func GetUserOutgoingFriendRequestsByUsernameOrEmail(userID uint64, search string, limit, offset int) ([]SocialUser, error) {
	query := `
		SELECT u.id, u.username, u.email, u.firstname, u.lastname, u.phone, u.birthdate, 
			'request' AS friend_status
		FROM users u
		JOIN friends f ON f.friend_2_id = u.id
		WHERE 
			f.friend_1_id = $1
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'request')
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC
		LIMIT $3 OFFSET $4
	`
	return queryUsers(query, userID, search, limit, offset)
}

// Get blocked users
func GetUserBlockListByUsernameOrEmail(userID uint64, search string, limit, offset int) ([]SocialUser, error) {
	query := `
		SELECT u.id, u.username, u.email, u.firstname, u.lastname, u.phone, u.birthdate, 
			'blocked' AS friend_status
		FROM users u
		JOIN friends f ON (f.friend_1_id = u.id OR f.friend_2_id = u.id)
		WHERE 
			(f.friend_1_id = $1 OR f.friend_2_id = $1)
			AND f.status_id = (SELECT id FROM friend_statuses WHERE name = 'blocked')
			AND u.is_deleted = FALSE
			AND u.is_banned = FALSE
			AND u.is_activated = TRUE
			AND ($2 = '' OR u.username ILIKE '%' || $2 || '%' OR u.email ILIKE '%' || $2 || '%')
		ORDER BY (u.username = $2 OR u.email = $2) DESC
		LIMIT $3 OFFSET $4
	`
	return queryUsers(query, userID, search, limit, offset)
}

// parsing users from sql rows
func createUsersFromSQLRows(rows *sql.Rows) ([]SocialUser, error) {
	var users []SocialUser

	for rows.Next() {
		var user user.UserDTO
		var friendStatus string
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Firstname, &user.Lastname, &user.Phone,
			&user.Birthdate, &user.IsDeleted, &user.IsBanned, &user.IsActivated, &friendStatus)
		if err != nil {
			return nil, err
		}

		socialUser := CreateSocialUser(&user, &friendStatus)

		users = append(users, *socialUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// helper function to execute a user query and return the results
func queryUsers(query string, params ...interface{}) ([]SocialUser, error) {
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

	return users, nil
}
