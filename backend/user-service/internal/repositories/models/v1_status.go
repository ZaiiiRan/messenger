package models

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
)

type V1StatusDal struct {
	Id                  int64      `db:"id" json:"id"`
	UserId              string     `db:"user_id" json:"user_id"`
	IsConfirmed         bool       `db:"is_confirmed" json:"is_confirmed"`
	IsPermanentlyBanned bool       `db:"is_permanently_banned" json:"is_permanently_banned"`
	BannedUntil         *time.Time `db:"banned_until" json:"banned_until"`
	IsDeleted           bool       `db:"is_deleted" json:"is_deleted"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deleted_at"`
}

func V1StatusDalFromDomain(userId string, s *status.Status) V1StatusDal {
	if s == nil {
		return V1StatusDal{UserId: userId}
	}
	return V1StatusDal{
		UserId:              userId,
		IsConfirmed:         s.IsConfirmed(),
		IsPermanentlyBanned: s.IsPermanentlyBanned(),
		BannedUntil:         s.GetBannedUntil(),
		IsDeleted:           s.IsDeleted(),
		DeletedAt:           s.GetDeletedAt(),
	}
}

func (s V1StatusDal) IsNull() bool { return false }
func (s V1StatusDal) Index(i int) any {
	switch i {
	case 0:
		return s.Id
	case 1:
		return s.UserId
	case 2:
		return s.IsConfirmed
	case 3:
		return s.IsPermanentlyBanned
	case 4:
		return s.BannedUntil
	case 5:
		return s.IsDeleted
	case 6:
		return s.DeletedAt
	default:
		return nil
	}
}

func (s V1StatusDal) ToDomain() *status.Status {
	return status.FromStorage(
		s.IsConfirmed,
		s.IsPermanentlyBanned,
		s.BannedUntil,
		s.IsDeleted,
		s.DeletedAt,
	)
}
