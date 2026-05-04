package authservice

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	codedomain "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	producerinterfaces "github.com/ZaiiiRan/messenger/backend/auth-service/internal/producers/interfaces"
	codeservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/code"
	passwordservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/password"
	tokenservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/token"
	userservice "github.com/ZaiiiRan/messenger/backend/auth-service/internal/services/user_service"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/utils"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	GetNewConfirmationCode(ctx context.Context, req *pb.GetNewConfirmationCodeRequest) (*pb.GetNewConfirmationCodeResponse, error)
	Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error)
	ConfirmByLink(ctx context.Context, req *pb.ConfirmByLinkRequest) (*pb.ConfirmByLinkResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error)
	Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
	ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error)
	ResetPasswordByCode(ctx context.Context, req *pb.ResetPasswordByCodeRequest) (*pb.ResetPasswordByCodeResponse, error)
	ResetPasswordByLink(ctx context.Context, req *pb.ResetPasswordByLinkRequest) (*pb.ResetPasswordByLinkResponse, error)
	GetActiveSessions(ctx context.Context, req *pb.GetActiveSessionsRequest) (*pb.GetActiveSessionsResponse, error)
	InvalidateSessions(ctx context.Context, req *pb.InvalidateSessionsRequest) (*pb.InvalidateSessionsResponse, error)
}

type service struct {
	codeService            codeservice.CodeService
	passwordService        passwordservice.PasswordService
	tokenService           tokenservice.TokenService
	userService            userservice.UserService
	emailCodeTasksProducer producerinterfaces.EmailCodeTasksProducer
	authDataProvider       *authDataProvider
	log                    *zap.SugaredLogger
}

func New(
	codeSvc codeservice.CodeService,
	passwordSvc passwordservice.PasswordService,
	tokenSvc tokenservice.TokenService,
	userSvc userservice.UserService,
	emailCodeTasksProducer producerinterfaces.EmailCodeTasksProducer,
	pgClient *postgres.PostgresClient,
	log *zap.SugaredLogger,
) AuthService {
	return &service{
		codeService:            codeSvc,
		passwordService:        passwordSvc,
		tokenService:           tokenSvc,
		userService:            userSvc,
		emailCodeTasksProducer: emailCodeTasksProducer,
		authDataProvider:       newAuthDataProvider(pgClient),
		log:                    log,
	}
}

func (s *service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	l := s.log.With("op", "register", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	user, err := s.userService.CreateUser(ctx, req.Username, req.Email, req.Profile)
	if err != nil {
		return nil, err
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.register_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	_, err = s.passwordService.CreatePassword(ctx, uow, user, req.Password)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	c, err := s.codeService.GenerateCode(ctx, uow, user.Id, codedomain.CodeTypeActivation)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if !errors.As(err, &cve) {
			return nil, status.Errorf(codes.Internal, "error.internal_server_error")
		}
		c = nil
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.register_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	lang := ctxmetadata.GetLangFromIncomingContext(ctx)

	if c != nil {
		s.emailCodeTasksProducer.ProduceEmailCodeTask(ctx, user.Email, c, lang)
	}

	return &pb.RegisterResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) GetNewConfirmationCode(ctx context.Context, req *pb.GetNewConfirmationCodeRequest) (*pb.GetNewConfirmationCodeResponse, error) {
	l := s.log.With("op", "get_new_confirmation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	user, err := s.getAndCheckUserForConfirmation(ctx)
	if err != nil {
		return nil, err
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.get_new_confirmation_code_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	c, err := s.codeService.GenerateCode(ctx, uow, user.Id, codedomain.CodeTypeActivation)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.get_new_confirmation_code_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	lang := ctxmetadata.GetLangFromIncomingContext(ctx)

	s.emailCodeTasksProducer.ProduceEmailCodeTask(ctx, user.Email, c, lang)

	l.Infow("auth.get_new_confirmation_code.success")
	return &pb.GetNewConfirmationCodeResponse{}, nil
}

func (s *service) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	l := s.log.With("op", "confirm", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	user, err := s.getAndCheckUserForConfirmation(ctx)
	if err != nil {
		return nil, err
	}

	if utf8.RuneCountInString(req.Code) != 6 {
		l.Errorw("auth.confirm_failed", "err", "invalid code")
		return nil, status.Errorf(codes.InvalidArgument, "%s", codedomain.ErrInvalidCode.Error())
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.confirm_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	valid, err := s.codeService.CheckCodeByCode(ctx, uow, user.Id, req.Code, codedomain.CodeTypeActivation)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "%s", codedomain.ErrInvalidCode.Error())
	}

	user, err = s.userService.ConfirmUser(ctx, user.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.confirm_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.confirm.success")
	return &pb.ConfirmResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) ConfirmByLink(ctx context.Context, req *pb.ConfirmByLinkRequest) (*pb.ConfirmByLinkResponse, error) {
	l := s.log.With("op", "confirm_by_link", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s", codedomain.ErrInvalidToken.Error())
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.confirm_by_link_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	userID, valid, err := s.codeService.CheckCodeByLinkToken(ctx, uow, req.Token, codedomain.CodeTypeActivation)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if !valid {
		return nil, status.Errorf(codes.NotFound, "%s", codedomain.ErrInvalidOrExpiredActivationLink.Error())
	}

	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Error(codes.PermissionDenied, userservice.ErrUserIsDeleted.Error())
	}
	if user.Status.IsConfirmed {
		return nil, status.Error(codes.FailedPrecondition, userservice.ErrUserAlreadyActivated.Error())
	}

	user, err = s.userService.ConfirmUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.confirm_by_link_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.confirm_by_link.success")
	return &pb.ConfirmByLinkResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	l := s.log.With("op", "login", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Login == "" {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidLoginOrPassword.Error())
	}
	req.Login = strings.ToLower(req.Login)

	user, err := s.userService.GetUserByUsername(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil {
		user, err = s.userService.GetUserByEmail(ctx, req.Login)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error.internal_server_error")
		}
	}

	if user == nil || user.Status.IsDeleted {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidLoginOrPassword.Error())
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	valid, err := s.passwordService.CheckPassword(ctx, uow, user, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if !valid {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidLoginOrPassword.Error())
	}

	uv, err := s.tokenService.GetUserVersion(ctx, uow, user.Id)
	if err != nil || uv == nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.login_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.login_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.login.success")
	return &pb.LoginResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	l := s.log.With("op", "refresh", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	refreshTokenStr, _ := ctxmetadata.GetRefreshTokenFromIncomingContext(ctx)
	if refreshTokenStr == "" {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	refreshToken, _, err := s.tokenService.ValidateRefreshToken(ctx, uow, refreshTokenStr)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	user, err := s.userService.GetUserByID(ctx, refreshToken.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.refresh_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, nil, refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.refresh_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.refresh.success", "user_id", user.Id)
	return &pb.RefreshResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	l := s.log.With("op", "logout", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	refreshTokenStr, _ := ctxmetadata.GetRefreshTokenFromIncomingContext(ctx)
	if refreshTokenStr == "" {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	if err := s.tokenService.ParseRefreshToken(refreshTokenStr); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	if err := s.tokenService.InvalidateRefreshToken(ctx, uow, refreshTokenStr); err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.logout.success")
	return &pb.LogoutResponse{}, nil
}

func (s *service) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	l := s.log.With("op", "change_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	if req.OldPassword == req.NewPassword {
		return nil, status.Errorf(codes.InvalidArgument, "%s", password.ErrOldAndNewPasswordSame.Error())
	}

	user, err := s.userService.GetUserByID(ctx, claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}
	if user.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(user.Status.BannedUntil) || !user.Status.IsConfirmed {
		return nil, status.Errorf(codes.PermissionDenied, "error.permission_denied")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.change_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	valid, err := s.passwordService.CheckPassword(ctx, uow, user, req.OldPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "%s", password.ErrOldPasswordIncorrect.Error())
	}

	_, err = s.passwordService.UpdatePassword(ctx, uow, user, req.NewPassword)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.change_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.change_password.success")
	return &pb.ChangePasswordResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	l := s.log.With("op", "forgot_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	login := strings.ToLower(req.Login)

	user, err := s.userService.GetUserByUsername(ctx, login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil {
		user, err = s.userService.GetUserByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error.internal_server_error")
		}
	}

	if user == nil || user.Status.IsDeleted {
		l.Infow("auth.forgot_password.user_not_found_or_deleted")
		return &pb.ForgotPasswordResponse{}, nil
	}

	if user.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(user.Status.BannedUntil) || !user.Status.IsConfirmed {
		return &pb.ForgotPasswordResponse{}, nil
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.forgot_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	c, err := s.codeService.GenerateCode(ctx, uow, user.Id, codedomain.CodeTypePasswordReset)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if errors.As(err, &cve) {
			l.Infow("auth.forgot_password.rate_limited")
			return &pb.ForgotPasswordResponse{}, nil
		}
		l.Errorw("auth.forgot_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.forgot_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	lang := ctxmetadata.GetLangFromIncomingContext(ctx)

	s.emailCodeTasksProducer.ProduceEmailCodeTask(ctx, user.Email, c, lang)

	l.Infow("auth.forgot_password.success")
	return &pb.ForgotPasswordResponse{}, nil
}

func (s *service) ResetPasswordByCode(ctx context.Context, req *pb.ResetPasswordByCodeRequest) (*pb.ResetPasswordByCodeResponse, error) {
	l := s.log.With("op", "reset_password_by_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	login := strings.ToLower(req.Login)

	user, err := s.userService.GetUserByUsername(ctx, login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil {
		user, err = s.userService.GetUserByEmail(ctx, login)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "error.internal_server_error")
		}
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
	}
	if user.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(user.Status.BannedUntil) || !user.Status.IsConfirmed {
		return nil, status.Errorf(codes.PermissionDenied, "error.permission_denied")
	}

	if utf8.RuneCountInString(req.Code) != 6 {
		return nil, status.Errorf(codes.InvalidArgument, "%s", codedomain.ErrInvalidCode.Error())
	}

	if err := password.ValidatePassword(req.NewPassword); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.reset_password_by_code_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	valid, err := s.codeService.CheckCodeByCode(ctx, uow, user.Id, req.Code, codedomain.CodeTypePasswordReset)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "%s", codedomain.ErrInvalidCode.Error())
	}

	if _, err := s.passwordService.ForceUpdatePassword(ctx, uow, user, req.NewPassword); err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.reset_password_by_code_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.reset_password_by_code.success")
	return &pb.ResetPasswordByCodeResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) ResetPasswordByLink(ctx context.Context, req *pb.ResetPasswordByLinkRequest) (*pb.ResetPasswordByLinkResponse, error) {
	l := s.log.With("op", "reset_password_by_link", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s", codedomain.ErrInvalidToken.Error())
	}

	if err := password.ValidatePassword(req.NewPassword); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.reset_password_by_link_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	userID, valid, err := s.codeService.CheckCodeByLinkToken(ctx, uow, req.Token, codedomain.CodeTypePasswordReset)
	if err != nil {
		var cve *codedomain.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if !valid {
		return nil, status.Errorf(codes.NotFound, "%s", codedomain.ErrInvalidOrExpiredResetLink.Error())
	}

	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Error(codes.PermissionDenied, userservice.ErrUserIsDeleted.Error())
	}
	if user.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(user.Status.BannedUntil) || !user.Status.IsConfirmed {
		return nil, status.Errorf(codes.PermissionDenied, "error.permission_denied")
	}

	if _, err := s.passwordService.ForceUpdatePassword(ctx, uow, user, req.NewPassword); err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.reset_password_by_link_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.reset_password_by_link.success")
	return &pb.ResetPasswordByLinkResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) GetActiveSessions(ctx context.Context, req *pb.GetActiveSessionsRequest) (*pb.GetActiveSessionsResponse, error) {
	l := s.log.With("op", "get_active_sessions", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	refreshTokenStr, _ := ctxmetadata.GetRefreshTokenFromIncomingContext(ctx)
	if refreshTokenStr == "" {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	refreshToken, uv, err := s.tokenService.ValidateRefreshToken(ctx, uow, refreshTokenStr)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	otherTokens, err := s.tokenService.GetRefreshTokens(ctx, uow, refreshToken, uv, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	sessions := make([]*pb.Session, 0, len(otherTokens))
	for _, t := range otherTokens {
		sessions = append(sessions, toSessionPb(t))
	}

	l.Infow("auth.get_active_sessions.success", "count", len(sessions))
	return &pb.GetActiveSessionsResponse{
		CurrentSession: toSessionPb(refreshToken),
		Sessions:       sessions,
	}, nil
}

func (s *service) InvalidateSessions(ctx context.Context, req *pb.InvalidateSessionsRequest) (*pb.InvalidateSessionsResponse, error) {
	l := s.log.With("op", "invalidate_sessions", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	all := req.All != nil && *req.All

	if !all && len(req.Ids) > 1000 {
		l.Warnw("auth.invalidate_sessions_failed", "err", "to many sessions", "count", len(req.Ids))
		return nil, status.Error(codes.InvalidArgument, ErrTooManySessions.Error())
	}
	if !all && len(req.Ids) == 0 {
		return &pb.InvalidateSessionsResponse{}, nil
	}

	refreshTokenStr, _ := ctxmetadata.GetRefreshTokenFromIncomingContext(ctx)
	if refreshTokenStr == "" {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	refreshToken, uv, err := s.tokenService.ValidateRefreshToken(ctx, uow, refreshTokenStr)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
		}
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.invalidate_sessions_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if !all {
		if err := s.tokenService.InvalidateRefreshTokensByIds(ctx, uow, refreshToken, req.Ids); err != nil {
			return nil, status.Errorf(codes.Internal, "error.internal_server_error")
		}

		if err := uow.Commit(ctx); err != nil {
			l.Errorw("auth.invalidate_sessions_failed", "err", err)
			return nil, status.Errorf(codes.Internal, "error.internal_server_error")
		}
		l.Infow("auth.invalidate_sessions.success", "count", len(req.Ids))
		return &pb.InvalidateSessionsResponse{}, nil
	}

	user, err := s.userService.GetUserByID(ctx, refreshToken.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Error(codes.PermissionDenied, userservice.ErrUserIsDeleted.Error())
	}
	if user.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(user.Status.BannedUntil) || !user.Status.IsConfirmed {
		return nil, status.Errorf(codes.PermissionDenied, "error.permission_denied")
	}

	uv, err = s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.invalidate_sessions_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}

	l.Infow("auth.invalidate_sessions.success")
	return &pb.InvalidateSessionsResponse{
		User:         user,
		AccessToken:  utils.StringPtr(access.GetToken()),
		RefreshToken: utils.StringPtr(refresh.GetToken()),
	}, nil
}

func (s *service) getAndCheckUserForConfirmation(ctx context.Context) (*userpb.User, error) {
	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}
	if claims.IsConfirmed {
		return nil, status.Error(codes.FailedPrecondition, userservice.ErrUserAlreadyActivated.Error())
	}
	if claims.IsDeleted {
		return nil, status.Error(codes.PermissionDenied, userservice.ErrUserIsDeleted.Error())
	}

	user, err := s.userService.GetUserByID(ctx, claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error.internal_server_error")
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "error.unauthorized")
	}
	if user.Status.IsConfirmed {
		return nil, status.Error(codes.FailedPrecondition, userservice.ErrUserAlreadyActivated.Error())
	}
	if user.Status.IsDeleted {
		return nil, status.Error(codes.PermissionDenied, userservice.ErrUserIsDeleted.Error())
	}
	return user, nil
}
