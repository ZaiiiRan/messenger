package userservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	usergrpcclient "github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/client/grpc/user_client"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID string, includePrivacySettings bool) (*userpb.User, error)
	GetUserByUsername(ctx context.Context, username string, includePrivacySettings bool) (*userpb.User, error)
	GetUsers(ctx context.Context, req *userpb.GetUsersRequest) ([]*userpb.User, error)
	UpdateMyPrivacySettingByUser(ctx context.Context, req *userpb.UpdateMyPrivacySettingsByUserRequest) (*userpb.UpdateMyPrivacySettingsByUserResponse, error)
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

func (s *service) GetUserByID(ctx context.Context, userID string, includePrivacySettings bool) (*userpb.User, error) {
	l := s.log.With("op", "get_user_by_id", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetUserByID(ctx, &userpb.GetUserByIDRequest{
		UserId:                 userID,
		IncludePrivacySettings: includePrivacySettings,
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

func (s *service) GetUserByUsername(ctx context.Context, username string, includePrivacySettings bool) (*userpb.User, error) {
	l := s.log.With("op", "get_user_by_username", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetUserByUsername(ctx, &userpb.GetUserByUsernameRequest{
		Username:               username,
		IncludePrivacySettings: includePrivacySettings,
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

func (s *service) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) ([]*userpb.User, error) {
	l := s.log.With("op", "get_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().GetUsers(ctx, req)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("user.get_users_failed.invalid_argument", "err", err)
		} else {
			l.Errorw("user.get_users_failed", "err", err)
		}
		return nil, err
	}

	l.Infow("user.get_users.success", "count", len(resp.Users))
	return resp.Users, nil
}

func (s *service) UpdateMyPrivacySettingByUser(
	ctx context.Context,
	req *userpb.UpdateMyPrivacySettingsByUserRequest,
) (*userpb.UpdateMyPrivacySettingsByUserResponse, error) {
	l := s.log.With("op", "update_my_privacy_settings", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	resp, err := s.userClient.UserClient().UpdateMyPrivacySettingsByUser(ctx, req)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			l.Warnw("user.update_my_privacy_settings_failed.invalid_argument", "err", err)
		} else if status.Code(err) == codes.FailedPrecondition {
			l.Warnw("user.update_my_privacy_settings_failed.failed_precondition", "err", err)
		} else {
			l.Warnw("user.update_my_privacy_settings_failed", "err", err)
		}
		return nil, err
	}

	l.Infow("user.update_my_privacy_settings.success")
	return resp, nil
}
