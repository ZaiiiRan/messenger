package status

import (
	"fmt"
	"time"
)

type Status struct {
	isConfirmed         bool
	isPermanentlyBanned bool
	bannedUntil         *time.Time
	isDeleted           bool
}

func New() *Status {
	return &Status{
		isConfirmed:         false,
		isPermanentlyBanned: false,
		bannedUntil:         nil,
		isDeleted:           false,
	}
}

func FromStorage(
	isConfirmed bool,
	isPermanentlyBanned bool,
	bannedUntil *time.Time,
	isDeleted bool,
) *Status {
	return &Status{
		isConfirmed:         isConfirmed,
		isPermanentlyBanned: isPermanentlyBanned,
		bannedUntil:         bannedUntil,
		isDeleted:           isDeleted,
	}
}

func (s *Status) IsConfirmed() bool          { return s.isConfirmed }
func (s *Status) IsPermanentlyBanned() bool  { return s.isPermanentlyBanned }
func (s *Status) GetBannedUntil() *time.Time { return s.bannedUntil }
func (s *Status) IsDeleted() bool            { return s.isDeleted }

func (s *Status) SetConfirmed(confirmed bool) {
	s.isConfirmed = confirmed
}

func (s *Status) SetPermanentlyBanned(permanentlyBanned bool) {
	s.isPermanentlyBanned = permanentlyBanned
}

func (s *Status) SetBannedUntil(bannedUntil *time.Time) error {
	if bannedUntil != nil && bannedUntil.Before(time.Now()) {
		return fmt.Errorf("banned until time cannot be in the past")
	}
	s.bannedUntil = bannedUntil
	return nil
}

func (s *Status) SetDeleted(deleted bool) {
	s.isDeleted = deleted
}
