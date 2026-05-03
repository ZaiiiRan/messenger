package tokenservice

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/auth/v1"
	userpb "github.com/ZaiiiRan/messenger/backend/auth-service/gen/go/user/v1"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	uow "github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/unitofwork/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/postgres"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/transport/redis"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/utils"
	"github.com/ZaiiiRan/messenger/backend/go-common/pkg/ctxmetadata"
	commonjwt "github.com/ZaiiiRan/messenger/backend/go-common/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenService interface {
	GenerateToken(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User, userVersion *userversion.UserVersion, existedRefreshToken *token.Token) (*token.Token, *token.Token, error)
	ValidateRefreshToken(ctx context.Context, uow *uow.UnitOfWork, refreshToken string) (*token.Token, *userversion.UserVersion, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (*commonjwt.UserClaims, error)
	ParseRefreshToken(tokenStr string) error
	InvalidateRefreshToken(ctx context.Context, uow *uow.UnitOfWork, refreshToken string) error
	GetUserVersion(ctx context.Context, uow *uow.UnitOfWork, userId string) (*userversion.UserVersion, error)
	UpdateUserVersion(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User) (*userversion.UserVersion, error)
	DeleteExpiredTokens(ctx context.Context, uow *uow.UnitOfWork, batchSize uint, workerID string) error
	GetRefreshTokens(ctx context.Context, uow *uow.UnitOfWork, refreshToken *token.Token, userVersion *userversion.UserVersion, req *pb.GetActiveSessionsRequest) ([]*token.Token, error)
	InvalidateRefreshTokensByIds(ctx context.Context, uow *uow.UnitOfWork, currentToken *token.Token, ids []int64) error
}

type tokenService struct {
	tokenDataProvider       *tokenDataProvider
	userVersionDataProvider *userVersionDataProvider
	jwtSettings             *settings.JWTSettings
	log                     *zap.SugaredLogger
}

func New(jwtSettings settings.JWTSettings, pgClient *postgres.PostgresClient, redisClient *redis.RedisClient, log *zap.SugaredLogger) TokenService {
	return &tokenService{
		tokenDataProvider:       newTokenDataProvider(pgClient, redisClient),
		userVersionDataProvider: newUserVersionDataProvider(pgClient, redisClient),
		jwtSettings:             &jwtSettings,
		log:                     log,
	}
}

func (s *tokenService) GenerateToken(
	ctx context.Context,
	uow *uow.UnitOfWork,
	user *userpb.User,
	userVersion *userversion.UserVersion,
	existedRefreshToken *token.Token,
) (*token.Token, *token.Token, error) {
	l := s.log.With("op", "generate_tokens", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "user_id", user.Id)

	var version int
	if existedRefreshToken != nil {
		version = existedRefreshToken.GetVersion()
	} else if userVersion != nil {
		version = userVersion.GetVersion()
	} else {
		l.Errorw("token.get_user_version", "err", "user version or existed refresh token is not provided")
		return nil, nil, fmt.Errorf("user version or existed refresh token is not provided")
	}

	c := &commonjwt.UserClaims{
		Id:                  user.Id,
		Username:            user.Username,
		Email:               user.Email,
		IsConfirmed:         user.Status.IsConfirmed,
		IsDeleted:           user.Status.IsDeleted,
		IsPermanentlyBanned: user.Status.IsPermanentlyBanned,
		IsTemporarilyBanned: utils.IsActiveTemporaryBan(user.Status.BannedUntil),
		Version:             version,
	}

	access, accessExp, err := signToken(*c, []byte(s.jwtSettings.AccessTokenSecret), time.Duration(s.jwtSettings.AccessTokenTTL)*time.Second)
	if err != nil {
		l.Errorw("token.sign_access_failed", "err", err)
		return nil, nil, err
	}

	refresh, refreshExp, err := signToken(*c, []byte(s.jwtSettings.RefreshTokenSecret), time.Duration(s.jwtSettings.RefreshTokenTTL)*time.Second)
	if err != nil {
		l.Errorw("token.sign_refresh_failed", "err", err)
		return nil, nil, err
	}

	accessToken := token.New(user.Id, access, token.AccessTokenType, version, "", "", "", "", "", accessExp)

	var refreshToken *token.Token
	if existedRefreshToken != nil {
		if err := s.tokenDataProvider.deleteFromCache(ctx, existedRefreshToken.GetToken()); err != nil {
			l.Errorw("token.delete_existed_from_cache", "err", err)
			return nil, nil, err
		}
		refreshToken = existedRefreshToken
		refreshToken.SetToken(refresh, refreshExp)
	} else {
		ip, country, city, os, browser := extractSessionInfo(ctx)
		refreshToken = token.New(user.Id, refresh, token.RefreshTokenType, version, ip, country, city, os, browser, refreshExp)
	}

	if err := s.tokenDataProvider.save(ctx, refreshToken, uow); err != nil {
		l.Errorw("token.save_token_failed", "err", err)
		return nil, nil, err
	}

	l.Infow("token.generate.success")
	return accessToken, refreshToken, nil
}

func (s *tokenService) ValidateRefreshToken(ctx context.Context, uow *uow.UnitOfWork, refreshToken string) (*token.Token, *userversion.UserVersion, error) {
	l := s.log.With("op", "validate_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := commonjwt.ParseUserToken(refreshToken, []byte(s.jwtSettings.RefreshTokenSecret))
	if err != nil {
		l.Warnw("token.refresh_token_parse_failed", "err", err)
		return nil, nil, commonjwt.ErrInvalidToken
	}

	t, err := s.tokenDataProvider.getByToken(ctx, refreshToken, uow)
	if err != nil {
		l.Errorw("token.get_token_failed", "err", err)
		return nil, nil, err
	}
	if t == nil || cl.Id != t.GetUserID() || cl.Version != t.GetVersion() {
		l.Warnw("token.refresh_token_invalid")
		return nil, nil, commonjwt.ErrInvalidToken
	}

	userVersion, err := s.userVersionDataProvider.getByUserId(ctx, t.GetUserID(), uow)
	if err != nil {
		l.Errorw("token.get_user_version_failed", "err", err)
		return nil, nil, err
	}
	if userVersion == nil || t.GetVersion() != userVersion.GetVersion() {
		l.Warnw("token.refresh_token_invalid")
		return nil, nil, commonjwt.ErrInvalidToken
	}

	l.Infow("token.refresh_token_valid", "user_id", cl.Id)
	return t, userVersion, nil
}

func (s *tokenService) ParseRefreshToken(tokenStr string) error {
	if _, err := commonjwt.ParseUserToken(tokenStr, []byte(s.jwtSettings.RefreshTokenSecret)); err != nil {
		return commonjwt.ErrInvalidToken
	}
	return nil
}

func (s *tokenService) ValidateAccessToken(ctx context.Context, accessToken string) (*commonjwt.UserClaims, error) {
	l := s.log.With("op", "validate_access_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	cl, err := commonjwt.ParseUserToken(accessToken, []byte(s.jwtSettings.AccessTokenSecret))
	if err != nil {
		l.Warnw("token.access_token_parse_failed", "err", err)
		return nil, commonjwt.ErrInvalidToken
	}

	l.Infow("token.access_token_valid", "user_id", cl.Id)
	return cl, nil
}

func (s *tokenService) InvalidateRefreshToken(ctx context.Context, uow *uow.UnitOfWork, refreshToken string) error {
	l := s.log.With("op", "invalidate_refresh_token", "req_id", ctxmetadata.GetReqIdFromContext(ctx))

	err := s.tokenDataProvider.delete(ctx, refreshToken, uow)
	if err != nil {
		l.Errorw("token.delete_refresh_token_failed", "err", err)
		return err
	}

	l.Infow("token.invalidate_refresh_token.success")
	return nil
}

func (s *tokenService) InvalidateRefreshTokensByIds(ctx context.Context, uow *uow.UnitOfWork, currentToken *token.Token, ids []int64) error {
	l := s.log.With("op", "invalidate_refresh_tokens_by_ids", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "user_id", currentToken.GetUserID())

	if err := s.tokenDataProvider.deleteByIds(ctx, currentToken.GetUserID(), ids, currentToken.GetID(), uow); err != nil {
		l.Errorw("token.invalidate_refresh_tokens_by_ids_failed", "err", err)
		return err
	}
	l.Infow("token.invalidate_refresh_tokens_by_ids.success", "count", len(ids))
	return nil
}

func (s *tokenService) GetUserVersion(ctx context.Context, uow *uow.UnitOfWork, userId string) (*userversion.UserVersion, error) {
	l := s.log.With("op", "get_user_version", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "user_id", userId)

	uv, err := s.userVersionDataProvider.getByUserId(ctx, userId, uow)
	if err != nil {
		l.Errorw("token.get_user_version_failed", "err", err)
		return nil, err
	}
	l.Infow("token.get_user_version.success")
	return uv, nil
}

func (s *tokenService) UpdateUserVersion(ctx context.Context, uow *uow.UnitOfWork, user *userpb.User) (*userversion.UserVersion, error) {
	l := s.log.With("op", "update_user_version", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "user_id", user.Id)

	var uv *userversion.UserVersion

	existedUserVersion, err := s.userVersionDataProvider.getByUserId(ctx, user.Id, uow)
	if err != nil {
		l.Errorw("token.get_existed_user_version_failed", "err", err)
		return nil, err
	}

	if existedUserVersion != nil {
		uv = existedUserVersion
		uv.IncrementVersion()
	} else {
		uv = userversion.New(user.Id)
	}

	if err := s.userVersionDataProvider.save(ctx, uv, uow); err != nil {
		l.Errorw("token.save_user_version_failed", "err", err)
		return nil, err
	}

	l.Infow("token.update_user_version.success")
	return uv, nil
}

func (s *tokenService) GetRefreshTokens(
	ctx context.Context,
	uow *uow.UnitOfWork,
	refreshToken *token.Token,
	userVersion *userversion.UserVersion,
	req *pb.GetActiveSessionsRequest,
) ([]*token.Token, error) {
	l := s.log.With("op", "get_refresh_tokens", "req_id", ctxmetadata.GetReqIdFromContext(ctx), "user_id", refreshToken.GetUserID())

	query := models.NewQueryTokensDal(refreshToken.GetUserID(), refreshToken.GetToken(), userVersion.GetVersion(), int(req.Page), int(req.PageSize))

	tokens, err := s.tokenDataProvider.getActiveByUserId(ctx, query, uow)
	if err != nil {
		l.Errorw("token.get_refresh_tokens_failed", "err", err)
		return nil, err
	}
	l.Infow("token.get_refresh_tokens.success", "count", len(tokens))
	return tokens, nil
}

func (s *tokenService) DeleteExpiredTokens(ctx context.Context, uow *uow.UnitOfWork, batchSize uint, workerID string) error {
	l := s.log.With("op", "delete_expired_tokens", "worker_id", workerID)

	tokens, err := s.tokenDataProvider.deleteExpiredTokens(ctx, batchSize, uow)
	if err != nil {
		l.Errorw("token.delete_expired_tokens_failed", "err", err)
		return err
	}
	l.Infow("token.delete_expired_tokens.success", "count", len(tokens))
	return nil
}

func signToken(c commonjwt.UserClaims, key []byte, ttl time.Duration) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(ttl)
	safeNbf := now.Add(-10 * time.Second)

	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(safeNbf),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	str, err := t.SignedString(key)
	if err != nil {
		return "", expiresAt, err
	}

	return str, expiresAt, nil
}

func extractSessionInfo(ctx context.Context) (ip, country, city, os, browser string) {
	ip, _ = ctxmetadata.GetRealIPFromIncomingContext(ctx)
	country, _ = ctxmetadata.GetCountryNameFromIncomingContext(ctx)
	city, _ = ctxmetadata.GetCityFromIncomingContext(ctx)

	ua, err := ctxmetadata.GetUAFromIncomingContext(ctx)
	if err == nil && ua != "" {
		parsed := utils.ParseUserAgent(ua)
		os = parsed.OS
		browser = parsed.Browser
	}
	return
}
