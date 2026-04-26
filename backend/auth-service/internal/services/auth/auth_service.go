package authservice

import (
	"context"
	"errors"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	producerinteraces "github.com/ZaiiiRan/messenger/backend/auth-service/internal/producers/interfaces"
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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	GetNewConfirmationCode(ctx context.Context, req *pb.GetNewConfirmationCodeRequest) (*pb.GetNewConfirmationCodeResponse, error)
	Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error)
	Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error)
	ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error)
}

type service struct {
	codeService            codeservice.CodeService
	passwordService        passwordservice.PasswordService
	tokenService           tokenservice.TokenService
	userService            userservice.UserService
	emailCodeTasksProducer producerinteraces.EmailCodeTasksProducer
	authDataProvider       *authDataProvider
	log                    *zap.SugaredLogger
}

func New(
	codeSvc codeservice.CodeService,
	passwordSvc passwordservice.PasswordService,
	tokenSvc tokenservice.TokenService,
	userSvc userservice.UserService,
	emailCodeTasksProducer producerinteraces.EmailCodeTasksProducer,
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
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	_, err = s.passwordService.CreatePassword(ctx, uow, user, req.Password)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	c, err := s.codeService.GenerateConfiramtionCode(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.register_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	s.emailCodeTasksProducer.ProduceEmailCodeTask(ctx, user.Email, c)

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
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	c, err := s.codeService.GenerateConfiramtionCode(ctx, uow, user)
	if err != nil {
		var cve *code.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.get_new_confirmation_code_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	s.emailCodeTasksProducer.ProduceEmailCodeTask(ctx, user.Email, c)

	l.Infow("auth.get_new_confirmation_code.success")

	return &pb.GetNewConfirmationCodeResponse{}, nil
}

func (s *service) Confirm(ctx context.Context, req *pb.ConfirmRequest) (*pb.ConfirmResponse, error) {
	l := s.log.With("op", "confirm", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	user, err := s.getAndCheckUserForConfirmation(ctx)
	if err != nil {
		return nil, err
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()
	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.confirm_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	valid, err := s.codeService.CheckConfirmationCode(ctx, uow, user, req.Code)
	if err != nil {
		var cve *code.CodeValidationError
		if errors.As(err, &cve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid code")
	}

	user, err = s.userService.ConfirmUser(ctx, user.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.confirm_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.confirm.success")

	return &pb.ConfirmResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	l := s.log.With("op", "login", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Login == "" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid login or password")
	}
	user, err := s.userService.GetUserByUsername(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if user == nil {
		user, err = s.userService.GetUserByEmail(ctx, req.Login)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal server error")
		}
	}

	if user == nil || user.Status.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "invalid login or password")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	valid, err := s.passwordService.CheckPassword(ctx, uow, user, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid login or password")
	}

	uv, err := s.tokenService.GetUserVersion(ctx, uow, user.Id)
	if err != nil || uv == nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	refreshTokenStr := md.Get("x-refresh-token")
	if len(refreshTokenStr) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	refreshToken, _, err := s.tokenService.ValidateRefreshToken(ctx, uow, refreshTokenStr[0])
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	user, err := s.userService.GetUserByID(ctx, refreshToken.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("auth.refresh_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, nil, refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.refresh_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	refreshTokenStr := md.Get("x-refresh-token")
	if len(refreshTokenStr) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	s.tokenService.InvalidateRefreshToken(ctx, uow, refreshTokenStr[0])

	l.Infow("auth.logout.success")
	return &pb.LogoutResponse{}, nil
}

func (s *service) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	l := s.log.With("op", "change_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil || claims.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if claims.IsPermanentlyBanned || claims.IsTemporarilyBanned || !claims.IsConfirmed {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	user, err := s.userService.GetUserByID(ctx, claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if user == nil || user.Status.IsDeleted {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if user.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(user.Status.BannedUntil) || !user.Status.IsConfirmed {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	if req.OldPassword == req.NewPassword {
		return nil, status.Errorf(codes.InvalidArgument, "old and new passwords are the same")
	}

	uow := s.authDataProvider.newUOW()
	defer uow.Close()

	valid, err := s.passwordService.CheckPassword(ctx, uow, user, req.OldPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if !valid {
		return nil, status.Errorf(codes.InvalidArgument, "old password is incorrect")
	}

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("auth.change_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	_, err = s.passwordService.UpdatePassword(ctx, uow, user, req.NewPassword)
	if err != nil {
		var pve *password.PasswordValidationError
		if errors.As(err, &pve) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	uv, err := s.tokenService.UpdateUserVersion(ctx, uow, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	access, refresh, err := s.tokenService.GenerateToken(ctx, uow, user, uv, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("auth.change_password_failed", "err", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	l.Infow("auth.change_password.success")

	return &pb.ChangePasswordResponse{
		User:         user,
		AccessToken:  access.GetToken(),
		RefreshToken: refresh.GetToken(),
	}, nil
}

func (s *service) getAndCheckUserForConfirmation(ctx context.Context) (*userpb.User, error) {
	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if claims.IsConfirmed {
		return nil, status.Errorf(codes.FailedPrecondition, "user is already activated")
	}
	if claims.IsDeleted {
		return nil, status.Errorf(codes.PermissionDenied, "user is deleted")
	}

	user, err := s.userService.GetUserByID(ctx, claims.Id)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	if claims.IsConfirmed {
		return nil, status.Errorf(codes.FailedPrecondition, "user is already activated")
	}
	if claims.IsDeleted {
		return nil, status.Errorf(codes.PermissionDenied, "user is deleted")
	}
	return user, nil
}
