package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
)

type V1CodeDal struct {
	Id                int64     `db:"id" json:"id"`
	UserId            string    `db:"user_id" json:"user_id"`
	Code              string    `db:"code" json:"code"`
	LinkToken         string    `db:"link_token" json:"link_token"`
	GenerationsLeft   int       `db:"generations_left" json:"generations_left"`
	VerificationsLeft int       `db:"verifications_left" json:"verifications_left"`
	ExpiresAt         time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

func V1CodeDalFromDomain(c *code.Code) V1CodeDal {
	if c == nil {
		return V1CodeDal{}
	}

	return V1CodeDal{
		Id:                c.GetID(),
		UserId:            c.GetUserID(),
		Code:              c.GetCode(),
		LinkToken:         c.GetLinkToken(),
		GenerationsLeft:   c.GetGenerationsLeft(),
		VerificationsLeft: c.GetVerificationsLeft(),
		ExpiresAt:         c.GetExpiresAt(),
		CreatedAt:         c.GetCreatedAt(),
		UpdatedAt:         c.GetUpdatedAt(),
	}
}

func (c V1CodeDal) IsNull() bool { return false }
func (c V1CodeDal) Index(i int) any {
	switch i {
	case 0:
		return c.Id
	case 1:
		return c.UserId
	case 2:
		return c.Code
	case 3:
		return c.LinkToken
	case 4:
		return c.GenerationsLeft
	case 5:
		return c.VerificationsLeft
	case 6:
		return c.ExpiresAt
	case 7:
		return c.CreatedAt
	case 8:
		return c.UpdatedAt
	default:
		return nil
	}
}

func (c V1CodeDal) ToDomain(codeType code.CodeType) *code.Code {
	return code.FromStorage(
		c.Id, c.UserId,
		c.Code, c.LinkToken,
		codeType,
		c.GenerationsLeft, c.VerificationsLeft,
		c.ExpiresAt, c.CreatedAt, c.UpdatedAt,
	)
}
