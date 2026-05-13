package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	emailchangecode "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code/email_change_code"
)

type V1EmailCodeDal struct {
	Id                int64     `db:"id" json:"id"`
	UserId            string    `db:"user_id" json:"user_id"`
	Email             string    `db:"email" json:"email"`
	Code              string    `db:"code" json:"code"`
	LinkToken         string    `db:"link_token" json:"link_token"`
	GenerationsLeft   int       `db:"generations_left" json:"generations_left"`
	VerificationsLeft int       `db:"verifications_left" json:"verifications_left"`
	ExpiresAt         time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

func V1EmailCodeDalFromDomain(c *emailchangecode.EmailChangeCode) V1EmailCodeDal {
	if c == nil {
		return V1EmailCodeDal{}
	}

	return V1EmailCodeDal{
		Id:                c.GetID(),
		UserId:            c.GetUserID(),
		Email:             c.GetEmail(),
		Code:              c.GetCode(),
		LinkToken:         c.GetLinkToken(),
		GenerationsLeft:   c.GetGenerationsLeft(),
		VerificationsLeft: c.GetVerificationsLeft(),
		ExpiresAt:         c.GetExpiresAt(),
		CreatedAt:         c.GetCreatedAt(),
		UpdatedAt:         c.GetUpdatedAt(),
	}
}

func (c V1EmailCodeDal) IsNull() bool { return false }
func (c V1EmailCodeDal) Index(i int) any {
	switch i {
	case 0:
		return c.Id
	case 1:
		return c.UserId
	case 2:
		return c.Email
	case 3:
		return c.Code
	case 4:
		return c.LinkToken
	case 5:
		return c.GenerationsLeft
	case 6:
		return c.VerificationsLeft
	case 7:
		return c.ExpiresAt
	case 8:
		return c.CreatedAt
	case 9:
		return c.UpdatedAt
	default:
		return nil
	}
}

func (c V1EmailCodeDal) ToDomain() *emailchangecode.EmailChangeCode {
	code := code.FromStorage(
		c.Id, c.UserId,
		c.Code, c.LinkToken,
		code.CodeTypeEmailChange,
		c.GenerationsLeft, c.VerificationsLeft,
		c.ExpiresAt, c.CreatedAt, c.UpdatedAt,
	)
	return emailchangecode.FromStorage(code, c.Email)
}
