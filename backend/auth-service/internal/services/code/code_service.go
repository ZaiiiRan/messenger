package codeservice

import (
	"context"
	"errors"

	userpb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	codedomain "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"go.uber.org/zap"
)

type CodeService interface {
	CheckConfirmationCode(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawCode string) (bool, error)
	GenerateConfiramtionCode(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User) (*codedomain.Code, error)
}

type codeService struct {
	codeDataProvider *codeDataProvider
	log              *zap.SugaredLogger
}

func New(pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) CodeService {
	return &codeService{
		codeDataProvider: newCodeDataProvider(pgClient, redisClient),
		log:              log,
	}
}

func (s *codeService) CheckConfirmationCode(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, rawCode string) (bool, error) {
	l := s.log.With("op", "check_confirmation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	code, err := s.codeDataProvider.getByUserID(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("code.check_confirmation_code_failed", "err", err)
		return false, err
	}
	if code == nil {
		l.Errorw("code.check_confirmation_code_failed", "err", "confirmation code not found")
		return false, nil
	}

	valid, err := code.CheckCode(rawCode)
	if err != nil {
		l.Warnw("code.check_confirmation_code_failed", "err", err)
		return false, err
	}
	if !valid {
		l.Warnw("code.check_confirmation_code_failed", "err", "invalid code")
		return false, nil
	}

	if err := s.codeDataProvider.delete(ctx, code, uow); err != nil {
		l.Errorw("code.check_confirmation_code_failed", "err", err)
		return false, err
	}

	l.Infow("code.check_confirmation_code.success")
	return true, nil
}

func (s *codeService) GenerateConfiramtionCode(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User) (*codedomain.Code, error) {
	l := s.log.With("op", "generate_confirmation_code", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	var c *codedomain.Code
	existedCode, err := s.codeDataProvider.getByUserIDLocked(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("code.generate_confirmation_code_failed", "err", err)
		return nil, err
	}
	if existedCode != nil {
		c = existedCode
		err = c.GenerateCode()
		if err != nil {
			var valErr *codedomain.CodeValidationError
			if errors.As(err, &valErr) {
				l.Warnw("code.generate_confirmation_code_failed.validation_error", "err", err)
			} else {
				l.Errorw("code.generate_confirmation_code_failed", "err", err)
			}
			return nil, err
		}
	} else {
		c, err = codedomain.New(user.Id)
		if err != nil {
			var valErr *codedomain.CodeValidationError
			if errors.As(err, &valErr) {
				l.Warnw("code.generate_confirmation_code_failed.validation_error", "err", err)
			} else {
				l.Errorw("code.generate_confirmation_code_failed", "err", err)
			}
			return nil, err
		}
	}

	if err := s.codeDataProvider.save(ctx, c, uow); err != nil {
		l.Errorw("code.generate_confirmation_code_failed", "err", err)
		return nil, err
	}

	l.Infow("code.generate_confirmation_code.success")
	return c, nil
}
