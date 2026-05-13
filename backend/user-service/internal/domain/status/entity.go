package status

import (
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
)

type Status struct {
	isConfirmed          bool
	isPermanentlyBanned  bool
	bannedUntil          *time.Time
	isDeleted            bool
	deletedAt            *time.Time
	isPermanentlyDeleted bool
}

func New() *Status {
	return &Status{
		isConfirmed:          false,
		isPermanentlyBanned:  false,
		bannedUntil:          nil,
		isDeleted:            false,
		deletedAt:            nil,
		isPermanentlyDeleted: false,
	}
}

func FromStorage(
	isConfirmed bool,
	isPermanentlyBanned bool,
	bannedUntil *time.Time,
	isDeleted bool,
	deletedAt *time.Time,
	isPermanentlyDeleted bool,
) *Status {
	return &Status{
		isConfirmed:          isConfirmed,
		isPermanentlyBanned:  isPermanentlyBanned,
		bannedUntil:          bannedUntil,
		isDeleted:            isDeleted,
		deletedAt:            deletedAt,
		isPermanentlyDeleted: isPermanentlyDeleted,
	}
}

func (s *Status) IsConfirmed() bool          { return s.isConfirmed }
func (s *Status) IsPermanentlyBanned() bool  { return s.isPermanentlyBanned }
func (s *Status) GetBannedUntil() *time.Time { return s.bannedUntil }
func (s *Status) IsDeleted() bool            { return s.isDeleted }
func (s *Status) GetDeletedAt() *time.Time   { return s.deletedAt }
func (s *Status) IsPermanentlyDeleted() bool { return s.isPermanentlyDeleted }
func (s *Status) IsTemporarilyBanned(t *time.Time) bool {
	if s.bannedUntil == nil {
		return false
	}

	var checkTime time.Time
	if t == nil {
		checkTime = time.Now()
	} else {
		checkTime = *t
	}

	bannedUntil := *s.bannedUntil
	return bannedUntil.After(checkTime)
}

func (s *Status) SetConfirmed(confirmed bool) {
	s.isConfirmed = confirmed
}

func (s *Status) SetPermanentlyBanned(permanentlyBanned bool) {
	s.isPermanentlyBanned = permanentlyBanned
}

func (s *Status) SetBannedUntil(bannedUntil *time.Time) error {
	if bannedUntil != nil && bannedUntil.Before(time.Now()) {
		return ErrBannedUntilInPast
	}
	s.bannedUntil = bannedUntil
	return nil
}

func (s *Status) SetDeleted(deleted bool, now *time.Time) {
	s.isDeleted = deleted
	if now == nil {
		now = utils.TimePtr(time.Now())
	}
	if deleted {
		s.deletedAt = now
	} else {
		s.deletedAt = nil
	}
}

func (s *Status) SetPermanentlyDeleted(permanentlyDeleted bool) error {
	if permanentlyDeleted && !s.isDeleted {
		return ErrPermanentlyDeletedIfNotDeleted
	}
	s.isPermanentlyDeleted = permanentlyDeleted
	return nil
}
