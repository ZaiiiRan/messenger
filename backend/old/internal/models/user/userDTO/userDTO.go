package userDTO

import (
	userModel "backend/internal/models/user"
	"time"
)

type UserDTO struct {
	userModel.BaseUser
	Email     string     `json:"email"`
	Phone     *string    `json:"phone"`
	Birthdate *time.Time `json:"birthdate"`
}

// Converting user object to user dto
func CreateUserDTOFromUserObj(user *userModel.User) *UserDTO {
	dto := &UserDTO{
		BaseUser: userModel.BaseUser{
			ID:          user.ID,
			Username:    user.Username,
			Firstname:   user.Firstname,
			Lastname:    user.Lastname,
			IsDeleted:   user.IsDeleted,
			IsBanned:    user.IsBanned,
			IsActivated: user.IsActivated,
		},
		Email:     user.Email,
		Phone:     user.Phone,
		Birthdate: user.Birthdate,
	}
	return dto
}

// Converting user objects to user dtos
func CreateUserDTOsFromUserObjects(users []userModel.User) []UserDTO {
	dtos := make([]UserDTO, len(users))
	for i, user := range users {
		dtos[i] = *CreateUserDTOFromUserObj(&user)
	}
	return dtos
}
