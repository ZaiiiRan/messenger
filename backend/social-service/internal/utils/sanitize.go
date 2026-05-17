package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
)

func SanitizeGetUserByIdRequest(req *pb.GetUserByIdRequest) {
	if req == nil {
		return
	}
	req.Id = strings.TrimSpace(req.Id)
}

func SanitizeGetUsersByIdsRequest(req *pb.GetUsersByIdsRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeGetUserByUsernameRequest(req *pb.GetUserByUsernameRequest) {
	if req == nil {
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
}

func SanitizeGetUsersByUsernamesRequest(req *pb.GetUsersByUsernamesRequest) {
	if req == nil {
		return
	}
	req.Usernames = sanitizeUniqueArray(sanitizeLowerStringArray(req.Usernames))
}

func SanitizeAddUsersToFriendsRequest(req *pb.AddUsersToFriendsRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeRemoveUsersFromFriendsRequest(req *pb.RemoveUsersFromFriendsRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeBlockUsersRequest(req *pb.BlockUsersRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

func SanitizeUnblockUsersRequest(req *pb.UnblockUsersRequest) {
	if req == nil {
		return
	}
	req.Ids = sanitizeUniqueArray(sanitizeStringArray(req.Ids))
}

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
