package authservice

import (
	codeservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	userservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_service"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"go.uber.org/zap"
)

type AuthService interface {
}

type service struct {
	codeService      codeservice.CodeService
	passwordService  passwordservice.PasswordService
	tokenService     tokenservice.TokenService
	userService      userservice.UserService
	authDataProvider *authDataProvider
	log              *zap.SugaredLogger
}

func New(
	codeSvc codeservice.CodeService,
	passwordSvc passwordservice.PasswordService,
	tokenSvc tokenservice.TokenService,
	userSvc userservice.UserService,
	pgClient *postgres.PostgresClient,
	log *zap.SugaredLogger,
) AuthService {
	return &service{
		codeService:      codeSvc,
		passwordService:  passwordSvc,
		tokenService:     tokenSvc,
		userService:      userSvc,
		authDataProvider: newAuthDataProvider(pgClient),
		log:              log,
	}
}
