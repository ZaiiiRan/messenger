package userrelationshipservice

import (
	"context"
	"time"

	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	userpb "github.com/ZaiiiRan/messenger/backend/social-service/gen/go/user/v1"
	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
)

func (s *service) AddUsersToFriends(ctx context.Context, actor *userpb.User, friendCandidates []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error) {
	if len(friendCandidates) == 0 {
		return []*userrelationship.UserRelationship{}, nil
	}

	l := s.log.With("op", "add_users_to_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	candidateIDs := make([]string, len(friendCandidates))
	for i, fc := range friendCandidates {
		candidateIDs[i] = fc.Id
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.add_users_to_friends_failed.begin_transaction_error", "err", err)
			return nil, ErrAddUserToFriends
		}
	}

	actorID := actor.Id
	query := models.NewQueryUserRelationshipsDal(&actorID, candidateIDs, nil, 1, len(candidateIDs), false)
	existingList, err := s.dataProvider.getUserRelationshipsLocked(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.add_users_to_friends_failed.get_relationships_error", "err", err)
		return nil, ErrAddUserToFriends
	}

	existingMap := buildRelationshipMap(existingList, actorID)

	results := make([]*userrelationship.UserRelationship, len(friendCandidates))
	var toCreate, toUpdate []*userrelationship.UserRelationship

	for i, fc := range friendCandidates {
		ur, _, err := applyAddFriend(actor, fc, existingMap[fc.Id], now)
		results[i] = ur
		if err != nil {
			continue
		}
		if ur.IsPersisted() {
			toUpdate = append(toUpdate, ur)
		} else {
			toCreate = append(toCreate, ur)
		}
	}

	if len(toCreate) > 0 {
		if err := s.dataProvider.createUserRelationships(ctx, toCreate, uow); err != nil {
			l.Errorw("user_relationship.add_users_to_friends_failed.create_error", "err", err)
			return nil, ErrAddUserToFriends
		}
	}
	if len(toUpdate) > 0 {
		if err := s.dataProvider.updateUserRelationships(ctx, toUpdate, uow); err != nil {
			l.Errorw("user_relationship.add_users_to_friends_failed.update_error", "err", err)
			return nil, ErrAddUserToFriends
		}
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.add_users_to_friends_failed.commit_error", "err", err)
			return nil, ErrAddUserToFriends
		}
	}

	l.Infow("user_relationship.add_users_to_friends.success")
	return results, nil
}

func (s *service) RemoveUsersFromFriends(ctx context.Context, actor *userpb.User, friends []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error) {
	if len(friends) == 0 {
		return []*userrelationship.UserRelationship{}, nil
	}

	l := s.log.With("op", "remove_users_from_friends", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	friendIDs := make([]string, len(friends))
	for i, f := range friends {
		friendIDs[i] = f.Id
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.remove_users_from_friends_failed.begin_transaction_error", "err", err)
			return nil, ErrRemoveFromFriends
		}
	}

	actorID := actor.Id
	query := models.NewQueryUserRelationshipsDal(&actorID, friendIDs, nil, 1, len(friendIDs), false)
	existingList, err := s.dataProvider.getUserRelationshipsLocked(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.remove_users_from_friends_failed.get_relationships_error", "err", err)
		return nil, ErrRemoveFromFriends
	}

	existingMap := buildRelationshipMap(existingList, actorID)

	results := make([]*userrelationship.UserRelationship, len(friends))
	var toDelete []*userrelationship.UserRelationship

	for i, f := range friends {
		if friends[i].Status.IsDeleted {
			continue
		}
		ur, needToDelete := applyRemoveFriend(existingMap[f.Id])
		results[i] = ur
		if needToDelete {
			toDelete = append(toDelete, ur)
		}
	}

	if len(toDelete) > 0 {
		if err := s.dataProvider.deleteUserRelationships(ctx, toDelete, uow); err != nil {
			l.Errorw("user_relationship.remove_users_from_friends_failed.delete_error", "err", err)
			return nil, ErrRemoveFromFriends
		}
		for _, ur := range toDelete {
			ur.MarkDeleted()
		}
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.remove_users_from_friends_failed.commit_error", "err", err)
			return nil, ErrRemoveFromFriends
		}
	}

	l.Infow("user_relationship.remove_users_from_friends.success")
	return results, nil
}

func (s *service) BlockUsers(ctx context.Context, actor *userpb.User, blockCandidates []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error) {
	if len(blockCandidates) == 0 {
		return []*userrelationship.UserRelationship{}, nil
	}

	l := s.log.With("op", "block_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	candidateIDs := make([]string, len(blockCandidates))
	for i, bc := range blockCandidates {
		candidateIDs[i] = bc.Id
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.block_users_failed.begin_transaction_error", "err", err)
			return nil, ErrBlockUser
		}
	}

	actorID := actor.Id
	query := models.NewQueryUserRelationshipsDal(&actorID, candidateIDs, nil, 1, len(candidateIDs), false)
	existingList, err := s.dataProvider.getUserRelationshipsLocked(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.block_users_failed.get_relationships_error", "err", err)
		return nil, ErrBlockUser
	}

	existingMap := buildRelationshipMap(existingList, actorID)

	results := make([]*userrelationship.UserRelationship, len(blockCandidates))
	var toCreate, toUpdate []*userrelationship.UserRelationship

	for i, bc := range blockCandidates {
		ur, err := applyBlockUser(actor, bc, existingMap[bc.Id], now)
		results[i] = ur
		if err != nil {
			continue
		}
		if ur.IsPersisted() {
			toUpdate = append(toUpdate, ur)
		} else {
			toCreate = append(toCreate, ur)
		}
	}

	if len(toCreate) > 0 {
		if err := s.dataProvider.createUserRelationships(ctx, toCreate, uow); err != nil {
			l.Errorw("user_relationship.block_users_failed.create_error", "err", err)
			return nil, ErrBlockUser
		}
	}
	if len(toUpdate) > 0 {
		if err := s.dataProvider.updateUserRelationships(ctx, toUpdate, uow); err != nil {
			l.Errorw("user_relationship.block_users_failed.update_error", "err", err)
			return nil, ErrBlockUser
		}
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.block_users_failed.commit_error", "err", err)
			return nil, ErrBlockUser
		}
	}

	l.Infow("user_relationship.block_users.success")
	return results, nil
}

func (s *service) UnblockUsers(ctx context.Context, actor *userpb.User, unblockCandidates []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error) {
	if len(unblockCandidates) == 0 {
		return []*userrelationship.UserRelationship{}, nil
	}

	l := s.log.With("op", "unblock_users", "req_id", ctxmetadata.GetReqIdFromContext(ctx))
	now := time.Now()

	candidateIDs := make([]string, len(unblockCandidates))
	for i, uc := range unblockCandidates {
		candidateIDs[i] = uc.Id
	}

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		needToCommit = true
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.unblock_users_failed.begin_transaction_error", "err", err)
			return nil, ErrUnblockUser
		}
	}

	actorID := actor.Id
	query := models.NewQueryUserRelationshipsDal(&actorID, candidateIDs, nil, 1, len(candidateIDs), false)
	existingList, err := s.dataProvider.getUserRelationshipsLocked(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.unblock_users_failed.get_relationships_error", "err", err)
		return nil, ErrUnblockUser
	}

	existingMap := buildRelationshipMap(existingList, actorID)

	results := make([]*userrelationship.UserRelationship, len(unblockCandidates))
	var toUpdate []*userrelationship.UserRelationship
	var toDelete []*userrelationship.UserRelationship

	for i, uc := range unblockCandidates {
		ur, needToDelete, skip := applyUnblockUser(actor, existingMap[uc.Id], now)
		results[i] = ur
		if skip {
			continue
		}
		if needToDelete {
			toDelete = append(toDelete, ur)
		} else {
			toUpdate = append(toUpdate, ur)
		}
	}

	if len(toUpdate) > 0 {
		if err := s.dataProvider.updateUserRelationships(ctx, toUpdate, uow); err != nil {
			l.Errorw("user_relationship.unblock_users_failed.update_error", "err", err)
			return nil, ErrUnblockUser
		}
	}
	if len(toDelete) > 0 {
		if err := s.dataProvider.deleteUserRelationships(ctx, toDelete, uow); err != nil {
			l.Errorw("user_relationship.unblock_users_failed.delete_error", "err", err)
			return nil, ErrUnblockUser
		}
		for _, ur := range toDelete {
			ur.MarkDeleted()
		}
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.unblock_users_failed.commit_error", "err", err)
			return nil, ErrUnblockUser
		}
	}

	l.Infow("user_relationship.unblock_users.success")
	return results, nil
}

func (s *service) GetUserRelationships(ctx context.Context, actor *userpb.User, targets []*userpb.User, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error) {
	if len(targets) == 0 {
		return []*userrelationship.UserRelationship{}, nil
	}

	l := s.log.With("op", "get_user_relationships", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	actorID := actor.Id
	targetIDs := make([]string, len(targets))
	for i, t := range targets {
		if t.Id == actorID {
			continue
		}
		targetIDs[i] = t.Id
	}

	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
	}

	query := models.NewQueryUserRelationshipsDal(&actorID, targetIDs, nil, 1, len(targetIDs), false)
	existingList, err := s.dataProvider.getUserRelationships(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.get_user_relationships_failed.get_user_relationships_error", "err", err)
		return nil, ErrGetUserRelationship
	}

	existingMap := buildRelationshipMap(existingList, actorID)
	results := make([]*userrelationship.UserRelationship, len(targets))
	for i, t := range targets {
		ur := existingMap[t.Id]
		results[i] = ur
	}

	l.Infow("user_relationship.get_user_relationships.success")
	return results, nil
}

func (s *service) GetUserRelationshipsByQuery(ctx context.Context, query *models.QueryUserRelationshipsDal, uow *uow.UnitOfWork) ([]*userrelationship.UserRelationship, error) {
	if query == nil {
		return nil, nil
	}

	l := s.log.With("op", "get_user_relationships_by_query", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
	}

	list, err := s.dataProvider.getUserRelationships(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.get_user_relationships_by_query_failed.get_user_relationships_error", "err", err)
		return nil, ErrGetUserRelationship
	}

	l.Infow("user_relationship.get_user_relationships_by_query.success")
	return list, nil
}
