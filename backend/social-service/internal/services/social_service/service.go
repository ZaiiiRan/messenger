package socialservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/commonerror"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/errors/validationerror"
	pb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/social/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	userrelationshipservice "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_relationship"
	userrelationshipchangestasks "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_relationship_changes_tasks"
	userservice "github.com/ZaiiiRan/messenger/backend/social-service/internal/services/user_service"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

const maxRelationshipsPerList = 5000
const maxPrivacySettingListSize = 1000

type SocialService interface {
	GetUserByID(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error)
	GetUsersByIDs(ctx context.Context, req *pb.GetUsersByIdsRequest) (*pb.GetUsersByIdsResponse, error)
	GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error)
	GetUsersByUsernames(ctx context.Context, req *pb.GetUsersByUsernamesRequest) (*pb.GetUsersByUsernamesResponse, error)
	AddUsersToFriends(ctx context.Context, req *pb.AddUsersToFriendsRequest) (*pb.AddUsersToFriendsResponse, error)
	RemoveUsersFromFriends(ctx context.Context, req *pb.RemoveUsersFromFriendsRequest) (*pb.RemoveUsersFromFriendsResponse, error)
	BlockUsers(ctx context.Context, req *pb.BlockUsersRequest) (*pb.BlockUsersResponse, error)
	UnblockUsers(ctx context.Context, req *pb.UnblockUsersRequest) (*pb.UnblockUsersResponse, error)
	SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error)
	GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error)
	GetIncomingFriendRequests(ctx context.Context, req *pb.GetIncomingFriendRequestsRequest) (*pb.GetIncomingFriendRequestsResponse, error)
	GetOutgoingFriendRequests(ctx context.Context, req *pb.GetOutgoingFriendRequestsRequest) (*pb.GetOutgoingFriendRequestsResponse, error)
	GetBlockedUsers(ctx context.Context, req *pb.GetBlockedUsersRequest) (*pb.GetBlockedUsersResponse, error)
	UpdateMyPrivacySettings(ctx context.Context, req *pb.UpdateMyPrivacySettingsRequest) (*pb.UpdateMyPrivacySettingsResponse, error)
}

type service struct {
	dataProvider                        *socialDataProvider
	userRelationshipService             userrelationshipservice.UserRelationshipService
	userService                         userservice.UserService
	userRelationshipChangesTasksService userrelationshipchangestasks.UserRelationshipChangesTasksService
	log                                 *zap.SugaredLogger
}

func New(
	pgClient *postgres.PostgresClient,
	userRelationshipService userrelationshipservice.UserRelationshipService,
	userService userservice.UserService,
	userRelationshipChangesTasksService userrelationshipchangestasks.UserRelationshipChangesTasksService,
	log *zap.SugaredLogger,
) SocialService {
	return &service{
		dataProvider:                        newSocialDataProvider(pgClient),
		userRelationshipService:             userRelationshipService,
		userService:                         userService,
		userRelationshipChangesTasksService: userRelationshipChangesTasksService,
		log:                                 log,
	}
}

func (s *service) GetUserByID(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	l := s.log.With("op", "get_user_by_id", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Id == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUserIdIsRequired.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if a.Id == req.Id {
		l.Infow("user.get_user_by_id.success", "user_id", a.GetId())
		return &pb.GetUserByIdResponse{User: toSocialUserProto(a, a, nil, req.IncludePrivacySettings)}, nil
	}

	u, err := s.userService.GetUserByID(ctx, req.Id, true)
	if err != nil {
		return nil, err
	}
	if u == nil || !u.Status.IsConfirmed {
		return nil, grpcstatus.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	ur, err := s.userRelationshipService.GetUserRelationship(ctx, a, u, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	s.processUserWithPrivacySettings(a, u, ur)

	l.Infow("user.get_user_by_id.success", "user_id", u.GetId())
	return &pb.GetUserByIdResponse{User: toSocialUserProto(a, u, ur, req.IncludePrivacySettings)}, nil
}

func (s *service) GetUsersByIDs(ctx context.Context, req *pb.GetUsersByIdsRequest) (*pb.GetUsersByIdsResponse, error) {
	l := s.log.With("op", "get_users_by_ids", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Ids) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUserIdsIsRequired.Error())
	}
	if len(req.Ids) > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyUserIds.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Ids) == 1 && req.Ids[0] == a.Id {
		l.Infow("user.get_users_by_ids.success", "count", 1)
		return &pb.GetUsersByIdsResponse{Users: []*pb.ShortSocialUser{toShortSocialUserProto(a, a, nil, req.IncludePrivacySettings)}}, nil
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:                    req.Ids,
		IsConfirmed:            utils.BoolPtr(true),
		IncludePrivacySettings: true,
		Page:                   1,
		PageSize:               int32(len(req.Ids)),
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		l.Infow("user.get_users_by_ids.success", "count", 0)
		return &pb.GetUsersByIdsResponse{}, nil
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	urs, err := s.userRelationshipService.GetUserRelationships(ctx, a, users, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, user := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(a, user, ur)
		shortSocialUser := toShortSocialUserProto(a, user, ur, req.IncludePrivacySettings)
		result = append(result, shortSocialUser)
	}

	l.Infow("user.get_users_by_ids.success", "count", len(result))
	return &pb.GetUsersByIdsResponse{Users: result}, nil
}

func (s *service) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	l := s.log.With("op", "get_user_by_username", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Username == "" {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUsernameIsRequired.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if a.Username == req.Username {
		l.Infow("user.get_user_by_username.success", "user_id", a.GetId())
		return &pb.GetUserByUsernameResponse{User: toSocialUserProto(a, a, nil, req.IncludePrivacySettings)}, nil
	}

	u, err := s.userService.GetUserByUsername(ctx, req.Username, true)
	if err != nil {
		return nil, err
	}
	if u == nil || !u.Status.IsConfirmed {
		return nil, grpcstatus.Error(codes.NotFound, ErrUserNotFound.Error())
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	ur, err := s.userRelationshipService.GetUserRelationship(ctx, a, u, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	s.processUserWithPrivacySettings(a, u, ur)

	l.Infow("user.get_user_by_username.success", "user_id", u.GetId())
	return &pb.GetUserByUsernameResponse{User: toSocialUserProto(a, u, ur, req.IncludePrivacySettings)}, nil
}

func (s *service) GetUsersByUsernames(ctx context.Context, req *pb.GetUsersByUsernamesRequest) (*pb.GetUsersByUsernamesResponse, error) {
	l := s.log.With("op", "get_users_by_usernames", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Usernames) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUsernamesIsRequired.Error())
	}
	if len(req.Usernames) > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyUsernames.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Usernames) == 1 && req.Usernames[0] == a.Username {
		l.Infow("user.get_users_by_usernames.success", "count", 1)
		return &pb.GetUsersByUsernamesResponse{Users: []*pb.ShortSocialUser{toShortSocialUserProto(a, a, nil, req.IncludePrivacySettings)}}, nil
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		FullUsernames:          req.Usernames,
		IsConfirmed:            utils.BoolPtr(true),
		IsDeleted:              utils.BoolPtr(false),
		IsPermanentlyDeleted:   utils.BoolPtr(false),
		IncludePrivacySettings: true,
		Page:                   1,
		PageSize:               int32(len(req.Usernames)),
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		l.Infow("user.get_users_by_usernames.success", "count", 0)
		return &pb.GetUsersByUsernamesResponse{}, nil
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	urs, err := s.userRelationshipService.GetUserRelationships(ctx, a, users, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, user := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(a, user, ur)
		shortSocialUser := toShortSocialUserProto(a, user, ur, req.IncludePrivacySettings)
		result = append(result, shortSocialUser)
	}

	l.Infow("user.get_users_by_usernames.success", "count", len(result))
	return &pb.GetUsersByUsernamesResponse{Users: result}, nil
}

func (s *service) AddUsersToFriends(ctx context.Context, req *pb.AddUsersToFriendsRequest) (*pb.AddUsersToFriendsResponse, error) {
	l := s.log.With("op", "add_users_to_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Ids) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUserIdsIsRequired.Error())
	}
	if len(req.Ids) > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyUserIds.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Ids) == 1 {
		if req.Ids[0] == a.Id {
			return nil, grpcstatus.Error(codes.InvalidArgument, ErrCannotAddYourselfToFriends.Error())
		}

		u, err := s.userService.GetUserByID(ctx, req.Ids[0], true)
		if err != nil {
			return nil, err
		}
		if u == nil || !u.Status.IsConfirmed {
			return nil, grpcstatus.Error(codes.NotFound, ErrUserNotFound.Error())
		}
		if u.Status.IsDeleted {
			return nil, grpcstatus.Error(codes.InvalidArgument, ErrCannotAddDeletedUserToFriends.Error())
		}

		uow := s.dataProvider.newUOW()
		defer uow.Close()
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("social.add_users_to_friends_failed.begin_transaction_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		ur, err := s.userRelationshipService.AddUserToFriends(ctx, a, u, uow)
		if err != nil {
			if errors.Is(err, userrelationshipservice.ErrAddUserToFriends) {
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			return nil, grpcstatus.Error(codes.InvalidArgument, err.Error())
		}

		if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(
			ctx, []*userrelationship.UserRelationship{ur}, userrelationshipchangestasks.AddToFriends, uow,
		); err != nil {
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		if err := uow.Commit(ctx); err != nil {
			l.Errorw("social.add_users_to_friends_failed.commit_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		s.processUserWithPrivacySettings(a, u, ur)

		l.Infow("user.add_users_to_friends.success")
		return &pb.AddUsersToFriendsResponse{Users: []*pb.ShortSocialUser{toShortSocialUserProto(a, u, ur, false)}}, nil
	}

	addIds := make([]string, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id != a.Id {
			addIds = append(addIds, id)
		}
	}
	if len(addIds) == 0 {
		l.Infow("user.add_users_to_friends.success")
		return &pb.AddUsersToFriendsResponse{}, nil
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:                    addIds,
		IsConfirmed:            utils.BoolPtr(true),
		IncludePrivacySettings: true,
		Page:                   1,
		PageSize:               int32(len(addIds)),
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		return nil, grpcstatus.Error(codes.NotFound, ErrUsersNotFound.Error())
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("social.add_users_to_friends_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	urs, changedUrs, err := s.userRelationshipService.AddUsersToFriends(ctx, a, users, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if len(changedUrs) > 0 {
		if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(ctx, changedUrs, userrelationshipchangestasks.AddToFriends, uow); err != nil {
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("social.add_users_to_friends_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, user := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(a, user, ur)
		shortSocialUser := toShortSocialUserProto(a, user, ur, false)
		result = append(result, shortSocialUser)
	}

	l.Infow("user.add_users_to_friends.success")
	return &pb.AddUsersToFriendsResponse{Users: result}, nil
}

func (s *service) RemoveUsersFromFriends(ctx context.Context, req *pb.RemoveUsersFromFriendsRequest) (*pb.RemoveUsersFromFriendsResponse, error) {
	l := s.log.With("op", "remove_users_from_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Ids) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUserIdsIsRequired.Error())
	}
	if len(req.Ids) > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyUserIds.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Ids) == 1 {
		if req.Ids[0] == a.Id {
			return nil, grpcstatus.Error(codes.InvalidArgument, ErrCannotRemoveYourselfFromFriends.Error())
		}

		u, err := s.userService.GetUserByID(ctx, req.Ids[0], true)
		if err != nil {
			return nil, err
		}
		if u == nil || !u.Status.IsConfirmed {
			return nil, grpcstatus.Error(codes.NotFound, ErrUserNotFound.Error())
		}

		uow := s.dataProvider.newUOW()
		defer uow.Close()
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("social.remove_users_from_friends_failed.begin_transaction_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		ur, err := s.userRelationshipService.RemoveUserFromFriends(ctx, a, u, uow)
		if err != nil {
			if errors.Is(err, userrelationshipservice.ErrRemoveFromFriends) {
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			return nil, grpcstatus.Error(codes.InvalidArgument, err.Error())
		}

		if ur != nil && ur.GetStatus() == userrelationship.None {
			if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(
				ctx, []*userrelationship.UserRelationship{ur}, userrelationshipchangestasks.RemoveFromFriends, uow,
			); err != nil {
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
		}

		if err := uow.Commit(ctx); err != nil {
			l.Errorw("social.remove_users_from_friends_failed.commit_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		s.processUserWithPrivacySettings(a, u, ur)

		l.Infow("user.remove_users_from_friends.success")
		return &pb.RemoveUsersFromFriendsResponse{Users: []*pb.ShortSocialUser{toShortSocialUserProto(a, u, ur, false)}}, nil
	}

	removeIds := make([]string, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id != a.Id {
			removeIds = append(removeIds, id)
		}
	}
	if len(removeIds) == 0 {
		l.Infow("user.remove_users_from_friends.success")
		return &pb.RemoveUsersFromFriendsResponse{}, nil
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:                    removeIds,
		IsConfirmed:            utils.BoolPtr(true),
		IncludePrivacySettings: true,
		Page:                   1,
		PageSize:               int32(len(removeIds)),
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		return nil, grpcstatus.Error(codes.NotFound, ErrUsersNotFound.Error())
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("social.remove_users_from_friends_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	urs, err := s.userRelationshipService.RemoveUsersFromFriends(ctx, a, users, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	deletedUrs := make([]*userrelationship.UserRelationship, 0, len(urs))
	for _, ur := range urs {
		if ur != nil && ur.GetStatus() == userrelationship.None {
			deletedUrs = append(deletedUrs, ur)
		}
	}
	if len(deletedUrs) > 0 {
		if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(ctx, deletedUrs, userrelationshipchangestasks.RemoveFromFriends, uow); err != nil {
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("social.remove_users_from_friends_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, user := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(a, user, ur)
		shortSocialUser := toShortSocialUserProto(a, user, ur, false)
		result = append(result, shortSocialUser)
	}

	l.Infow("user.remove_users_from_friends.success")
	return &pb.RemoveUsersFromFriendsResponse{Users: result}, nil
}

func (s *service) BlockUsers(ctx context.Context, req *pb.BlockUsersRequest) (*pb.BlockUsersResponse, error) {
	l := s.log.With("op", "block_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Ids) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUserIdsIsRequired.Error())
	}
	if len(req.Ids) > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyUserIds.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Ids) == 1 {
		if req.Ids[0] == a.Id {
			return nil, grpcstatus.Error(codes.InvalidArgument, ErrCannotBlockYourself.Error())
		}

		u, err := s.userService.GetUserByID(ctx, req.Ids[0], true)
		if err != nil {
			return nil, err
		}
		if u == nil || !u.Status.IsConfirmed {
			return nil, grpcstatus.Error(codes.NotFound, ErrUserNotFound.Error())
		}

		uow := s.dataProvider.newUOW()
		defer uow.Close()
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("social.block_users_failed.begin_transaction_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		ur, err := s.userRelationshipService.BlockUser(ctx, a, u, uow)
		if err != nil {
			if errors.Is(err, userrelationshipservice.ErrBlockUser) {
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			return nil, grpcstatus.Error(codes.InvalidArgument, err.Error())
		}

		if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(
			ctx, []*userrelationship.UserRelationship{ur}, userrelationshipchangestasks.Block, uow,
		); err != nil {
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		if err := uow.Commit(ctx); err != nil {
			l.Errorw("social.block_users_failed.commit_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		s.processUserWithPrivacySettings(a, u, ur)

		l.Infow("user.block_users.success")
		return &pb.BlockUsersResponse{Users: []*pb.ShortSocialUser{toShortSocialUserProto(a, u, ur, false)}}, nil
	}

	blockIds := make([]string, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id != a.Id {
			blockIds = append(blockIds, id)
		}
	}
	if len(blockIds) == 0 {
		l.Infow("user.block_users.success")
		return &pb.BlockUsersResponse{}, nil
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:                    blockIds,
		IsConfirmed:            utils.BoolPtr(true),
		IncludePrivacySettings: true,
		Page:                   1,
		PageSize:               int32(len(blockIds)),
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		return nil, grpcstatus.Error(codes.NotFound, ErrUsersNotFound.Error())
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("social.block_users_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	urs, changedUrs, err := s.userRelationshipService.BlockUsers(ctx, a, users, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if len(changedUrs) > 0 {
		if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(ctx, changedUrs, userrelationshipchangestasks.Block, uow); err != nil {
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("social.block_users_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, user := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(a, user, ur)
		shortSocialUser := toShortSocialUserProto(a, user, ur, false)
		result = append(result, shortSocialUser)
	}

	l.Infow("user.block_users.success")
	return &pb.BlockUsersResponse{Users: result}, nil
}

func (s *service) UnblockUsers(ctx context.Context, req *pb.UnblockUsersRequest) (*pb.UnblockUsersResponse, error) {
	l := s.log.With("op", "unblock_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Ids) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrUserIdsIsRequired.Error())
	}
	if len(req.Ids) > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrTooManyUserIds.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.Ids) == 1 {
		if req.Ids[0] == a.Id {
			return nil, grpcstatus.Error(codes.InvalidArgument, ErrCannotUnblockYourself.Error())
		}

		u, err := s.userService.GetUserByID(ctx, req.Ids[0], true)
		if err != nil {
			return nil, err
		}
		if u == nil || !u.Status.IsConfirmed {
			return nil, grpcstatus.Error(codes.NotFound, ErrUserNotFound.Error())
		}

		uow := s.dataProvider.newUOW()
		defer uow.Close()
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("social.unblock_users_failed.begin_transaction_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		ur, changed, err := s.userRelationshipService.UnblockUser(ctx, a, u, uow)
		if err != nil {
			if errors.Is(err, userrelationshipservice.ErrUnblockUser) {
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
			return nil, grpcstatus.Error(codes.InvalidArgument, err.Error())
		}

		if changed {
			if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(
				ctx, []*userrelationship.UserRelationship{ur}, userrelationshipchangestasks.Unblock, uow,
			); err != nil {
				return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
			}
		}

		if err := uow.Commit(ctx); err != nil {
			l.Errorw("social.unblock_users_failed.commit_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}

		s.processUserWithPrivacySettings(a, u, ur)

		l.Infow("user.unblock_users.success")
		return &pb.UnblockUsersResponse{Users: []*pb.ShortSocialUser{toShortSocialUserProto(a, u, ur, false)}}, nil
	}

	unblockIds := make([]string, 0, len(req.Ids))
	for _, id := range req.Ids {
		if id != a.Id {
			unblockIds = append(unblockIds, id)
		}
	}
	if len(unblockIds) == 0 {
		l.Infow("user.unblock_users.success")
		return &pb.UnblockUsersResponse{}, nil
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		Ids:                    unblockIds,
		IsConfirmed:            utils.BoolPtr(true),
		IncludePrivacySettings: true,
		Page:                   1,
		PageSize:               int32(len(unblockIds)),
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		return nil, grpcstatus.Error(codes.NotFound, ErrUsersNotFound.Error())
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	if _, err := uow.BeginTransaction(ctx); err != nil {
		l.Errorw("social.unblock_users_failed.begin_transaction_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	urs, changedUrs, err := s.userRelationshipService.UnblockUsers(ctx, a, users, uow)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	if len(changedUrs) > 0 {
		if err := s.userRelationshipChangesTasksService.CreateUserRelationshipChangesTasks(ctx, changedUrs, userrelationshipchangestasks.Unblock, uow); err != nil {
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
	}

	if err := uow.Commit(ctx); err != nil {
		l.Errorw("social.unblock_users_failed.commit_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, user := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(a, user, ur)
		shortSocialUser := toShortSocialUserProto(a, user, ur, false)
		result = append(result, shortSocialUser)
	}

	l.Infow("user.unblock_users.success")
	return &pb.UnblockUsersResponse{Users: result}, nil
}

func (s *service) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	l := s.log.With("op", "search_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Request == nil {
		req.Request = &pb.UsersRequest{}
	}
	if req.Request.SearchFilter == nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrSearchFilterIsRequired.Error())
	}
	if len(*req.Request.SearchFilter) < 3 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrSearchFilterTooShort.Error())
	}
	if len(*req.Request.SearchFilter) > 120 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrSearchFilterTooLong.Error())
	}
	if req.Request.PageSize <= 0 {
		req.Request.PageSize = 50
	}
	if req.Request.Page <= 0 {
		req.Request.Page = 1
	}
	if req.Request.PageSize > 100 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrPagesizeTooLarge.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	users, err := s.userService.GetUsers(ctx, &userpb.GetUsersRequest{
		ExcludeIds:             []string{a.Id},
		SearchFilter:           req.Request.SearchFilter,
		SortByUsername:         true,
		Page:                   req.Request.Page,
		PageSize:               req.Request.PageSize,
		IsConfirmed:            utils.BoolPtr(true),
		IsDeleted:              utils.BoolPtr(false),
		IsPermanentlyDeleted:   utils.BoolPtr(false),
		IsPermanentlyBanned:    utils.BoolPtr(false),
		IncludePrivacySettings: req.Request.IncludePrivacySettings,
	})
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if len(users) == 0 {
		l.Infow("user.search_users.success", "count", 0)
		return &pb.SearchUsersResponse{}, nil
	}

	result, err := s.buildShortSocialUsers(ctx, a, users, req.Request.IncludePrivacySettings)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("user.search_users.success", "count", len(result))
	return &pb.SearchUsersResponse{Users: result}, nil
}

func (s *service) GetFriends(ctx context.Context, req *pb.GetFriendsRequest) (*pb.GetFriendsResponse, error) {
	l := s.log.With("op", "get_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Request == nil {
		req.Request = &pb.UsersRequest{}
	}
	if err := s.validateUsersRequest(req.Request); err != nil {
		return nil, err
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	sortByUpdatedAt := req.Request.SortByUpdatedAt
	users, err := s.fetchUsersFromRelationshipIDs(
		ctx,
		req.Request.Page, req.Request.PageSize,
		req.Request.IncludePrivacySettings,
		req.Request.SearchFilter,
		!sortByUpdatedAt,
		sortByUpdatedAt,
		func(ctx context.Context) ([]string, error) {
			query := models.NewQueryUserRelationshipsDal(
				&a.Id, nil,
				[]userrelationship.UserRelationshipStatus{userrelationship.Friends},
				1, maxRelationshipsPerList, sortByUpdatedAt,
			)

			uow := s.dataProvider.newUOW()
			defer uow.Close()

			list, err := s.userRelationshipService.GetUserRelationshipsByQuery(ctx, query, uow)
			if err != nil {
				return nil, err
			}

			ids := make([]string, 0, len(list))
			for _, ur := range list {
				ids = append(ids, ur.OtherUserID(a.Id))
			}

			return ids, nil
		},
	)
	if err != nil {
		l.Errorw("user.get_friends_failed", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result, err := s.buildShortSocialUsers(ctx, a, users, req.Request.IncludePrivacySettings)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("user.get_friends.success", "count", len(result))
	return &pb.GetFriendsResponse{Users: result}, nil
}

func (s *service) GetIncomingFriendRequests(ctx context.Context, req *pb.GetIncomingFriendRequestsRequest) (*pb.GetIncomingFriendRequestsResponse, error) {
	l := s.log.With("op", "get_incoming_friend_requests", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Request == nil {
		req.Request = &pb.UsersRequest{}
	}
	if err := s.validateUsersRequest(req.Request); err != nil {
		return nil, err
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	sortByUpdatedAt := req.Request.SortByUpdatedAt
	users, err := s.fetchUsersFromRelationshipIDs(
		ctx,
		req.Request.Page, req.Request.PageSize,
		req.Request.IncludePrivacySettings,
		req.Request.SearchFilter,
		!sortByUpdatedAt,
		sortByUpdatedAt,
		func(ctx context.Context) ([]string, error) {
			query := models.NewQueryUserRelationshipsDal(
				&a.Id, nil,
				[]userrelationship.UserRelationshipStatus{userrelationship.FriendRequestBy1, userrelationship.FriendRequestBy2},
				1, maxRelationshipsPerList, sortByUpdatedAt,
			)
			query.DirectionFilter = models.DirectionIncoming

			uow := s.dataProvider.newUOW()
			defer uow.Close()

			list, err := s.userRelationshipService.GetUserRelationshipsByQuery(ctx, query, uow)
			if err != nil {
				return nil, err
			}

			ids := make([]string, 0, len(list))
			for _, ur := range list {
				ids = append(ids, ur.OtherUserID(a.Id))
			}

			return ids, nil
		},
	)
	if err != nil {
		l.Errorw("user.get_incoming_friend_requests_failed", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result, err := s.buildShortSocialUsers(ctx, a, users, req.Request.IncludePrivacySettings)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("user.get_incoming_friend_requests.success", "count", len(result))
	return &pb.GetIncomingFriendRequestsResponse{Users: result}, nil
}

func (s *service) GetOutgoingFriendRequests(ctx context.Context, req *pb.GetOutgoingFriendRequestsRequest) (*pb.GetOutgoingFriendRequestsResponse, error) {
	l := s.log.With("op", "get_outgoing_friend_requests", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Request == nil {
		req.Request = &pb.UsersRequest{}
	}
	if err := s.validateUsersRequest(req.Request); err != nil {
		return nil, err
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	sortByUpdatedAt := req.Request.SortByUpdatedAt
	users, err := s.fetchUsersFromRelationshipIDs(
		ctx,
		req.Request.Page, req.Request.PageSize,
		req.Request.IncludePrivacySettings,
		req.Request.SearchFilter,
		!sortByUpdatedAt,
		sortByUpdatedAt,
		func(ctx context.Context) ([]string, error) {
			query := models.NewQueryUserRelationshipsDal(
				&a.Id, nil,
				[]userrelationship.UserRelationshipStatus{userrelationship.FriendRequestBy1, userrelationship.FriendRequestBy2},
				1, maxRelationshipsPerList, sortByUpdatedAt,
			)
			query.DirectionFilter = models.DirectionOutgoing

			uow := s.dataProvider.newUOW()
			defer uow.Close()

			list, err := s.userRelationshipService.GetUserRelationshipsByQuery(ctx, query, uow)
			if err != nil {
				return nil, err
			}

			ids := make([]string, 0, len(list))
			for _, ur := range list {
				ids = append(ids, ur.OtherUserID(a.Id))
			}

			return ids, nil
		},
	)
	if err != nil {
		l.Errorw("user.get_outgoing_friend_requests_failed", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result, err := s.buildShortSocialUsers(ctx, a, users, req.Request.IncludePrivacySettings)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("user.get_outgoing_friend_requests.success", "count", len(result))
	return &pb.GetOutgoingFriendRequestsResponse{Users: result}, nil
}

func (s *service) GetBlockedUsers(ctx context.Context, req *pb.GetBlockedUsersRequest) (*pb.GetBlockedUsersResponse, error) {
	l := s.log.With("op", "get_blocked_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if req.Request == nil {
		req.Request = &pb.UsersRequest{}
	}
	if err := s.validateUsersRequest(req.Request); err != nil {
		return nil, err
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	sortByUpdatedAt := req.Request.SortByUpdatedAt
	users, err := s.fetchUsersFromRelationshipIDs(
		ctx,
		req.Request.Page, req.Request.PageSize,
		req.Request.IncludePrivacySettings,
		req.Request.SearchFilter,
		!sortByUpdatedAt,
		sortByUpdatedAt,
		func(ctx context.Context) ([]string, error) {
			query := models.NewQueryUserRelationshipsDal(
				&a.Id, nil,
				[]userrelationship.UserRelationshipStatus{userrelationship.BlockedBy1, userrelationship.BlockedBy2, userrelationship.BlockedByBoth},
				1, maxRelationshipsPerList, sortByUpdatedAt,
			)
			query.DirectionFilter = models.DirectionActorBlocked

			uow := s.dataProvider.newUOW()
			defer uow.Close()

			list, err := s.userRelationshipService.GetUserRelationshipsByQuery(ctx, query, uow)
			if err != nil {
				return nil, err
			}

			ids := make([]string, 0, len(list))
			for _, ur := range list {
				ids = append(ids, ur.OtherUserID(a.Id))
			}

			return ids, nil
		},
	)
	if err != nil {
		l.Errorw("user.get_blocked_users_failed", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	result, err := s.buildShortSocialUsers(ctx, a, users, req.Request.IncludePrivacySettings)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("user.get_blocked_users.success", "count", len(result))
	return &pb.GetBlockedUsersResponse{Users: result}, nil
}

func (s *service) fetchUsersFromRelationshipIDs(
	ctx context.Context,
	page, pageSize int32,
	includePrivacySettings bool,
	searchFilter *string,
	sortByUsername bool,
	preserveOrder bool,
	getIDs func(ctx context.Context) ([]string, error),
) ([]*userpb.User, error) {
	allIDs, err := getIDs(ctx)
	if err != nil {
		return nil, err
	}
	if len(allIDs) == 0 {
		return nil, nil
	}

	req := &userpb.GetUsersRequest{
		Ids:                    allIDs,
		IsConfirmed:            utils.BoolPtr(true),
		IncludePrivacySettings: includePrivacySettings,
		SortByUsername:         sortByUsername,
		SortInactiveLast:       true,
		PreserveIdsOrder:       preserveOrder,
		Page:                   page,
		PageSize:               pageSize,
	}
	if searchFilter != nil {
		req.SearchFilter = searchFilter
	}

	return s.userService.GetUsers(ctx, req)
}

type privSettingState struct {
	newValue      *string
	exceptions    []string
	hasExceptions bool
	favourites    []string
	hasFavourites bool
}

func buildPrivacySettingsUpdate(perSetting map[string]*privSettingState) (*userpb.UpdateUserPrivacySettings, []string) {
	out := &userpb.UpdateUserPrivacySettings{}
	fields := make([]string, 0)

	setSetting := func(name string, s *userpb.UpdateUserPrivacySetting) {
		switch name {
		case "avatar":
			out.Avatar = s
		case "photos":
			out.Photos = s
		case "phone_number":
			out.PhoneNumber = s
		case "email":
			out.Email = s
		case "birthdate":
			out.Birthdate = s
		case "online_status":
			out.OnlineStatus = s
		case "first_dialogs_init":
			out.FirstDialogsInit = s
		case "group_chat_invites":
			out.GroupChatInvites = s
		}
	}

	for settingName, st := range perSetting {
		s := &userpb.UpdateUserPrivacySetting{}
		setSetting(settingName, s)

		if st.newValue != nil {
			v := *st.newValue
			s.Value = &v
			fields = append(fields, settingName+".value")
		}

		willClearFavourites := st.newValue != nil && (*st.newValue == "all" || *st.newValue == "friends")
		willClearExceptions := st.newValue != nil && *st.newValue == "none"

		if st.hasFavourites || willClearFavourites {
			if !willClearFavourites {
				s.Favourites = st.favourites
			}
			fields = append(fields, settingName+".favourites")
		}

		if st.hasExceptions || willClearExceptions {
			if !willClearExceptions {
				s.Exceptions = st.exceptions
			}
			fields = append(fields, settingName+".exceptions")
		}
	}

	return out, fields
}

func (s *service) UpdateMyPrivacySettings(ctx context.Context, req *pb.UpdateMyPrivacySettingsRequest) (*pb.UpdateMyPrivacySettingsResponse, error) {
	l := s.log.With("op", "update_my_privacy_settings", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if len(req.Fields) == 0 {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrFieldsAreRequired.Error())
	}
	if req.PrivacySettings == nil {
		return nil, grpcstatus.Error(codes.InvalidArgument, ErrPrivacySettingsAreRequired.Error())
	}

	a, err := s.getAndCheckActor(ctx)
	if err != nil {
		return nil, err
	}

	getUpdateSetting := func(name string) *userpb.UpdateUserPrivacySetting {
		switch name {
		case "avatar":
			return req.PrivacySettings.GetAvatar()
		case "photos":
			return req.PrivacySettings.GetPhotos()
		case "phone_number":
			return req.PrivacySettings.GetPhoneNumber()
		case "email":
			return req.PrivacySettings.GetEmail()
		case "birthdate":
			return req.PrivacySettings.GetBirthdate()
		case "online_status":
			return req.PrivacySettings.GetOnlineStatus()
		case "first_dialogs_init":
			return req.PrivacySettings.GetFirstDialogsInit()
		case "group_chat_invites":
			return req.PrivacySettings.GetGroupChatInvites()
		}
		return nil
	}

	validValues := map[string]bool{"all": true, "friends": true, "none": true}
	validSettingNames := map[string]bool{
		"avatar": true, "photos": true, "phone_number": true, "email": true,
		"birthdate": true, "online_status": true, "first_dialogs_init": true, "group_chat_invites": true,
	}

	perSetting := make(map[string]*privSettingState)
	verr := make(validationerror.ValidationError)

	for _, field := range req.Fields {
		settingName, subField, ok := strings.Cut(field, ".")
		if !ok || !validSettingNames[settingName] {
			continue
		}
		upd := getUpdateSetting(settingName)
		if upd == nil {
			continue
		}
		st, exists := perSetting[settingName]
		if !exists {
			st = &privSettingState{}
			perSetting[settingName] = st
		}
		switch subField {
		case "value":
			v := strings.TrimSpace(upd.GetValue())
			if !validValues[v] {
				verr["privacy_settings."+field] = ErrInvalidPrivacySettingValue.Error()
				continue
			}
			st.newValue = &v
		case "favourites":
			st.hasFavourites = true
			st.favourites = upd.Favourites
		case "exceptions":
			st.hasExceptions = true
			st.exceptions = upd.Exceptions
		}
	}

	if len(verr) > 0 {
		return nil, verr.ToStatus()
	}

	type listEntry struct {
		fieldKey string
		ids      []string
	}
	toValidate := make([]listEntry, 0)
	allUniqueIDs := make(map[string]struct{})

	for settingName, st := range perSetting {
		willClearFavourites := st.newValue != nil && (*st.newValue == "all" || *st.newValue == "friends")
		willClearExceptions := st.newValue != nil && *st.newValue == "none"

		if st.hasFavourites && !willClearFavourites {
			fKey := "privacy_settings." + settingName + ".favourites"
			if len(st.favourites) > maxPrivacySettingListSize {
				verr[fKey] = ErrPrivacyListTooLong.Error()
			} else {
				toValidate = append(toValidate, listEntry{fKey, st.favourites})
				for _, id := range st.favourites {
					allUniqueIDs[id] = struct{}{}
				}
			}
		}

		if st.hasExceptions && !willClearExceptions {
			eKey := "privacy_settings." + settingName + ".exceptions"
			if len(st.exceptions) > maxPrivacySettingListSize {
				verr[eKey] = ErrPrivacyListTooLong.Error()
			} else {
				toValidate = append(toValidate, listEntry{eKey, st.exceptions})
				for _, id := range st.exceptions {
					allUniqueIDs[id] = struct{}{}
				}
			}
		}
	}

	if len(verr) > 0 {
		return nil, verr.ToStatus()
	}

	friendSet := make(map[string]struct{})
	if len(allUniqueIDs) > 0 {
		ids := make([]string, 0, len(allUniqueIDs))
		for id := range allUniqueIDs {
			ids = append(ids, id)
		}
		query := models.NewQueryUserRelationshipsDal(
			&a.Id, ids,
			[]userrelationship.UserRelationshipStatus{userrelationship.Friends},
			1, len(ids), false,
		)
		uow := s.dataProvider.newUOW()
		defer uow.Close()

		urs, err := s.userRelationshipService.GetUserRelationshipsByQuery(ctx, query, uow)
		if err != nil {
			l.Errorw("user.update_my_privacy_settings_failed.get_relationships_error", "err", err)
			return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
		}
		for _, ur := range urs {
			friendSet[ur.OtherUserID(a.Id)] = struct{}{}
		}
	}

	for _, entry := range toValidate {
		for i, id := range entry.ids {
			if _, ok := friendSet[id]; !ok {
				verr[fmt.Sprintf("%s.%d", entry.fieldKey, i)] = ErrUserNotFriend.Error()
			}
		}
	}

	if len(verr) > 0 {
		return nil, verr.ToStatus()
	}

	outSettings, outFields := buildPrivacySettingsUpdate(perSetting)
	if len(outFields) == 0 {
		return nil, grpcstatus.Error(codes.FailedPrecondition, ErrNothingToUpdate.Error())
	}

	resp, err := s.userService.UpdateMyPrivacySettingByUser(ctx, &userpb.UpdateMyPrivacySettingsByUserRequest{
		Fields:          outFields,
		PrivacySettings: outSettings,
	})
	if err != nil {
		if grpcstatus.Code(err) == codes.FailedPrecondition {
			return nil, grpcstatus.Error(codes.FailedPrecondition, ErrNothingToUpdate.Error())
		}
		l.Errorw("user.update_my_privacy_settings_failed.user_service_error", "err", err)
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}

	l.Infow("user.update_my_privacy_settings.success")
	return &pb.UpdateMyPrivacySettingsResponse{PrivacySettings: resp.PrivacySettings}, nil
}

func (s *service) buildShortSocialUsers(
	ctx context.Context,
	actor *userpb.User,
	users []*userpb.User,
	includePrivacySettings bool,
) ([]*pb.ShortSocialUser, error) {
	if len(users) == 0 {
		return []*pb.ShortSocialUser{}, nil
	}

	uow := s.dataProvider.newUOW()
	defer uow.Close()

	urs, err := s.userRelationshipService.GetUserRelationships(ctx, actor, users, uow)
	if err != nil {
		return nil, err
	}

	result := make([]*pb.ShortSocialUser, 0, len(users))
	for i, u := range users {
		ur := urs[i]
		s.processUserWithPrivacySettings(actor, u, ur)
		result = append(result, toShortSocialUserProto(actor, u, ur, includePrivacySettings))
	}
	return result, nil
}

func (s *service) validateUsersRequest(req *pb.UsersRequest) error {
	if req.PageSize <= 0 {
		req.PageSize = 50
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize > 100 {
		return grpcstatus.Error(codes.InvalidArgument, ErrPagesizeTooLarge.Error())
	}
	if req.SearchFilter != nil && *req.SearchFilter != "" {
		if len(*req.SearchFilter) < 3 {
			return grpcstatus.Error(codes.InvalidArgument, ErrSearchFilterTooShort.Error())
		}
		if len(*req.SearchFilter) > 120 {
			return grpcstatus.Error(codes.InvalidArgument, ErrSearchFilterTooLong.Error())
		}
	}
	return nil
}

func (s *service) processUserWithPrivacySettings(actor, user *userpb.User, ur *userrelationship.UserRelationship) {
	if actor.Id == user.Id {
		return
	}

	if user.Status != nil && user.Status.IsDeleted {
		user.Email = ""
		if user.Profile != nil {
			user.Profile.Phone = nil
			user.Profile.Birthdate = nil
			user.Profile.Bio = nil
		}
		return
	}

	if ur != nil && ((ur.GetStatus() == userrelationship.BlockedBy1 && ur.RoleOf(actor.Id) == 2) ||
		(ur.GetStatus() == userrelationship.BlockedBy2 && ur.RoleOf(actor.Id) == 1) ||
		ur.GetStatus() == userrelationship.BlockedByBoth) {
		user.Email = ""
		if user.Profile != nil {
			user.Profile.Phone = nil
			user.Profile.Birthdate = nil
			user.Profile.Bio = nil
		}
		return
	}

	if user.PrivacySettings == nil {
		return
	}

	isFriend := ur != nil && ur.GetStatus() == userrelationship.Friends

	if !privacyAllows(user.PrivacySettings.Email, actor.Id, isFriend) {
		user.Email = ""
	}
	if user.Profile != nil {
		if !privacyAllows(user.PrivacySettings.PhoneNumber, actor.Id, isFriend) {
			user.Profile.Phone = nil
		}
		if !privacyAllows(user.PrivacySettings.Birthdate, actor.Id, isFriend) {
			user.Profile.Birthdate = nil
		}
	}
}

func privacyAllows(setting *userpb.UserPrivacySetting, actorID string, isFriend bool) bool {
	if setting == nil {
		return true
	}
	for _, id := range setting.Exceptions {
		if id == actorID && isFriend {
			return false
		}
	}
	for _, id := range setting.Favourites {
		if id == actorID && isFriend {
			return true
		}
	}
	switch setting.Value {
	case "all":
		return true
	case "friends":
		return isFriend
	case "none":
		return false
	default:
		return false
	}
}

func (s *service) getAndCheckActor(ctx context.Context) (*userpb.User, error) {
	claims, _ := ctxmetadata.GetUserClaimsFromContext(ctx)
	if claims == nil {
		return nil, grpcstatus.Error(codes.Unauthenticated, commonerror.ErrUnauthorized.Error())
	}

	a, err := s.userService.GetUserByID(ctx, claims.Id, false)
	if err != nil {
		return nil, grpcstatus.Error(codes.Internal, commonerror.ErrInternal.Error())
	}
	if a == nil {
		return nil, grpcstatus.Error(codes.Unauthenticated, commonerror.ErrUnauthorized.Error())
	}
	if a.Status == nil || a.Status.IsDeleted {
		return nil, grpcstatus.Error(codes.Unauthenticated, commonerror.ErrUnauthorized.Error())
	}
	if a.Status.IsPermanentlyBanned || utils.IsActiveTemporaryBan(a.Status.BannedUntil) || !a.Status.IsConfirmed {
		return nil, grpcstatus.Error(codes.PermissionDenied, commonerror.ErrPermissionDenied.Error())
	}
	return a, nil
}
