package userrelationshipservice

import (
	"context"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	postgresimpl "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/transport/redis"
)

type userRelationshipDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newUserRelationshipDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *userRelationshipDataProvider {
	return &userRelationshipDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (udp *userRelationshipDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(udp.pg)
}

func (udp *userRelationshipDataProvider) getUserRelationship(
	ctx context.Context,
	userID1, userID2 string,
	uow *uow.UnitOfWork,
) (*userrelationship.UserRelationship, error) {
	if userID1 == userID2 {
		return nil, nil
	}

	cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
	ur, err := cacheRepo.GetUserRelationship(ctx, userID1, userID2)
	if err == nil && ur != nil {
		return ur, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	query := models.NewQueryUserRelationshipsDal(&userID1, []string{userID2}, nil, 1, 1, false)
	list, err := dbRepo.QueryUserRelationships(ctx, query, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetUserRelationship(ctx, list[0])

	return list[0], nil
}

func (udp *userRelationshipDataProvider) getUserRelationshipLocked(
	ctx context.Context,
	userID1, userID2 string,
	uow *uow.UnitOfWork,
) (*userrelationship.UserRelationship, error) {
	if userID1 == userID2 {
		return nil, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	query := models.NewQueryUserRelationshipsDal(&userID1, []string{userID2}, nil, 1, 1, false)
	list, err := dbRepo.QueryUserRelationships(ctx, query, true)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list[0], nil
}

func (udp *userRelationshipDataProvider) getUserRelationships(
	ctx context.Context,
	query *models.QueryUserRelationshipsDal,
	uow *uow.UnitOfWork,
) ([]*userrelationship.UserRelationship, error) {
	cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
	list, err := cacheRepo.GetUserRelationshipsList(ctx, query)
	if err == nil && list != nil {
		return list, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	list, err = dbRepo.QueryUserRelationships(ctx, query, false)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetUserRelationshipsList(ctx, query, list)

	return list, nil
}

func (udp *userRelationshipDataProvider) getUserRelationshipsLocked(
	ctx context.Context,
	query *models.QueryUserRelationshipsDal,
	uow *uow.UnitOfWork,
) ([]*userrelationship.UserRelationship, error) {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	list, err := dbRepo.QueryUserRelationships(ctx, query, true)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (udp *userRelationshipDataProvider) createUserRelationships(
	ctx context.Context,
	urs []*userrelationship.UserRelationship,
	uow *uow.UnitOfWork,
) error {
	if len(urs) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	if err := dbRepo.CreateUserRelationships(ctx, urs); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
		toInvalidate := make(map[string]struct{})
		for _, ur := range urs {
			cacheRepo.SetUserRelationship(ctx, ur)
			toInvalidate[ur.GetUserID1()] = struct{}{}
			toInvalidate[ur.GetUserID2()] = struct{}{}
		}
		for uid := range toInvalidate {
			cacheRepo.InvalidateUserRelationshipsLists(ctx, uid)
		}
	})
	return nil
}

func (udp *userRelationshipDataProvider) updateUserRelationships(
	ctx context.Context,
	urs []*userrelationship.UserRelationship,
	uow *uow.UnitOfWork,
) error {
	if len(urs) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	if err := dbRepo.UpdateUserRelationships(ctx, urs); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
		toInvalidate := make(map[string]struct{})
		for _, ur := range urs {
			cacheRepo.SetUserRelationship(ctx, ur)
			toInvalidate[ur.GetUserID1()] = struct{}{}
			toInvalidate[ur.GetUserID2()] = struct{}{}
		}
		for uid := range toInvalidate {
			cacheRepo.InvalidateUserRelationshipsLists(ctx, uid)
		}
	})
	return nil
}

func (udp *userRelationshipDataProvider) deleteUserRelationships(
	ctx context.Context,
	urs []*userrelationship.UserRelationship,
	uow *uow.UnitOfWork,
) error {
	if len(urs) == 0 {
		return nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	if err := dbRepo.DeleteUserRelationships(ctx, urs); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
		toInvalidate := make(map[string]struct{})
		for _, ur := range urs {
			cacheRepo.DelUserRelationship(ctx, ur.GetUserID1(), ur.GetUserID2())
			toInvalidate[ur.GetUserID1()] = struct{}{}
			toInvalidate[ur.GetUserID2()] = struct{}{}
		}
		for uid := range toInvalidate {
			cacheRepo.InvalidateUserRelationshipsLists(ctx, uid)
		}
	})
	return nil
}

func (udp *userRelationshipDataProvider) save(ctx context.Context, ur *userrelationship.UserRelationship, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)

	if !ur.IsPersisted() {
		if err := dbRepo.CreateUserRelationships(ctx, []*userrelationship.UserRelationship{ur}); err != nil {
			return err
		}
	} else {
		if err := dbRepo.UpdateUserRelationships(ctx, []*userrelationship.UserRelationship{ur}); err != nil {
			return err
		}
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
		cacheRepo.SetUserRelationship(ctx, ur)
		cacheRepo.InvalidateUserRelationshipsLists(ctx, ur.GetUserID1())
		cacheRepo.InvalidateUserRelationshipsLists(ctx, ur.GetUserID2())
	})
	return nil
}

func (udp *userRelationshipDataProvider) delete(ctx context.Context, ur *userrelationship.UserRelationship, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRelationshipsRepository(pgConn)
	if err := dbRepo.DeleteUserRelationships(ctx, []*userrelationship.UserRelationship{ur}); err != nil {
		return err
	}

	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserRelationshipsCacheRepository(udp.redis)
		cacheRepo.DelUserRelationship(ctx, ur.GetUserID1(), ur.GetUserID2())
		cacheRepo.InvalidateUserRelationshipsLists(ctx, ur.GetUserID1())
		cacheRepo.InvalidateUserRelationshipsLists(ctx, ur.GetUserID2())
	})
	return nil
}

