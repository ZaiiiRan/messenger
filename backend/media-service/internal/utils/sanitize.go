package utils

import (
	"strings"
)

func sanitizeStringArray(arr []string) []string {
	var sanitized []string
	for _, str := range arr {
		str = strings.TrimSpace(str)
		if str != "" {
			sanitized = append(sanitized, str)
		}
	}
	return sanitized
}

func sanitizeLowerStringArray(arr []string) []string {
	var sanitized []string
	for _, str := range arr {
		str = strings.ToLower(strings.TrimSpace(str))
		if str != "" {
			sanitized = append(sanitized, str)
		}
	}
	return sanitized
}

func sanitizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

func sanitizeLowerStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.ToLower(strings.TrimSpace(*s))
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func sanitizeUniqueArray[T comparable](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	out := s[:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
