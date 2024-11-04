package user

import (
	pgDB "backend/internal/dbs/pgDB"
	appErr "backend/internal/errors/appError"
	"database/sql"
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
		if err.Error() != "inernal server error" {
			return nil, appErr.BadRequest("first" + err.Error())
		}
		return nil, err
	}
	if err := validateName(lastname); err != nil {
		if err.Error() != "inernal server error" {
			return nil, appErr.BadRequest("last" + err.Error())
		}
		return nil, err
	}
	if phone != nil {
		if err := validatePhone(*phone); err != nil {
			return nil, err
		}
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
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
			return appErr.InternalServerError("internal server error")
		}
	} else {
		// existing user
		query := `UPDATE users SET username=$1, email=$2, password=$3, phone=$4, firstname=$5, lastname=$6, 
					birthdate=$7, is_deleted=$8, is_banned=$9, is_activated=$10 WHERE id=$11`
		_, err := db.Exec(query, u.Username, u.Email, u.Password, u.Phone, u.Firstname, u.Lastname, u.Birthdate,
			u.IsDeleted, u.IsBanned, u.IsActivated, u.ID)
		if err != nil {
			return appErr.InternalServerError("internal server error")
		}
	}
	return nil
}

// Deletion user from DataBase (soft delete)
func (u *User) Delete() error {
	if u.ID == 0 {
		return appErr.BadRequest("user not found")
	}
	u.IsDeleted = true
	return u.Save()
}

// Check user password
func (u *User) CheckPassword(password string) bool {
	return comparePasswords(u.Password, password)
}

// Get user by id
func GetUserByID(ID uint64) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE id = $1`, ID)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.BadRequest("user not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// Get user by username
func GetUserByUsername(username string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE username = $1`, username)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.BadRequest("user not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// Get user by email
func GetUserByEmail(email string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE email = $1`, email)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.BadRequest("user not found")
	} else if err != nil {
		return nil, appErr.InternalServerError("internal server error")
	}
	return user, nil
}

// Get user by phone
func GetUserByPhone(phone string) (*User, error) {
	db := pgDB.GetDB()
	row := db.QueryRow(`SELECT * FROM users WHERE phone = $1`, phone)
	user, err := createUserFromSQLRow(row)
	if err == sql.ErrNoRows {
		return nil, appErr.BadRequest("user not found")
	} else if err != nil {
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

// hashing password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", appErr.InternalServerError("internal server error")
	}
	return string(hash), nil
}

// compare passwords
func comparePasswords(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// validate username
func validateUsername(username string) error {
	if username == "" {
		return appErr.BadRequest("username is empty")
	}

	candidate, err := GetUserByUsername(username)
	if err != nil && err.Error() != "user not found" {
		return appErr.InternalServerError("internal server error")
	}
	if candidate != nil {
		return appErr.BadRequest("user with the same username already exists")
	} else if len(username) < 5 {
		return appErr.BadRequest("username must be at least 5 characters")
	}
	return nil
}

// validate email
func validateEmail(email string) error {
	if email == "" {
		return appErr.BadRequest("email is empty")
	}

	candidate, err := GetUserByEmail(email)
	if err != nil && err.Error() != "user not found" {
		return appErr.InternalServerError("internal server error")
	}
	if candidate != nil {
		return appErr.BadRequest("user with the same email already exists")
	}

	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegex.MatchString(email) {
		return appErr.BadRequest("invalid email format")
	}
	return nil
}

// validate phone
func validatePhone(phone string) error {
	if phone == "" {
		return nil
	}

	candidate, err := GetUserByPhone(phone)
	if err != nil && err.Error() != "user not found" {
		return appErr.InternalServerError("internal server error")
	}
	if candidate != nil {
		return appErr.BadRequest("user with the same phone number already exists")
	}

	phoneRegex := regexp.MustCompile(`^\+7\(9\d{2}\)-\d{3}-\d{2}-\d{2}$`)
	if !phoneRegex.MatchString(phone) {
		return appErr.BadRequest("phone must be in format +7(9xx)-xxx-xx-xx or empty")
	}
	return nil
}

// validate names (firstname and lastname)
func validateName(name string) error {
	if name == "" {
		return appErr.BadRequest("name is empty")
	}

	nameRegex := regexp.MustCompile(`^[A-ZА-Я][a-zа-я]+(-[A-ZА-Я][a-zа-я]+)?$`)
	if !nameRegex.MatchString(name) {
		return appErr.BadRequest("name must start with a capital letter")
	}
	return nil
}

// validate password
func validatePassword(password string) error {
	if password == "" {
		return appErr.BadRequest("password is empty")
	}

	var (
		hasUpperCase   = regexp.MustCompile(`[A-ZА-ЯЁ]`).MatchString(password)
		hasLowerCase   = regexp.MustCompile(`[a-zа-яё]`).MatchString(password)
		hasNumber      = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecialChar = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	if len(password) < 8 {
		return appErr.BadRequest("password must be at least 8 characters")
	}
	if !hasUpperCase {
		return appErr.BadRequest("password must contain at least one uppercase letter")
	}
	if !hasLowerCase {
		return appErr.BadRequest("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return appErr.BadRequest("password must contain at least one digit")
	}
	if !hasSpecialChar {
		return appErr.BadRequest("password must contain at least one special character")
	}
	return nil
}
