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

const tokenColumns = `id, user_id, token, version, ip, country, city, os, browser, expires_at, created_at, updated_at`

type scanner interface {
	Scan(dest ...any) error
}

func scanToken(row scanner, res *models.V1RefreshTokenDal) error {
	return row.Scan(
		&res.Id, &res.UserId, &res.Token, &res.Version,
		&res.IP, &res.Country, &res.City, &res.OS, &res.Browser,
		&res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	)
}

func (r *TokenRepository) CreateToken(ctx context.Context, t *token.Token) error {
	dal := models.V1RefreshTokenFromDomain(t)
	sql := `
		INSERT INTO refresh_tokens (user_id, token, version, ip, country, city, os, browser, expires_at)
		SELECT (i).user_id, (i).token, (i).version, (i).ip, (i).country, (i).city, (i).os, (i).browser, (i).expires_at
		FROM UNNEST($1::v1_refresh_token[]) i
		RETURNING ` + tokenColumns

	var res models.V1RefreshTokenDal
	if err := scanToken(r.conn.QueryRow(ctx, sql, []models.V1RefreshTokenDal{dal}), &res); err != nil {
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
		RETURNING ` + tokenColumns

	var res models.V1RefreshTokenDal
	if err := scanToken(r.conn.QueryRow(ctx, sql, []models.V1RefreshTokenDal{dal}), &res); err != nil {
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
	sb.WriteString(`SELECT ` + tokenColumns + ` FROM refresh_tokens WHERE 1=1`)
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	appendEqual(&sb, "token", query.Token, &args, &argPos)
	appendEqual(&sb, "version", query.Version, &args, &argPos)
	sb.WriteString(" ORDER BY id DESC LIMIT 1")

	var res models.V1RefreshTokenDal
	if err := scanToken(r.conn.QueryRow(ctx, sb.String(), args...), &res); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query token: %w", err)
	}
	return res.ToDomain(), nil
}

func (r *TokenRepository) QueryActiveTokens(ctx context.Context, query *models.QueryTokensDal) ([]*token.Token, error) {
	sql := `
		SELECT ` + tokenColumns + `
		FROM refresh_tokens rt
		WHERE rt.user_id::text = $1
			AND rt.token != $2
			AND rt.expires_at > NOW()
			AND rt.version = $3
		ORDER BY rt.id DESC
		LIMIT $4 OFFSET $5`

	rows, err := r.conn.Query(ctx, sql, query.UserId, query.ExcludeToken, query.Version, query.Limit, query.Offset)
	if err != nil {
		return nil, fmt.Errorf("query active tokens: %w", err)
	}
	defer rows.Close()

	var result []*token.Token
	for rows.Next() {
		var res models.V1RefreshTokenDal
		if err := scanToken(rows, &res); err != nil {
			return nil, fmt.Errorf("scan active token: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows active tokens: %w", err)
	}
	return result, nil
}

func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context, query *models.QueryExpiredTokensDal) ([]*token.Token, error) {
	querySql := `
		DELETE FROM refresh_tokens AS rt
		WHERE rt.id IN (
			SELECT rt.id
			FROM refresh_tokens rt
			JOIN user_versions uv ON uv.user_id = rt.user_id
			WHERE uv.version != rt.version
			OR rt.expires_at < NOW()
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		RETURNING ` + tokenColumns

	rows, err := r.conn.Query(ctx, querySql, query.PageSize)
	if err != nil {
		return nil, fmt.Errorf("delete expired tokens: %w", err)
	}
	defer rows.Close()

	var result []*token.Token
	for rows.Next() {
		var res models.V1RefreshTokenDal
		if err := scanToken(rows, &res); err != nil {
			return nil, fmt.Errorf("scan expired tokens: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows expired tokens: %w", err)
	}
	return result, nil
}
