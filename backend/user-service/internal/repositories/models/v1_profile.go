package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
)

type V1ProfileDal struct {
	Id        int64      `db:"id"`
	UserId    string     `db:"user_id"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	Phone     *string    `db:"phone"`
	Birthdate *time.Time `db:"birthdate"`
	Bio       *string    `db:"bio"`
}

func V1ProfileDalFromDomain(userId string, p *profile.Profile) V1ProfileDal {
	if p == nil {
		return V1ProfileDal{UserId: userId}
	}
	return V1ProfileDal{
		UserId:    userId,
		FirstName: p.GetFirstName(),
		LastName:  p.GetLastName(),
		Phone:     p.GetPhone(),
		Birthdate: p.GetBirthdate(),
		Bio:       p.GetBio(),
	}
}

func (p V1ProfileDal) IsNull() bool { return false }
func (p V1ProfileDal) Index(i int) any {
	switch i {
	case 0:
		return p.Id
	case 1:
		return p.UserId
	case 2:
		return p.FirstName
	case 3:
		return p.LastName
	case 4:
		return p.Phone
	case 5:
		return p.Birthdate
	case 6:
		return p.Bio
	default:
		return nil
	}
}

func (p V1ProfileDal) ToDomain() *profile.Profile {
	return profile.FromStorage(
		p.FirstName,
		p.LastName,
		p.Phone,
		p.Birthdate,
		p.Bio,
	)
}
