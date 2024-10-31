package user

import (
	pgDB "backend/internal/dbs/pgDB"
	"database/sql"
	"errors"
	"regexp"
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

// Creating user object
func CreateUser(username, email, password, firstname, lastname string, phone *string, birthdate *time.Time) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	if err := validateName(firstname); err != nil {
		return nil, err
	}
	if err := validateName(lastname); err != nil {
		return nil, err
	}
	if phone != nil {
		if err := validatePhone(*phone); err != nil {
			return nil, err
		}
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	user := &User{
		Username:    username,
		Email:       email,
		Password:    hashedPassword,
		Phone:       phone,
		Firstname:   firstname,
		Lastname:    lastname,
		Birthdate:   birthdate,
		IsDeleted:   false,
		IsBanned:    false,
		IsActivated: false,
	}
	return user, nil
}

// Saving user in DataBase
func (u *User) Save() error {
	db := pgDB.GetDB()
	if u.ID == 0 {
		// new user
		query := `INSERT INTO users (username, email, password, phone, firstname, lastname, birthdate, is_deleted, is_banned, is_activated, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at`
		err := db.QueryRow(query, u.Username, u.Email, u.Password, u.Phone, u.Firstname, u.Lastname, u.Birthdate,
			u.IsDeleted, u.IsBanned, u.IsActivated, u.CreatedAt).Scan(&u.ID, &u.CreatedAt)
		if err != nil {
			return errors.New("internal server error")
		}
	} else {
		// existing user
		query := `UPDATE users SET username=$1, email=$2, password=$3, phone=$4, firstname=$5, lastname=$6, 
					birthdate=$7, is_deleted=$8, is_banned=$9, is_activated=$10 WHERE id=$11`
		_, err := db.Exec(query, u.Username, u.Email, u.Password, u.Phone, u.Firstname, u.Lastname, u.Birthdate,
			u.IsDeleted, u.IsBanned, u.IsActivated, u.ID)
		if err != nil {
			return errors.New("internal server error")
		}
	}
	return nil
}

// Deletion user from DataBase (soft delete)
func (u *User) Delete() error {
	if u.ID == 0 {
		return errors.New("user not found")
	}
	u.IsDeleted = true
	return u.Save()
}

// Get user by id
func GetUserByID(ID uint64) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE id = $1`, ID)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, errors.New("internal server error")
	}
	return user, nil
}

// Get user by username
func GetUserByUsername(username string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE username = $1`, username)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, errors.New("internal server error")
	}
	return user, nil
}

// Get user by email
func GetUserByEmail(email string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE email = $1`, email)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, errors.New("internal server error")
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
		return nil, errors.New("internal server error")
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

// validate username
func validateUsername(username string) error {
	candidate, err := GetUserByUsername(username)
	if candidate != nil {
		return errors.New("user with the same username already exists")
	} else if err != nil {
		return errors.New("internal server error")
	}

	if len(username) < 5 {
		return errors.New("username must be at least 5 characters")
	}
	return nil
}

// validate email
func validateEmail(email string) error {
	candidate, err := GetUserByEmail(email)
	if candidate != nil {
		return errors.New("user with the same email already exists")
	} else if err != nil {
		return errors.New("internal server error")
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// validate phone
func validatePhone(phone string) error {
	if phone == "" {
		return nil
	}

	candidate, err := GetUserByPhone(phone)
	if candidate != nil {
		return errors.New("user with the same phone number already exists")
	} else if err != nil {
		return errors.New("internal server error")
	}

	phoneRegex := regexp.MustCompile(`^\+7\(9\d{2}\)-\d{3}-\d{2}-\d{2}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("phone must be in format +7(9xx)-xxx-xx-xx or empty")
	}
	return nil
}

// validate names (firstname and lastname)
func validateName(name string) error {
	nameRegex := regexp.MustCompile(`^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$`)
	if !nameRegex.MatchString(name) {
		return errors.New("name must start with a capital letter")
	}
	return nil
}

// validate password
func validatePassword(password string) error {
	var (
		hasUpperCase   = regexp.MustCompile(`[A-ZА-ЯЁ]`).MatchString(password)
		hasLowerCase   = regexp.MustCompile(`[a-zа-яё]`).MatchString(password)
		hasNumber      = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecialChar = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !hasUpperCase {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLowerCase {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one digit")
	}
	if !hasSpecialChar {
		return errors.New("password must contain at least one special character")
	}
	return nil
}