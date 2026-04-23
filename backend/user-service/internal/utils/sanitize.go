package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
)

func SanitizeCreateUserRequest(req *pb.CreateUserRequest) {
	if req == nil {
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	sanitizeProfile(req.Profile)
}

func SanitizeGetUsersRequest(req *pb.GetUsersRequest) {
	if req == nil {
		return
	}
	if req.BannedUntil != nil {
		trimmed := strings.TrimSpace(*req.BannedUntil)
		req.BannedUntil = &trimmed
	}

	req.Ids = sanitizeStringArray(req.Ids)
	req.FullUsernames = sanitizeStringArray(req.FullUsernames)
	req.PartialUsernames = sanitizeStringArray(req.PartialUsernames)
	req.FullEmails = sanitizeStringArray(req.FullEmails)
	req.PartialEmails = sanitizeStringArray(req.PartialEmails)
	req.PhoneNumbers = sanitizeStringArray(req.PhoneNumbers)
	req.PartialNames = sanitizeStringArray(req.PartialNames)
}

func SanitizeConfirmUserRequest(req *pb.ConfirmUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
}

func SanitizeBanUserRequest(req *pb.BanUserRequest) {
	if req == nil {
		return
	}
	req.UserId = strings.TrimSpace(req.UserId)
	if req.BannedUntil != nil {
		trimmed := strings.TrimSpace(*req.BannedUntil)
		req.BannedUntil = &trimmed
	}
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

	if p.Phone != nil {
		trimmed := strings.TrimSpace(*p.Phone)
		p.Phone = &trimmed
	}
	if p.Birthdate != nil {
		trimmed := strings.TrimSpace(*p.Birthdate)
		p.Birthdate = &trimmed
	}
	if p.Bio != nil {
		trimmed := strings.TrimSpace(*p.Bio)
		p.Bio = &trimmed
	}
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
