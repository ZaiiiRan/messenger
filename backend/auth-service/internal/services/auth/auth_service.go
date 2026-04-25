package authservice

import (
	userservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_service"
	"go.uber.org/zap"
)

type AuthService interface {
}

type service struct {
	userService userservice.UserService
	log         *zap.SugaredLogger
}

func New(
	userSvc userservice.UserService,
	log *zap.SugaredLogger,
) AuthService {
	return &service{
		userService: userSvc,
		log:         log,
	}
}
