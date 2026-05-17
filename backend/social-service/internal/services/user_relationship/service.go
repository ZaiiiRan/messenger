package userrelationshipservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/redis"
	"go.uber.org/zap"
)

type UserRelationshipService interface {
	AddUserToFriends(ctx context.Context, actor, friendCandidate *userpb.User) (*userrelationship.UserRelationship, error)
	RemoveUserFromFriends(ctx context.Context, actor, friend *userpb.User) (*userrelationship.UserRelationship, error)
	BlockUser(ctx context.Context, actor, blockCandidate *userpb.User) (*userrelationship.UserRelationship, error)
	UnblockUser(ctx context.Context, actor, unblockCandidate *userpb.User) (*userrelationship.UserRelationship, error)
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

func (s *service) AddUserToFriends(ctx context.Context, actor, friendCandidate *userpb.User) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "add_user_to_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()
	blocked := false

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_relationship.add_user_to_friends_failed.begin_transaction_error", "err", err)
		return nil, ErrAddUserToFriends
	}

	ur, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, friendCandidate.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.add_user_to_friends_failed.get_user_relationship_error", "err", err)
		return nil, ErrAddUserToFriends
	}
	if ur == nil {
		var friendRequest userrelationship.UserRelationshipStatus
		if actor.Id > friendCandidate.Id {
			friendRequest = userrelationship.FriendRequestBy2
		} else {
			friendRequest = userrelationship.FriendRequestBy1
		}
		ur = userrelationship.New(actor.Id, friendCandidate.Id, friendRequest)
	} else {
		actorRole := ur.RoleOf(actor.Id)
		status := ur.GetStatus()

		var newStatus userrelationship.UserRelationshipStatus
		switch {
		case status == userrelationship.Friends:
			return ur, ErrAlreadyFriends

		case (status == userrelationship.BlockedBy1 && actorRole == 2) ||
			(status == userrelationship.BlockedBy2 && actorRole == 1):
			return ur, ErrBlockedByFriendCandidate

		case (status == userrelationship.FriendRequestBy1 && actorRole == 2) ||
			(status == userrelationship.FriendRequestBy2 && actorRole == 1):
			newStatus = userrelationship.Friends

		case status == userrelationship.FriendRequestBy1 || status == userrelationship.FriendRequestBy2:
			return ur, ErrFriendRequestAlreadySent

		case (status == userrelationship.BlockedBy1 && actorRole == 1) ||
			(status == userrelationship.BlockedBy2 && actorRole == 2):
			if actorRole == 1 {
				newStatus = userrelationship.FriendRequestBy1
			} else {
				newStatus = userrelationship.FriendRequestBy2
			}

		case status == userrelationship.BlockedByBoth:
			blocked = true
			if actorRole == 1 {
				newStatus = userrelationship.BlockedBy2
			} else {
				newStatus = userrelationship.BlockedBy1
			}
		}

		if err := ur.SetStatus(newStatus); err != nil {
			l.Errorw("user_relationship.add_user_to_friends_failed.set_status_error", "err", err)
			return nil, ErrAddUserToFriends
		}
		ur.SetUpdatedAt(&now)
	}

	if err := s.dataProvider.save(ctx, ur, actor.Id, uow); err != nil {
		l.Errorw("user_relationship.add_user_to_friends_failed.save_error", "err", err)
		return nil, ErrAddUserToFriends
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_relationship.add_user_to_friends_failed.commit_error", "err", err)
		return nil, ErrAddUserToFriends
	}
	s.dataProvider.saveCache(ctx, ur)
	s.dataProvider.invalidateLists(ctx, ur, actor.Id)
	if blocked {
		return ur, ErrBlockedByFriendCandidate
	}

	l.Infow("user_relationship.add_user_to_friends.success")
	return ur, nil
}

func (s *service) RemoveUserFromFriends(ctx context.Context, actor, friend *userpb.User) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "remove_user_from_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_relationship.remove_user_from_friends_failed.begin_transaction_error", "err", err)
		return nil, ErrRemoveFromFriends
	}

	ur, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, friend.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.remove_user_from_friends_failed.get_user_relationship_error", "err", err)
		return nil, ErrRemoveFromFriends
	}
	if ur == nil {
		return nil, nil
	}

	status := ur.GetStatus()

	if status == userrelationship.BlockedBy1 || status == userrelationship.BlockedBy2 || status == userrelationship.BlockedByBoth {
		l.Infow("user_relationship.remove_user_from_friends.success")
		return ur, nil
	}

	if err := s.dataProvider.delete(ctx, ur, actor.Id, uow); err != nil {
		l.Errorw("user_relationship.remove_user_from_friends_failed.delete_error", "err", err)
		return nil, ErrRemoveFromFriends
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_relationship.remove_user_from_friends_failed.commit_error", "err", err)
		return nil, ErrRemoveFromFriends
	}
	s.dataProvider.delPairCache(ctx, ur)
	s.dataProvider.invalidateLists(ctx, ur, actor.Id)

	l.Infow("user_relationship.remove_user_from_friends.success")
	return nil, nil
}

func (s *service) BlockUser(ctx context.Context, actor, blockCandidate *userpb.User) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "block_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_relationship.block_user_failed.begin_transaction_error", "err", err)
		return nil, ErrBlockUser
	}

	ur, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, blockCandidate.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.block_user_failed.get_user_relationship_error", "err", err)
		return nil, ErrBlockUser
	}
	if ur == nil {
		var block userrelationship.UserRelationshipStatus
		if actor.Id > blockCandidate.Id {
			block = userrelationship.BlockedBy2
		} else {
			block = userrelationship.BlockedBy1
		}
		ur = userrelationship.New(actor.Id, blockCandidate.Id, block)
	} else {
		actorRole := ur.RoleOf(actor.Id)
		status := ur.GetStatus()

		var newStatus userrelationship.UserRelationshipStatus
		switch {
		case (status == userrelationship.BlockedBy1 && actorRole == 1) ||
			(status == userrelationship.BlockedBy2 && actorRole == 2) || status == userrelationship.BlockedByBoth:
			return ur, ErrAlreadyBlocked

		case (status == userrelationship.BlockedBy1 && actorRole == 2) ||
			(status == userrelationship.BlockedBy2 && actorRole == 1):
			newStatus = userrelationship.BlockedByBoth

		default:
			if actorRole == 1 {
				newStatus = userrelationship.BlockedBy1
			} else {
				newStatus = userrelationship.BlockedBy2
			}
		}

		if err := ur.SetStatus(newStatus); err != nil {
			l.Errorw("user_relationship.block_user_failed.set_status_error", "err", err)
			return nil, ErrBlockUser
		}
		ur.SetUpdatedAt(&now)
	}

	if err := s.dataProvider.save(ctx, ur, actor.Id, uow); err != nil {
		l.Errorw("user_relationship.block_user_failed.save_error", "err", err)
		return nil, ErrBlockUser
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_relationship.block_user_failed.commit_error", "err", err)
		return nil, ErrBlockUser
	}
	s.dataProvider.saveCache(ctx, ur)
	s.dataProvider.invalidateLists(ctx, ur, actor.Id)

	l.Infow("user_relationship.block_user.success")
	return ur, nil
}

func (s *service) UnblockUser(ctx context.Context, actor, unblockCandidate *userpb.User) (*userrelationship.UserRelationship, error) {
	l := s.log.With("op", "unblock_user", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	uow := s.dataProvider.newUOW()
	defer uow.Close()
	_, err := uow.BeginTransaction(ctx)
	if err != nil {
		l.Errorw("user_relationship.unblock_user_failed.begin_transaction_error", "err", err)
		return nil, ErrUnblockUser
	}

	ur, err := s.dataProvider.getUserRelationshipLocked(ctx, actor.Id, unblockCandidate.Id, uow)
	if err != nil {
		l.Errorw("user_relationship.unblock_user_failed.get_user_relationship_error", "err", err)
		return nil, ErrUnblockUser
	}
	if ur == nil {
		return nil, nil
	}

	status := ur.GetStatus()
	actorRole := ur.RoleOf(actor.Id)

	needToDelete := false
	switch {
	case (status == userrelationship.BlockedBy1 && actorRole == 2) ||
		(status == userrelationship.BlockedBy2 && actorRole == 1) ||
		status == userrelationship.FriendRequestBy1 || status == userrelationship.FriendRequestBy2 ||
		status == userrelationship.Friends:
		l.Infow("user_relationship.unblock_user.success")
		return ur, nil

	case (status == userrelationship.BlockedBy1 && actorRole == 1) ||
		(status == userrelationship.BlockedBy2 && actorRole == 2):
		needToDelete = true
	
	default:
		var newStatus userrelationship.UserRelationshipStatus
		if actorRole == 1 {
			newStatus = userrelationship.BlockedBy2
		} else {
			newStatus = userrelationship.BlockedBy1
		}
		
		if err := ur.SetStatus(newStatus); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.set_status_error", "err", err)
			return nil, ErrUnblockUser
		}
		ur.SetUpdatedAt(&now)
	}

	if needToDelete {
		if err := s.dataProvider.delete(ctx, ur, actor.Id, uow); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.delete_error", "err", err)
			return nil, ErrUnblockUser
		}
	} else {
		if err := s.dataProvider.save(ctx, ur, actor.Id, uow); err != nil {
			l.Errorw("user_relationship.unblock_user_failed.save_error", "err", err)
			return nil, ErrUnblockUser
		}
	}
	if err := uow.Commit(ctx); err != nil {
		l.Errorw("user_relationship.unblock_user_failed.commit_error", "err", err)
		return nil, ErrUnblockUser
	}

	if needToDelete {
		s.dataProvider.delPairCache(ctx, ur)
		s.dataProvider.invalidateLists(ctx, ur, actor.Id)
	} else {
		s.dataProvider.saveCache(ctx, ur)
		s.dataProvider.invalidateUserLists(ctx, actor.Id)
	}

	l.Infow("user_relationship.unblock_user.success")
	return ur, nil
}
