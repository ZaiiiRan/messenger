package models

import (
	"time"

	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
)

type V1UserVersionDal struct {
	Id        int64     `db:"id" json:"id"`
	UserId    string    `db:"user_id" json:"user_id"`
	Version   int       `db:"version" json:"version"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func V1UserVersionDalFromDomain(uv *userversion.UserVersion) V1UserVersionDal {
	if uv == nil {
		return V1UserVersionDal{}
	}

	return V1UserVersionDal{
		Id:        uv.GetID(),
		UserId:    uv.GetUserID(),
		Version:   uv.GetVersion(),
		CreatedAt: uv.GetCreatedAt(),
		UpdatedAt: uv.GetUpdatedAt(),
	}
}

func (p V1UserVersionDal) IsNull() bool { return false }
func (p V1UserVersionDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.Version
	case 3:
		return p.CreatedAt
	case 4:
		return p.UpdatedAt
	default:
		return nil
	}
}

func (p V1UserVersionDal) ToDomain() *userversion.UserVersion {
	return userversion.FromStorage(
		p.Id, p.UserId,
		p.Version,
		p.CreatedAt, p.UpdatedAt,
	)
}
