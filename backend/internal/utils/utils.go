package utils

import (
	"errors"
	"time"
)

// Parsing date from string
func ParseDate(date string) (*time.Time, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return nil, errors.New("invalid date format (use DD.MM.YYYY)")
	}
	return &parsedDate, nil
}

// Pointer in time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Pointer on string
func StringPtr(s string) *string {
	return &s
}
