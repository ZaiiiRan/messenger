package userrelationshipservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
)

func (s *service) DeleteDeletedUsersRelationships(ctx context.Context, workerID string, userID string, batchSize int, uow *uow.UnitOfWork) (bool, error) {
	l := s.log.With("op", "delete_deleted_users_relationships", "worker_id", workerID, "user_id", userID, "batch_size", batchSize)

	needToCommit := false
	if uow == nil {
		uow = s.dataProvider.newUOW()
		defer uow.Close()
		if _, err := uow.BeginTransaction(ctx); err != nil {
			l.Errorw("user_relationship.delete_deleted_users_relationships_failed.begin_transaction_error", "err", err)
			return false, ErrDeleteDeletedUsersRelationships
		}
		needToCommit = true
	}

	query := models.NewQueryUserRelationshipsDal(
		&userID, nil, nil,
		1, batchSize, false,
	)

	end := false
	urs, err := s.dataProvider.getUserRelationshipsLocked(ctx, query, uow)
	if err != nil {
		l.Errorw("user_relationship.delete_deleted_users_relationships_failed.get_user_relationships_error", "err", err)
		return false, ErrDeleteDeletedUsersRelationships
	}

	if urs == nil {
		l.Infow("user_relationship.delete_deleted_users_relationships.success", "end", true)
		return true, nil
	}
	if len(urs) < batchSize {
		end = true
	}

	if err := s.dataProvider.deleteUserRelationships(ctx, urs, uow); err != nil {
		l.Errorw("user_relationship.delete_deleted_users_relationships_failed.delete_user_relationships_error", "err", err)
		return false, ErrDeleteDeletedUsersRelationships
	}

	if needToCommit {
		if err := uow.Commit(ctx); err != nil {
			l.Errorw("user_relationship.delete_deleted_users_relationships_failed.commit_error", "err", err)
			return false, ErrDeleteDeletedUsersRelationships
		}
	}

	l.Infow("user_relationship.delete_deleted_users_relationships.success", "deleted_count", len(urs), "end", end)
	return end, nil
}
