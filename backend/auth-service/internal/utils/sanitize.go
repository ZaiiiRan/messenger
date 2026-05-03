package utils

import (
	"strings"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
)

func SanitizeRegisterRequest(req *pb.RegisterRequest) {
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeConfirmRequest(req *pb.ConfirmRequest) {
	req.Code = strings.TrimSpace(req.Code)
}

func SanitizeLoginRequest(req *pb.LoginRequest) {
	req.Login = strings.TrimSpace(req.Login)
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeChangePasswordRequest(req *pb.ChangePasswordRequest) {
	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
}
