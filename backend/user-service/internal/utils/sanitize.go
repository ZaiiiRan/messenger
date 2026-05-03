package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
)

func SanitizeCreateUserRequest(req *pb.CreateUserRequest) {
	if req == nil {
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	sanitizeProfile(req.Profile)
}

func SanitizeGetUsersRequest(req *pb.GetUsersRequest) {
	if req == nil {
		return
	}

	req.Ids = sanitizeStringArray(req.Ids)
	req.FullUsernames = sanitizeLowerStringArray(req.FullUsernames)
	req.PartialUsernames = sanitizeLowerStringArray(req.PartialUsernames)
	req.FullEmails = sanitizeLowerStringArray(req.FullEmails)
	req.PartialEmails = sanitizeLowerStringArray(req.PartialEmails)
	req.PhoneNumbers = sanitizeStringArray(req.PhoneNumbers)
	req.PartialNames = sanitizeStringArray(req.PartialNames)

	req.DeletedFrom = sanitizeStringPtr(req.DeletedFrom)
	req.DeletedTo = sanitizeStringPtr(req.DeletedTo)
	req.CreatedFrom = sanitizeStringPtr(req.CreatedFrom)
	req.CreatedTo = sanitizeStringPtr(req.CreatedTo)
	req.UpdatedFrom = sanitizeStringPtr(req.UpdatedFrom)
	req.UpdatedTo = sanitizeStringPtr(req.UpdatedTo)
}

func SanitizeConfirmUserRequest(req *pb.ConfirmUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeGetUserByIDRequest(req *pb.GetUserByIDRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeGetUserByUsernameRequest(req *pb.GetUserByUsernameRequest) {
	if req == nil {
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
}

func SanitizeGetUserByEmailRequest(req *pb.GetUserByEmailRequest) {
	if req == nil {
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
}

func SanitizeBanUserRequest(req *pb.BanUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
	req.BannedUntil = sanitizeStringPtr(req.BannedUntil)
}

func SanitizeUnbanUserRequest(req *pb.UnbanUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeDeleteUserRequest(req *pb.DeleteUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func sanitizeProfile(p *pb.Profile) {
	if p == nil {
		return
	}
	p.FirstName = strings.TrimSpace(p.FirstName)
	p.LastName = strings.TrimSpace(p.LastName)
	p.Phone = sanitizeStringPtr(p.Phone)
	p.Birthdate = sanitizeStringPtr(p.Birthdate)
	p.Bio = sanitizeStringPtr(p.Bio)
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
