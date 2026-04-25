package userservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
	pb "github.com/ZaiiiRan/messenger/backend/user-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/profile"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/status"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type UserService interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error)
	GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error)
	GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error)
	GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error)
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *userDataProvider
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) UserService {
	return &service{
		log:          log,
		dataProvider: newUserDataProvider(pgClient, redisClient),
	}
}

func (s *service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	l := s.log.With("op", "create_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Profile == nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, "user.profile is required")
	}

	birthdate, berr := utils.ParseDatePtr(req.Profile.Birthdate)
	prof, pverr := profile.New(
		req.Profile.FirstName, req.Profile.LastName,
		req.Profile.Phone,
		birthdate,
		req.Profile.Bio,
	)
	st := status.New()
	u, verr := user.New(req.Username, req.Email, prof, st)
	if verr == nil {
		verr = make(validationerror.ValidationError)
	}
	verr.Merge(pverr)
	if berr != nil {
		verr["profile.birthdate"] = berr.Error()
	}
	if len(verr) > 0 {
		l.Warnw("user.create_user_failed.validation_error", "err", verr)
		return nil, verr.ToStatus()
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	byEmail, err := s.dataProvider.getUserByFilter(ctx, models.UserFilterDal{
		Emails: []string{u.GetEmail()},
	}, uow)
	if err != nil {
		l.Errorw("user.create_user_failed.get_by_email_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if byEmail != nil && !isUniqueHolder(byEmail) {
		byEmail = nil
	}

	byUsername, err := s.dataProvider.getUserByFilter(ctx, models.UserFilterDal{
		Usernames: []string{u.GetUsername()},
	}, uow)
	if err != nil {
		l.Errorw("user.create_user_failed.get_by_username_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if byUsername != nil && !isUniqueHolder(byUsername) {
		byUsername = nil
	}

	switch {
	case byEmail == nil && byUsername == nil:

	case byEmail != nil && byUsername != nil &&
		byEmail.GetID() == byUsername.GetID() &&
		!byEmail.GetStatus().IsConfirmed():
		l.Infow("user.create_user.user_already_pending", "user_id", byEmail.GetID())
		return &pb.CreateUserResponse{User: userToProto(byEmail)}, nil

	default:
		vErr := make(validationerror.ValidationError)
		if byEmail != nil {
			vErr["profile.email"] = "user with this email already exists"
		}
		if byUsername != nil {
			vErr["profile.username"] = "user with this username already exists"
		}
		l.Warnw("user.create_user_failed.uniqueness_conflict", "err", vErr)
		return nil, vErr.ToStatus()
	}

	_, err = uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user.create_user_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if err := s.dataProvider.save(ctx, u, uow); err != nil {
		l.Errorw("user.create_user_failed.save_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user.create_user_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}

	l.Infow("user.create_user_success", "user_id", u.GetID())
	return &pb.CreateUserResponse{User: userToProto(u)}, nil
}

func (s *service) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	l := s.log.With("op", "get_user_by_id", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.UserId == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "user_id is required")
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	u, err := s.dataProvider.getByID(ctx, req.UserId, uow)
	if err != nil {
		l.Errorw("user.get_user_by_id_failed.get_by_id_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}

	if u == nil {
		return nil, grpcstatus.Error(codes.NotFound, "user not found")
	}

	l.Infow("user.get_user_by_id_success", "user_id", u.GetID())
	return &pb.GetUserByIDResponse{User: userToProto(u)}, nil
}

func (s *service) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	l := s.log.With("op", "get_user_by_username", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Username == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "username is required")
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	u, err := s.dataProvider.getUserByFilter(ctx, models.UserFilterDal{
		Usernames: []string{req.Username},
	}, uow)
	if err != nil {
		l.Errorw("user.get_user_by_username_failed.get_by_username_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}

	if u == nil {
		return nil, grpcstatus.Error(codes.NotFound, "user not found")
	}

	l.Infow("user.get_user_by_username_success", "user_id", u.GetID())
	return &pb.GetUserByUsernameResponse{User: userToProto(u)}, nil
}

func (s *service) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	l := s.log.With("op", "get_user_by_email", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Email == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "email is required")
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	u, err := s.dataProvider.getUserByFilter(ctx, models.UserFilterDal{
		Emails: []string{req.Email},
	}, uow)
	if err != nil {
		l.Errorw("user.get_user_by_email_failed.get_by_email_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}

	if u == nil {
		return nil, grpcstatus.Error(codes.NotFound, "user not found")
	}

	l.Infow("user.get_user_by_email_success", "user_id", u.GetID())
	return &pb.GetUserByEmailResponse{User: userToProto(u)}, nil
}

func (s *service) ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	l := s.log.With("op", "confirm_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.UserId == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, "user_id is required")
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	u, err := s.dataProvider.getByID(ctx, req.UserId, uow)
	if err != nil {
		l.Errorw("user.confirm_user_failed.get_by_id_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if u == nil {
		return nil, grpcstatus.Error(codes.NotFound, "user not found")
	}

	if u.GetStatus().IsConfirmed() {
		return nil, grpcstatus.Error(codes.FailedPrecondition, "user is already confirmed")
	}
	if u.GetStatus().IsDeleted() {
		return nil, grpcstatus.Error(codes.FailedPrecondition, "deleted user cannot be confirmed")
	}
	bannedUntil := u.GetStatus().GetBannedUntil()
	if u.GetStatus().IsPermanentlyBanned() || (bannedUntil != nil && bannedUntil.After(time.Now())) {
		return nil, grpcstatus.Error(codes.FailedPrecondition, "banned user cannot be confirmed")
	}

	u.GetStatus().SetConfirmed(true)
	u.SetUpdatedAt(utils.TimePtr(time.Now()))

	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("user.confirm_user_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if err := s.dataProvider.save(ctx, u, uow); err != nil {
		l.Errorw("user.confirm_user_failed.save_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user.confirm_user_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, "internal server error")
	}

	l.Infow("user.confirm_user_success", "user_id", u.GetID())
	return &pb.ConfirmUserResponse{User: userToProto(u)}, nil
}

func isUniqueHolder(u *user.User) bool {
	s := u.GetStatus()
	if s.IsPermanentlyBanned() {
		return true
	}
	if bu := s.GetBannedUntil(); bu != nil && bu.After(time.Now()) {
		return true
	}
	if !s.IsDeleted() {
		return true
	}
	if da := s.GetDeletedAt(); da != nil && time.Since(*da) < 30*24*time.Hour {
		return true
	}
	return false
}
