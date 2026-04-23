package utils

import "time"

const dateLayout = "2006-01-02"

func ParseDate(s string) (time.Time, error) {
	return time.Parse(dateLayout, s)
}

func ParseDatePtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(dateLayout, *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func FormatDate(t time.Time) string {
	return t.Format(dateLayout)
}

func FormatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(dateLayout)
	return &s
}
