package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/token"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	conn *pgxpool.Conn
}

func NewTokenRepository(conn *pgxpool.Conn) interfaces.TokenRepository {
	return &TokenRepository{conn: conn}
}

func (r *TokenRepository) CreateToken(ctx context.Context, t *token.Token) error {
	dal := models.V1RefreshTokenFromDomain(t)
	sql := `
		INSERT INTO refresh_tokens (user_id, token, version, expires_at)
		SELECT (i).user_id, (i).token, (i).version, (i).expires_at
		FROM UNNEST($1::v1_refresh_token[]) i
		RETURNING id, user_id, token, version, expires_at, created_at, updated_at
	`
	var res models.V1RefreshTokenDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1RefreshTokenDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Token, &res.Version, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	*t = *res.ToDomain()
	return nil
}

func (r *TokenRepository) UpdateToken(ctx context.Context, t *token.Token) error {
	dal := models.V1RefreshTokenFromDomain(t)
	sql := `
		UPDATE refresh_tokens AS rt
		SET
			token      = u.token,
			version    = u.version,
			expires_at = u.expires_at,
			updated_at = u.updated_at
		FROM UNNEST($1::v1_refresh_token[]) AS u
		WHERE rt.id = u.id
		RETURNING rt.id, rt.user_id, rt.token, rt.version, rt.expires_at, rt.created_at, rt.updated_at
	`
	var res models.V1RefreshTokenDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1RefreshTokenDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Token, &res.Version, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update token: %w", err)
	}
	*t = *res.ToDomain()
	return nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, tokenStr string) error {
	if _, err := r.conn.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, tokenStr); err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}

func (r *TokenRepository) DeleteTokensByUserId(ctx context.Context, userId string) error {
	if _, err := r.conn.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id::text = $1`, userId); err != nil {
		return fmt.Errorf("delete tokens by user_id: %w", err)
	}
	return nil
}

func (r *TokenRepository) QueryToken(ctx context.Context, query *models.QueryTokenDal) (*token.Token, error) {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)
	sb.WriteString(`
		SELECT id, user_id, token, version, expires_at, created_at, updated_at
		FROM refresh_tokens
		WHERE 1=1
	`)
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	appendEqual(&sb, "token", query.Token, &args, &argPos)
	appendEqual(&sb, "version", query.Version, &args, &argPos)
	sb.WriteString(" ORDER BY id DESC LIMIT 1")

	var res models.V1RefreshTokenDal
	err := r.conn.QueryRow(ctx, sb.String(), args...).Scan(
		&res.Id, &res.UserId, &res.Token, &res.Version, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query token: %w", err)
	}
	return res.ToDomain(), nil
}
