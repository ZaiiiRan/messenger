package passwordservice

import (
	"context"
	"errors"
	"fmt"

	userpb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	passworddomain "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"go.uber.org/zap"
)

type PasswordService interface {
	CreatePassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (*passworddomain.Password, error)
	CheckPassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (bool, error)
	UpdatePassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (*passworddomain.Password, error)
	ForceUpdatePassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (*passworddomain.Password, error)
}

type passwordService struct {
	passwordDataProvider *passwordDataProvider
	log                  *zap.SugaredLogger
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) PasswordService {
	return &passwordService{
		passwordDataProvider: newPasswordDataProvider(pgClient, redisClient),
		log:                  log,
	}
}

func (s *passwordService) CreatePassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (*passworddomain.Password, error) {
	l := s.log.With("op", "create_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	p, err := passworddomain.New(user.Id, rawPassword)
	if err != nil {
		var valErr *passworddomain.PasswordValidationError
		if errors.As(err, &valErr) {
			l.Warnw("password.create_password_failed.validation_error", "err", err)
		} else {
			l.Warnw("password.create_password_failed.validation_error", "err", err)
		}
		return nil, err
	}

	existed, err := s.passwordDataProvider.getByUserID(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("password.create_password_failed", "err", err)
		return nil, err
	}
	if existed != nil {
		return p, nil
	}

	if err := s.passwordDataProvider.save(ctx, p, uow); err != nil {
		l.Errorw("password.create_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.create_password.success")
	return p, nil
}

func (s *passwordService) CheckPassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (bool, error) {
	l := s.log.With("op", "check_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	p, err := s.passwordDataProvider.getByUserID(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("password.check_password_failed", "err", err)
		return false, err
	}
	if p == nil {
		l.Errorw("password.check_password_failed", "err", "password not found")
		return false, fmt.Errorf("password not found")
	}

	correct := p.CheckPassword(rawPassword)
	l.Infow("password.check_password.success")
	return correct, nil
}

func (s *passwordService) UpdatePassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (*passworddomain.Password, error) {
	l := s.log.With("op", "update_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	p, err := s.passwordDataProvider.getByUserID(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("password.update_password_failed", "err", err)
		return nil, err
	}
	if p == nil {
		l.Errorw("password.update_password_failed", "err", "password not found")
		return nil, fmt.Errorf("password not found")
	}

	if err := p.SetPassword(rawPassword); err != nil {
		var valErr *passworddomain.PasswordValidationError
		if errors.As(err, &valErr) {
			l.Warnw("password.update_password_failed.validation_error", "err", err)
		} else {
			l.Errorw("password.update_password_failed", "err", err)
		}
		return nil, err
	}

	if err := s.passwordDataProvider.save(ctx, p, uow); err != nil {
		l.Errorw("password.update_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.update_password.success")
	return p, nil
}

func (s *passwordService) ForceUpdatePassword(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawPassword string) (*passworddomain.Password, error) {
	l := s.log.With("op", "force_update_password", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	p, err := s.passwordDataProvider.getByUserID(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("password.force_update_password_failed", "err", err)
		return nil, err
	}
	if p == nil {
		l.Errorw("password.force_update_password_failed", "err", "password not found")
		return nil, fmt.Errorf("password not found")
	}

	if err := p.ForceSetPassword(rawPassword); err != nil {
		var valErr *passworddomain.PasswordValidationError
		if errors.As(err, &valErr) {
			l.Warnw("password.force_update_password_failed.validation_error", "err", err)
		} else {
			l.Errorw("password.force_update_password_failed", "err", err)
		}
		return nil, err
	}

	if err := s.passwordDataProvider.save(ctx, p, uow); err != nil {
		l.Errorw("password.force_update_password_failed", "err", err)
		return nil, err
	}

	l.Infow("password.force_update_password.success")
	return p, nil
}
