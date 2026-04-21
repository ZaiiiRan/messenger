package user

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
	"database/sql"
)

// get user by id from db
func getUserByIDFromDB(id uint64) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE id = $1`, id)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("user not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get user by id", id, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// get user by username from db
func getUserByUsernameFromDB(username string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE username = $1`, username)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("user not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get user by username", username, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// get user by email from db
func getUserByEmailFromDB(email string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE email = $1`, email)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("user not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get user by email", email, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// get user by phone from db
func getUserByPhoneFromDB(phone string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE phone = $1`, phone)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.NotFound("user not found")
	} else if err != nil {
		logger.GetInstance().Error(err.Error(), "get user by phone", phone, err)
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// parsing user from sql row
func createUserFromSQLRow(row *sql.Row) (*User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Phone,
		&user.Firstname, &user.Lastname, &user.Birthdate,
		&user.IsDeleted, &user.IsBanned, &user.IsActivated, &user.CreatedAt)
	return &user, err
}

// insert user to db
func insertUserToDB(u *User) error {
	db := pgDB.GetDB()
	query := `INSERT INTO users (username, email, password, phone, firstname, lastname, birthdate, is_deleted, is_banned, is_activated, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at`
	err := db.QueryRow(query, u.Username, u.Email, u.Password, u.Phone, u.Firstname, u.Lastname, u.Birthdate,
		u.IsDeleted, u.IsBanned, u.IsActivated, u.CreatedAt).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "user inserting", u, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

// update user in db
func updateUserInDB(u *User) error {
	db := pgDB.GetDB()
	query := `UPDATE users SET username=$1, email=$2, password=$3, phone=$4, firstname=$5, lastname=$6, 
					birthdate=$7, is_deleted=$8, is_banned=$9, is_activated=$10 WHERE id=$11`
	_, err := db.Exec(query, u.Username, u.Email, u.Password, u.Phone, u.Firstname, u.Lastname, u.Birthdate,
		u.IsDeleted, u.IsBanned, u.IsActivated, u.ID)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "user updating", u, err)
		return appErr.InternalServerError("internal server error")
	}
	return nil
}

