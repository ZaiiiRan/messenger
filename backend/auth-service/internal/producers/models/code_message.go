package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
)

type CodeMessage struct {
	Id        int64     `json:"id"`
	UserId    string    `json:"user_id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	LinkToken string    `json:"link_token"`
	CodeType  string    `json:"code_type"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CodeMessageFromDomain(c *code.Code, email string) CodeMessage {
	if c == nil {
		return CodeMessage{}
	}

	return CodeMessage{
		Id:        c.GetID(),
		UserId:    c.GetUserID(),
		Email:     email,
		Code:      c.GetCode(),
		LinkToken: c.GetLinkToken(),
		CodeType:  string(c.GetCodeType()),
		ExpiresAt: c.GetExpiresAt(),
		CreatedAt: c.GetCreatedAt(),
		UpdatedAt: c.GetUpdatedAt(),
	}
}
