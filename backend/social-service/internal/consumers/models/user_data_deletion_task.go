package models

import "time"

type UserDataDeletionTask struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	IsConfirmed bool      `json:"is_confirmed"`
	IsDeleted   bool      `json:"is_deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
