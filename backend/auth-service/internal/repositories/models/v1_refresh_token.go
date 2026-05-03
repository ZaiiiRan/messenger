package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
)

type V1RefreshTokenDal struct {
	Id        int64     `db:"id" json:"id"`
	UserId    string    `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	Version   int       `db:"version" json:"version"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func V1RefreshTokenFromDomain(t *token.Token) V1RefreshTokenDal {
	if t == nil {
		return V1RefreshTokenDal{}
	}

	return V1RefreshTokenDal{
		Id:        t.GetID(),
		UserId:    t.GetUserID(),
		Token:     t.GetToken(),
		Version:   t.GetVersion(),
		ExpiresAt: t.GetExpiresAt(),
		CreatedAt: t.GetCreatedAt(),
		UpdatedAt: t.GetUpdatedAt(),
	}
}

func (p V1RefreshTokenDal) IsNull() bool { return false }
func (p V1RefreshTokenDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.Token
	case 3:
		return p.Version
	case 4:
		return p.ExpiresAt
	case 5:
		return p.CreatedAt
	case 6:
		return p.UpdatedAt
	default:
		return nil
	}
}

func (p V1RefreshTokenDal) ToDomain() *token.Token {
	return token.FromStorage(
		p.Id,
		p.UserId,
		p.Token,
		token.RefreshTokenType,
		p.Version,
		p.ExpiresAt,
		p.CreatedAt,
		p.UpdatedAt,
	)
}
