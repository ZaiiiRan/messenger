package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	emailchangecode "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code/email_change_code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type emailChangeCodeRepository struct {
	conn *pgxpool.Conn
}

func NewEmailChangeCodeRepository(conn *pgxpool.Conn) interfaces.EmailChangeCodeRepository {
	return &emailChangeCodeRepository{
		conn: conn,
	}
}

func (r *emailChangeCodeRepository) CreateCode(ctx context.Context, c *emailchangecode.EmailChangeCode) error {
	dal := models.V1EmailCodeDalFromDomain(c)
	sql := `
		INSERT INTO change_email_tokens (user_id, email, code, link_token, generations_left, verifications_left, expires_at)
		SELECT (i).user_id, (i).email, (i).code, (i).link_token, (i).generations_left, (i).verifications_left, (i).expires_at
		FROM UNNEST($1::v1_email_code[]) i
		RETURNING id, user_id, email, code, link_token, generations_left, verifications_left, expires_at, created_at, updated_at
	`

	var res models.V1EmailCodeDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1EmailCodeDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Email, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create code: %w", err)
	}
	*c = *res.ToDomain()
	return nil
}

func (r *emailChangeCodeRepository) UpdateCode(ctx context.Context, c *emailchangecode.EmailChangeCode) error {
	dal := models.V1EmailCodeDalFromDomain(c)
	sql := `
		UPDATE change_email_tokens AS t
		SET
			email              = u.email,
			code               = u.code,
			link_token         = u.link_token,
			generations_left   = u.generations_left,
			verifications_left = u.verifications_left,
			expires_at         = u.expires_at,
			updated_at         = u.updated_at
		FROM UNNEST($1::v1_email_code[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.user_id, t.email, t.code, t.link_token, t.generations_left, t.verifications_left, t.expires_at, t.created_at, t.updated_at
	`

	var res models.V1EmailCodeDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1EmailCodeDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Email, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update code: %w", err)
	}
	*c = *res.ToDomain()
	return nil
}

func (r *emailChangeCodeRepository) DeleteCode(ctx context.Context, c *emailchangecode.EmailChangeCode) error {
	var (
		query string
		arg   any
	)
	if c.GetID() != 0 {
		query = `DELETE FROM change_email_tokens WHERE id = $1`
		arg = c.GetID()
	} else {
		query = `DELETE FROM change_email_tokens WHERE user_id::text = $1`
		arg = c.GetUserID()
	}
	if _, err := r.conn.Exec(ctx, query, arg); err != nil {
		return fmt.Errorf("delete code: %w", err)
	}
	return nil
}

func (r *emailChangeCodeRepository) QueryCode(ctx context.Context, query *models.QueryCodeDal) (*emailchangecode.EmailChangeCode, error) {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)
	sb.WriteString(`
		SELECT id, user_id, email, code, link_token, generations_left, verifications_left, expires_at, created_at, updated_at
		FROM change_email_tokens
		WHERE 1=1
	`)
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	appendEqual(&sb, "link_token", query.LinkToken, &args, &argPos)
	sb.WriteString(" LIMIT 1")
	if query.ForUpdate {
		sb.WriteString(" FOR UPDATE")
	}

	var res models.V1EmailCodeDal
	err := r.conn.QueryRow(ctx, sb.String(), args...).Scan(
		&res.Id, &res.UserId, &res.Email, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query code: %w", err)
	}
	return res.ToDomain(), nil
}

func (r *emailChangeCodeRepository) DeleteExpiredCodes(ctx context.Context, query *models.QueryExpiredCodesDal) ([]*emailchangecode.EmailChangeCode, error) {
	sql := `
		DELETE FROM change_email_tokens AS c
		WHERE c.id IN (
			SELECT c.id
			FROM change_email_tokens AS c
			WHERE c.expires_at < NOW()
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		RETURNING id, user_id, email, code, link_token, generations_left, verifications_left, expires_at, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, query.PageSize)
	if err != nil {
		return nil, fmt.Errorf("delete expired codes: %w", err)
	}
	defer rows.Close()

	var result []*emailchangecode.EmailChangeCode
	for rows.Next() {
		var res models.V1EmailCodeDal
		if err := rows.Scan(
			&res.Id, &res.UserId, &res.Email, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan expired codes: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expired codes: %w", err)
	}

	return result, nil
}
