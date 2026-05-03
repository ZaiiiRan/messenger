package userservice

import (
	"context"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	usergrpcclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/client/grpc/user_client"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	CreateUser(ctx context.Context, username, email string, profile *pb.Profile) (*pb.User, error)
	ConfirmUser(ctx context.Context, userId string) (*pb.User, error)
	GetUserByID(ctx context.Context, userId string) (*pb.User, error)
	GetUserByUsername(ctx context.Context, username string) (*pb.User, error)
	GetUserByEmail(ctx context.Context, email string) (*pb.User, error)
}

type service struct {
	userClient *usergrpcclient.Client
	log        *zap.SugaredLogger
}

func New(userClient *usergrpcclient.Client, log *zap.SugaredLogger) UserService {
	return &service{
		userClient: userClient,
		log:        log,
	}
}

func (s *service) CreateUser(ctx context.Context, username, email string, profile *pb.Profile) (*pb.User, error) {
	l := s.log.With("op", "create_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().CreateUser(ctx, &pb.CreateUserRequest{
		Username: username,
		Email:    email,
		Profile:  profile,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("user.create_user_failed.invalid_argument", "err", err)
		} else {
			l.Errorw("user.create_user_failed.create_user_error", "err", err)
		}
		return nil, err
	}

	l.Infow("user.create_user.success", "user_id", resp.User.Id)
	return resp.User, nil
}

func (s *service) ConfirmUser(ctx context.Context, userId string) (*pb.User, error) {
	l := s.log.With("op", "confirm_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().ConfirmUser(ctx, &pb.ConfirmUserRequest{
		UserId: userId,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument || status.Code(err) == codes.NotFound ||
			status.Code(err) == codes.FailedPrecondition {
			l.Warnw("user.confirm_user_failed.invalid_argument", "err", err)
		} else {
			l.Errorw("user.confirm_user_failed.confirm_user_error", "err", err)
		}
		return nil, err
	}

	l.Infow("user.confirm_user.success", "user_id", resp.User.Id)
	return resp.User, nil
}

func (s *service) GetUserByID(ctx context.Context, userId string) (*pb.User, error) {
	l := s.log.With("op", "get_user_by_id", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetUserByID(ctx, &pb.GetUserByIDRequest{
		UserId: userId,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("user.get_user_by_id_failed.invalid_argument", "err", err)
		} else if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_user_by_id_failed.not_found", "err", err)
			return nil, nil
		} else {
			l.Errorw("user.get_user_by_id_failed.get_user_by_id_error", "err", err)
		}
		return nil, err
	}

	l.Infow("user.get_user_by_id.success", "user_id", resp.User.Id)
	return resp.User, nil
}

func (s *service) GetUserByUsername(ctx context.Context, username string) (*pb.User, error) {
	l := s.log.With("op", "get_user_by_username", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("user.get_user_by_username_failed.invalid_argument", "err", err)
		} else if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_user_by_username_failed.not_found", "err", err)
			return nil, nil
		} else {
			l.Errorw("user.get_user_by_username_failed.get_user_by_username_error", "err", err)
		}
		return nil, err
	}

	l.Infow("user.get_user_by_username.success", "user_id", resp.User.Id)
	return resp.User, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	l := s.log.With("op", "get_user_by_email", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("user.get_user_by_email_failed.invalid_argument", "err", err)
		} else if status.Code(err) == codes.NotFound {
			l.Warnw("user.get_user_by_email_failed.not_found", "err", err)
			return nil, nil
		} else {
			l.Errorw("user.get_user_by_email_failed.get_user_by_email_error", "err", err)
		}
		return nil, err
	}

	l.Infow("user.get_user_by_email.success", "user_id", resp.User.Id)
	return resp.User, nil
}
