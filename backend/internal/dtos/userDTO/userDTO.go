package userDTO

import (
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