package userDTO

import (
	"time"
	"backend/internal/models/user"
)

type UserDTO struct {
	ID          uint64     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       *string    `json:"phone"`
	Firstname   string     `json:"firstname"`
	Lastname    string     `json:"lastname"`
	Birthdate   *time.Time `json:"birthdate"`
	IsDeleted   bool       `json:"is_deleted"`
	IsBanned    bool       `json:"is_banned"`
	IsActivated bool       `json:"is_activated"`
}

// Converting user object to user dto
func CreateUserDTOFromUserObj(user *user.User) *UserDTO {
	dto := &UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Birthdate:   user.Birthdate,
		IsDeleted:   user.IsDeleted,
		IsBanned:    user.IsBanned,
		IsActivated: user.IsActivated,
	}
	return dto
}

// Converting user objects to user dtos
func CreateUserDTOsFromUserObjects(users []user.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, user := range users {
		dtos[i] = *CreateUserDTOFromUserObj(&user)
	}
	return dtos
}
