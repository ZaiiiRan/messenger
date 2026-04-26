package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
)

type V1UserDal struct {
	Id        string    `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func V1UserDalFromDomain(u *user.User) V1UserDal {
	if u == nil {
		return V1UserDal{}
	}

	return V1UserDal{
		Id:        u.GetID(),
		Username:  u.GetUsername(),
		Email:     u.GetEmail(),
		CreatedAt: u.GetCreatedAt(),
		UpdatedAt: u.GetUpdatedAt(),
	}
}

func (u V1UserDal) IsNull() bool { return false }
func (u V1UserDal) Index(i int) any {
	switch i {
	case 0:
		return u.Id
	case 1:
		return u.Username
	case 2:
		return u.Email
	case 3:
		return u.CreatedAt
	case 4:
		return u.UpdatedAt
	default:
		return nil
	}
}

func (u V1UserDal) ToDomain(profile *profile.Profile, status *status.Status) *user.User {
	return user.FromStorage(
		u.Id,
		u.Username,
		u.Email,
		profile,
		status,
		u.CreatedAt,
		u.UpdatedAt,
	)
}
