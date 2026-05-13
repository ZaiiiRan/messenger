package emailchangecode

import "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"

type EmailChangeCode struct {
	code.Code
	email string
}

func NewEmailChangeCode(code *code.Code, email string) *EmailChangeCode {
	return &EmailChangeCode{
		Code:  *code,
		email: email,
	}
}

func FromStorage(code *code.Code, email string) *EmailChangeCode {
	return &EmailChangeCode{
		Code:  *code,
		email: email,
	}
}

func (c *EmailChangeCode) GetEmail() string { return c.email }

func (c *EmailChangeCode) SetEmail(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	c.email = email
	return nil
}
