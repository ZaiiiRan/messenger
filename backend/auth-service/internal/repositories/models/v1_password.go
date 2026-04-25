package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
)

type V1PasswordDal struct {
	Id           int64     `db:"id" json:"id"`
	UserId       string    `db:"user_id" json:"user_id"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func V1PasswordDalFromDomain(p *password.Password) V1PasswordDal {
	if p == nil {
		return V1PasswordDal{}
	}

	return V1PasswordDal{
		Id:           p.GetID(),
		UserId:       p.GetUserID(),
		PasswordHash: p.GetPasswordHash(),
		CreatedAt:    p.GetCreatedAt(),
		UpdatedAt:    p.GetUpdatedAt(),
	}
}

func (p V1PasswordDal) IsNull() bool { return false }
func (p V1PasswordDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.PasswordHash
	case 3:
		return p.CreatedAt
	case 4:
		return p.UpdatedAt
	default:
		return nil
	}
}

func (p V1PasswordDal) ToDomain() *password.Password {
	return password.FromStorage(
		p.Id,
		p.UserId,
		p.PasswordHash,
		p.CreatedAt,
		p.UpdatedAt,
	)
}
