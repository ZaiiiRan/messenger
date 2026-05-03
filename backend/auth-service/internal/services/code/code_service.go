package codeservice

import (
	"context"
	"errors"
	"time"

	codedomain "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	pgclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	redisclient "github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	"go.uber.org/zap"
)

type CodeService interface {
	GenerateCode(ctx context.Context, uow *uow.UnitOfWork, userID string, codeType codedomain.CodeType) (*codedomain.Code, error)
	CheckCodeByCode(ctx context.Context, uow *uow.UnitOfWork, userID, rawCode string, codeType codedomain.CodeType) (bool, error)
	CheckCodeByLinkToken(ctx context.Context, uow *uow.UnitOfWork, linkToken string, codeType codedomain.CodeType) (userID string, valid bool, err error)
}

type codeService struct {
	codeDataProvider *codeDataProvider
	log              *zap.SugaredLogger
}

func New(pgClient *pgclient.PostgresClient, redisClient *redisclient.RedisClient, log *zap.SugaredLogger) CodeService {
	return &codeService{
		codeDataProvider: newCodeDataProvider(pgClient, redisClient),
		log:              log,
	}
}

func (s *codeService) GenerateCode(ctx context.Context, uow *uow.UnitOfWork, userID string, codeType codedomain.CodeType) (*codedomain.Code, error) {
	l := s.log.With("op", "generate_code", "code_type", codeType, "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	existedCode, err := s.codeDataProvider.getByUserIDLocked(ctx, userID, codeType, uow)
	if err != nil {
		l.Errorw("code.generate_failed", "err", err)
		return nil, err
	}

	var c *codedomain.Code
	if existedCode != nil {
		c = existedCode
		err = c.GenerateCode()
	} else {
		c, err = codedomain.New(userID, codeType)
	}
	if err != nil {
		var valErr *codedomain.CodeValidationError
		if errors.As(err, &valErr) {
			l.Warnw("code.generate_failed.validation_error", "err", err)
		} else {
			l.Errorw("code.generate_failed", "err", err)
		}
		return nil, err
	}

	if err := s.codeDataProvider.save(ctx, c, uow); err != nil {
		l.Errorw("code.generate_failed", "err", err)
		return nil, err
	}

	l.Infow("code.generate.success")
	return c, nil
}

func (s *codeService) CheckCodeByCode(ctx context.Context, uow *uow.UnitOfWork, userID, rawCode string, codeType codedomain.CodeType) (bool, error) {
	l := s.log.With("op", "check_code_by_code", "code_type", codeType, "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	c, err := s.codeDataProvider.getByUserID(ctx, userID, codeType, uow)
	if err != nil {
		l.Errorw("code.check_by_code_failed", "err", err)
		return false, err
	}
	if c == nil {
		l.Warnw("code.check_by_code_failed", "err", "code not found")
		return false, nil
	}

	valid, err := c.CheckCode(rawCode)
	if err != nil {
		l.Warnw("code.check_by_code_failed", "err", err)
		return false, err
	}
	if !valid {
		l.Warnw("code.check_by_code_failed", "err", "invalid code")
		if err := s.codeDataProvider.save(ctx, c, uow); err != nil {
			l.Errorw("code.check_by_code_failed", "err", err)
			return false, err
		}
		return false, nil
	}

	if err := s.codeDataProvider.delete(ctx, c, uow); err != nil {
		l.Errorw("code.check_by_code_failed", "err", err)
		return false, err
	}

	l.Infow("code.check_by_code.success")
	return true, nil
}

func (s *codeService) CheckCodeByLinkToken(ctx context.Context, uow *uow.UnitOfWork, linkToken string, codeType codedomain.CodeType) (string, bool, error) {
	l := s.log.With("op", "check_code_by_link_token", "code_type", codeType, "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	c, err := s.codeDataProvider.getByLinkTokenLocked(ctx, linkToken, codeType, uow)
	if err != nil {
		l.Errorw("code.check_by_link_token_failed", "err", err)
		return "", false, err
	}
	if c == nil {
		return "", false, nil
	}

	if time.Now().After(c.GetExpiresAt()) {
		return "", false, codedomain.NewCodeValidationError("link has expired")
	}

	if err := s.codeDataProvider.delete(ctx, c, uow); err != nil {
		l.Errorw("code.check_by_link_token_failed", "err", err)
		return "", false, err
	}

	l.Infow("code.check_by_link_token.success")
	return c.GetUserID(), true, nil
}
