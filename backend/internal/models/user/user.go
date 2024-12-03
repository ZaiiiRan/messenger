package user

import (
	appErr "backend/internal/errors/appError"
	"backend/internal/logger"
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
	if err := validateAllFields(username, email, password, firstname, lastname, phone, birthdate); err != nil {
		return nil, err
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

// validate all fields
func validateAllFields(username, email, password, firstname, lastname string, phone *string, birthdate *time.Time) error {
	if err := validateUsername(username); err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	if err := validatePassword(password); err != nil {
		return err
	}
	if err := validateName(firstname); err != nil {
		if err.Error() != "inernal server error" {
			return appErr.BadRequest("first" + err.Error())
		}
		return err
	}
	if err := validateName(lastname); err != nil {
		if err.Error() != "inernal server error" {
			return appErr.BadRequest("last" + err.Error())
		}
		return err
	}
	if phone != nil {
		if err := validatePhone(*phone); err != nil {
			return err
		}
		if *phone == "" {
			phone = nil
		}
	}
	if birthdate != nil {
		if err := validateBirthdate(birthdate); err != nil {
			return err
		}
	}
	return nil
}

// Saving user in DataBase
func (u *User) Save() error {
	if u.ID == 0 {
		// new user
		err := insertUserToDB(u)
		if err != nil {
			return err
		}
	} else {
		// existing user
		err := updateUserInDB(u)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deletion user from DataBase (soft delete)
func (u *User) Delete() error {
	if u.ID == 0 {
		return appErr.NotFound("user not found")
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
	return getUserByIDFromDB(ID)
}

// Get user by username
func GetUserByUsername(username string) (*User, error) {
	return getUserByUsernameFromDB(username)
}

// Get user by email
func GetUserByEmail(email string) (*User, error) {
	return getUserByEmailFromDB(email)
}

// Get user by phone
func GetUserByPhone(phone string) (*User, error) {
	return getUserByPhoneFromDB(phone)
}

// hashing password
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		logger.GetInstance().Error(err.Error(), "password hashing", "", err)
		return "", appErr.InternalServerError("internal server error")
	}
	return string(hash), nil
}

// compare passwords
func comparePasswords(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
