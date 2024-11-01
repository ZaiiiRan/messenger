package userDTO

import (
	"backend/internal/models/user"
	"time"
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

func CreateUserDTOFromUserObj(user *user.User) (*UserDTO) {
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
