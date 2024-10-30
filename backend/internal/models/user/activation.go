package user

import (
	"time"
)

type ActivationCode struct {
	ID        uint64
	UserID    uint64
	Code      string
	ExpiresAt time.Time
}