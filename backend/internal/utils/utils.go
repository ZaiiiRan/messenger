package utils

import (
	"time"
	"errors"
)

// Parsing date from string
func ParseDate(date string) (*time.Time, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return nil, errors.New("invalid date format (use DD.MM.YYYY)")
	}
	return &parsedDate, nil
}

// Pointer on string
func StringPtr(s string) *string {
	return &s
}