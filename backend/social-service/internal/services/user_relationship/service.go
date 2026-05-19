package userrelationshipservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/redis"
	"go.uber.org/zap"
)

type UserRelationshipService interface {
	AddUserToFriends(ctx context.Context, actor, friendCandidate *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error)
	AddUsersToFriends(ctx context.Context, actor *userpb.User, friendCandidates []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, []*userrelationship.UserRelationship, error)
	RemoveUserFromFriends(ctx context.Context, actor, friend *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error)
	RemoveUsersFromFriends(ctx context.Context, actor *userpb.User, friends []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error)
	BlockUser(ctx context.Context, actor, blockCandidate *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error)
	BlockUsers(ctx context.Context, actor *userpb.User, blockCandidates []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, []*userrelationship.UserRelationship, error)
	UnblockUser(ctx context.Context, actor, unblockCandidate *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, bool, error)
	UnblockUsers(ctx context.Context, actor *userpb.User, unblockCandidates []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, []*userrelationship.UserRelationship, error)
	GetUserRelationship(ctx context.Context, actor, target *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error)
	GetUserRelationships(ctx context.Context, actor *userpb.User, targets []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error)
	GetUserRelationshipsByQuery(ctx context.Context, query *models.QueryUserRelationshipsDal, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error)
}

type service struct {
	log          *zap.SugaredLogger
	dataProvider *userRelationshipDataProvider
}

func New(
	pgClient *postgres.PostgresClient,
	redisClient *redis.RedisClient,
	log *zap.SugaredLogger,
) UserRelationshipService {
	return &service{
		log:          log,
		dataProvider: newUserRelationshipDataProvider(pgClient, redisClient),
	}
}

func (s *service) AddUserToFriends(ctx context.Context, actor, friendCandidate *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "add_user_to_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.add_user_to_friends_failed.begin_transaction_error", "err", err)
			return nil, ErrAddUserToFriends
		}
	}

	existing, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, friendCandidate.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.add_user_to_friends_failed.get_user_relationship_error", "err", err)
		return nil, ErrAddUserToFriends
	}

	ur, blocked, err := applyAddFriend(actor, friendCandidate, existing, now)
	if err != nil {
		if err == ErrAddUserToFriends {
			l.Errorw("user_relationship.add_user_to_friends_failed.set_status_error", "err", err)
		}
		return ur, err
	}

	if err := s.dataProvider.save(ctx, ur, uow); err != nil {
		l.Errorw("user_relationship.add_user_to_friends_failed.save_error", "err", err)
		return nil, ErrAddUserToFriends
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.add_user_to_friends_failed.commit_error", "err", err)
			return nil, ErrAddUserToFriends
		}
	}
	if blocked {
		return ur, ErrBlockedByFriendCandidate
	}

	l.Infow("user_relationship.add_user_to_friends.success")
	return ur, nil
}

func (s *service) RemoveUserFromFriends(ctx context.Context, actor, friend *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "remove_user_from_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.remove_user_from_friends_failed.begin_transaction_error", "err", err)
			return nil, ErrRemoveFromFriends
		}
	}

	existing, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, friend.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.remove_user_from_friends_failed.get_user_relationship_error", "err", err)
		return nil, ErrRemoveFromFriends
	}

	ur, needToDelete := applyRemoveFriend(existing)
	if !needToDelete {
		l.Infow("user_relationship.remove_user_from_friends.success")
		return ur, nil
	}

	if err := s.dataProvider.delete(ctx, ur, uow); err != nil {
		l.Errorw("user_relationship.remove_user_from_friends_failed.delete_error", "err", err)
		return nil, ErrRemoveFromFriends
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.remove_user_from_friends_failed.commit_error", "err", err)
			return nil, ErrRemoveFromFriends
		}
	}

	ur.MarkDeleted()
	l.Infow("user_relationship.remove_user_from_friends.success")
	return ur, nil
}

func (s *service) BlockUser(ctx context.Context, actor, blockCandidate *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "block_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.block_user_failed.begin_transaction_error", "err", err)
			return nil, ErrBlockUser
		}
	}

	existing, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, blockCandidate.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.block_user_failed.get_user_relationship_error", "err", err)
		return nil, ErrBlockUser
	}

	ur, err := applyBlockUser(actor, blockCandidate, existing, now)
	if err != nil {
		return ur, err
	}

	if err := s.dataProvider.save(ctx, ur, uow); err != nil {
		l.Errorw("user_relationship.block_user_failed.save_error", "err", err)
		return nil, ErrBlockUser
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.block_user_failed.commit_error", "err", err)
			return nil, ErrBlockUser
		}
	}

	l.Infow("user_relationship.block_user.success")
	return ur, nil
}

func (s *service) UnblockUser(ctx context.Context, actor, unblockCandidate *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, bool, error) {
	l := s.log.With("op", "unblock_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.begin_transaction_error", "err", err)
			return nil, false, ErrUnblockUser
		}
	}

	existing, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, unblockCandidate.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.unblock_user_failed.get_user_relationship_error", "err", err)
		return nil, false, ErrUnblockUser
	}

	ur, needToDelete, skip := applyUnblockUser(actor, existing, now)
	if skip {
		l.Infow("user_relationship.unblock_user.success")
		return ur, false, nil
	}

	if needToDelete {
		if err := s.dataProvider.delete(ctx, ur, uow); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.delete_error", "err", err)
			return nil, false, ErrUnblockUser
		}
		ur.MarkDeleted()
	} else {
		if err := s.dataProvider.save(ctx, ur, uow); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.save_error", "err", err)
			return nil, false, ErrUnblockUser
		}
	}
	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.commit_error", "err", err)
			return nil, false, ErrUnblockUser
		}
	}

	l.Infow("user_relationship.unblock_user.success")
	return ur, true, nil
}

func (s *service) GetUserRelationship(ctx context.Context, actor, target *userpb.User, uow *uow.UnitOfWork) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "get_user_relationship", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
	}

	ur, err := s.dataProvider.getUserRelationship(ctx, actor.Id, target.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.get_user_relationship_failed.get_user_relationship_error", "err", err)
		return nil, ErrGetUserRelationship
	}

	l.Infow("user_relationship.get_user_relationship.success")
	return ur, nil
}
