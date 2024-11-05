package shortUser

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"database/sql"
)

type ShortUser struct {
	ID    uint64 `json:"user_id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func Search(userID uint64, search string, limit, offset int) ([]ShortUser, error) {
	db := pgDB.GetDB()
	query := `
		SELECT u.id, u.username, u.firstname, u.lastname
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

	rows, err := db.Query(query, userID, "%"+search+"%", search, limit, offset)
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

// parsing users from sql rows
func createUsersFromSQLRows(rows *sql.Rows) ([]ShortUser, error) {
	var users []ShortUser

	for rows.Next() {
		var user ShortUser
		err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname)
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