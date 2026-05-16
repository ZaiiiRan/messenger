package models

import (
	"time"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
)

type V1UserRelationshipDal struct {
	UserId1   string    `db:"user_id_1" json:"user_id_1"`
	UserId2   string    `db:"user_id_2" json:"user_id_2"`
	Status    int16     `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func V1UserRelationshipDalFromDomain(ur *userrelationship.UserRelationship) V1UserRelationshipDal {
	if ur == nil {
		return V1UserRelationshipDal{}
	}
	return V1UserRelationshipDal{
		UserId1:   ur.GetUserID1(),
		UserId2:   ur.GetUserID2(),
		Status:    int16(ur.GetStatus()),
		CreatedAt: ur.GetCreatedAt(),
		UpdatedAt: ur.GetUpdatedAt(),
	}
}

func (ur V1UserRelationshipDal) IsNull() bool { return false }
func (ur V1UserRelationshipDal) Index(i int) any {
	switch i {
	case 0:
		return ur.UserId1
	case 1:
		return ur.UserId2
	case 2:
		return ur.Status
	case 3:
		return ur.CreatedAt
	case 4:
		return ur.UpdatedAt
	default:
		return nil
	}
}

func (ur V1UserRelationshipDal) ToDomain() *userrelationship.UserRelationship {
	return userrelationship.FromStorage(
		ur.UserId1, ur.UserId2,
		ur.Status,
		ur.CreatedAt, ur.UpdatedAt,
	)
}
