package userservice

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"

	postgresimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/postgres"
	redisimpl "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/impl/redis"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/redis"
)

type userDataProvider struct {
	pg    *postgres.PostgresClient
	redis *redis.RedisClient
}

func newUserDataProvider(pg *postgres.PostgresClient, redis *redis.RedisClient) *userDataProvider {
	return &userDataProvider{
		pg:    pg,
		redis: redis,
	}
}

func (udp *userDataProvider) newUOW() *uow.UnitOfWork {
	return uow.New(udp.pg)
}

func (udp *userDataProvider) getByID(ctx context.Context, id string, uow *uow.UnitOfWork) (*user.User, error) {
	cacheRepo := redisimpl.NewUserCacheRepository(udp.redis)
	u, err := cacheRepo.GetUser(ctx, id)
	if err == nil && u != nil {
		return u, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)
	query := models.NewQueryUsersDal(
		models.UserFilterDal{
			Ids: []string{id},
		},
		models.UserSortDal{},
		1,
		1,
	)
	list, err := dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetUser(ctx, list[0])

	return list[0], nil
}

func (udp *userDataProvider) getUserByFilter(ctx context.Context, filter models.UserFilterDal, uow *uow.UnitOfWork) (*user.User, error) {
	query := models.NewQueryUsersDal(filter, models.UserSortDal{}, 1, 1)

	cacheRepo := redisimpl.NewUserCacheRepository(udp.redis)
	list, err := cacheRepo.GetUserList(ctx, query)
	if err == nil && len(list) > 0 {
		return list[0], nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)
	list, err = dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	cacheRepo.SetUserList(ctx, query, list)

	return list[0], nil
}

func (udp *userDataProvider) getUsersLocked(ctx context.Context, filter models.UserFilterDal, batch_size int, uow *uow.UnitOfWork) ([]*user.User, error) {
	query := models.NewQueryUsersDal(filter, models.UserSortDal{}, 1, batch_size)

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)
	list, err := dbRepo.QueryLocked(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return list, nil
}

func (udp *userDataProvider) deleteUsers(ctx context.Context, users []*user.User, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)
	if err := dbRepo.Delete(ctx, users); err != nil {
		return err
	}

	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.GetID()
	}
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserCacheRepository(udp.redis)
		for _, id := range ids {
			cacheRepo.DeleteUser(ctx, id)
		}
	})
	return nil
}

func (udp *userDataProvider) updateUsers(ctx context.Context, users []*user.User, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)
	if err := dbRepo.Update(ctx, users); err != nil {
		return err
	}

	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.GetID()
	}
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserCacheRepository(udp.redis)
		for _, id := range ids {
			cacheRepo.DeleteUser(ctx, id)
		}
	})
	return nil
}

func (udp *userDataProvider) getUserList(ctx context.Context, query *models.QueryUsersDal, uow *uow.UnitOfWork) ([]*user.User, error) {
	cacheRepo := redisimpl.NewUserCacheRepository(udp.redis)
	list, err := cacheRepo.GetUserList(ctx, query)
	if err == nil && len(list) > 0 {
		return list, nil
	}

	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return nil, err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)
	list, err = dbRepo.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return list, nil
	}

	cacheRepo.SetUserList(ctx, query, list)
	return list, nil
}

func (udp *userDataProvider) save(ctx context.Context, u *user.User, uow *uow.UnitOfWork) error {
	pgConn, err := uow.GetConn(ctx)
	if err != nil {
		return err
	}

	dbRepo := postgresimpl.NewUserRepository(pgConn)

	if u.GetID() == "" {
		if err := dbRepo.Create(ctx, []*user.User{u}); err != nil {
			return err
		}
	} else {
		if err := dbRepo.Update(ctx, []*user.User{u}); err != nil {
			return err
		}
	}

	id := u.GetID()
	email := u.GetEmail()
	username := u.GetUsername()
	uow.OnCommit(func() {
		cacheRepo := redisimpl.NewUserCacheRepository(udp.redis)
		cacheRepo.DeleteUser(ctx, id)
		emailQuery := models.NewQueryUsersDal(models.UserFilterDal{Emails: []string{email}}, models.UserSortDal{}, 1, 1)
		cacheRepo.InvalidateUserList(ctx, emailQuery)
		usernameQuery := models.NewQueryUsersDal(models.UserFilterDal{Usernames: []string{username}}, models.UserSortDal{}, 1, 1)
		cacheRepo.InvalidateUserList(ctx, usernameQuery)
		cacheRepo.SetUser(ctx, u)
		cacheRepo.SetUserList(ctx, emailQuery, []*user.User{u})
		cacheRepo.SetUserList(ctx, usernameQuery, []*user.User{u})
	})

	return nil
}
