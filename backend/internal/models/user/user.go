package user

import (
	pgDB "backend/internal/dbs/pgDB"
	"database/sql"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint64     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	Phone       *string    `json:"phone"`
	Firstname   string     `json:"firstname"`
	Lastname    string     `json:"lastname"`
	Birthdate   *time.Time `json:"birthdate"`
	IsDeleted   bool       `json:"is_deleted"`
	IsBanned    bool       `json:"is_banned"`
	IsActivated bool       `json:"is_activated"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Get user by id
func GetUserByID(ID uint64) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE id = $1`, ID)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

// Get user by username
func GetUserByUsername(username string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE username = $1`, username)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

// Get user by email
func GetUserByEmail(email string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE email = $1`, email)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

// Get user by phone
func GetUserByPhone(phone string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE phone = $1`, phone)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
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

// hashing password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	return string(hash), err
}

// compare passwords
func comparePasswords(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}