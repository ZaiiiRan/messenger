package utils

import "time"

func BoolPtr(b bool) *bool {
	return &b
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
