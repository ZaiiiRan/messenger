package user

import (
	"backend/internal/dbs/pgDB"
	"backend/internal/errors/appError"
	"database/sql"
)

// Get activated not deleted and not banned users by username or email (for infinite scroll)
func GetUsersByUsernameOrEmail(search string, limit, offset int) ([]User, error) {
	db := pgDB.GetDB()

	query := `
		SELECT * FROM users
		WHERE (username ILIKE $1 OR email ILIKE $2) AND is_deleted = FALSE
		AND is_banned = FALSE AND is_activated = TRUE
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
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

// parsing users from sql rows
func createUsersFromSQLRows(rows *sql.Rows) ([]User, error) {
	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Phone,
			&user.Firstname, &user.Lastname, &user.Birthdate, &user.IsDeleted, &user.IsBanned, &user.IsActivated, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
