package models

import (
	"time"
)

type CodeMessage struct {
	Id        int64     `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	LinkToken string    `json:"link_token"`
	CodeType  string    `json:"code_type"`
	Language  string    `json:"language"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
