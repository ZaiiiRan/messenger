package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/code"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type codeRepository struct {
	conn       *pgxpool.Conn
	tableName  string
	pgTypeName string
	codeType   code.CodeType
}

func NewActivationCodeRepository(conn *pgxpool.Conn) interfaces.CodeRepository {
	return &codeRepository{
		conn:       conn,
		tableName:  "confirmation_codes",
		pgTypeName: "v1_confirmation_code",
		codeType:   code.CodeTypeActivation,
	}
}

func NewPasswordResetCodeRepository(conn *pgxpool.Conn) interfaces.CodeRepository {
	return &codeRepository{
		conn:       conn,
		tableName:  "password_reset_tokens",
		pgTypeName: "v1_password_reset_token",
		codeType:   code.CodeTypePasswordReset,
	}
}

func (r *codeRepository) CreateCode(ctx context.Context, c *code.Code) error {
	dal := models.V1CodeDalFromDomain(c)
	sql := fmt.Sprintf(`
		INSERT INTO %s (user_id, code, link_token, generations_left, verifications_left, expires_at)
		SELECT (i).user_id, (i).code, (i).link_token, (i).generations_left, (i).verifications_left, (i).expires_at
		FROM UNNEST($1::%s[]) i
		RETURNING id, user_id, code, link_token, generations_left, verifications_left, expires_at, created_at, updated_at
	`, r.tableName, r.pgTypeName)

	var res models.V1CodeDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1CodeDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create code: %w", err)
	}
	*c = *res.ToDomain(r.codeType)
	return nil
}

func (r *codeRepository) UpdateCode(ctx context.Context, c *code.Code) error {
	dal := models.V1CodeDalFromDomain(c)
	sql := fmt.Sprintf(`
		UPDATE %s AS t
		SET
			code               = u.code,
			link_token         = u.link_token,
			generations_left   = u.generations_left,
			verifications_left = u.verifications_left,
			expires_at         = u.expires_at,
			updated_at         = u.updated_at
		FROM UNNEST($1::%s[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.user_id, t.code, t.link_token, t.generations_left, t.verifications_left, t.expires_at, t.created_at, t.updated_at
	`, r.tableName, r.pgTypeName)

	var res models.V1CodeDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1CodeDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update code: %w", err)
	}
	*c = *res.ToDomain(r.codeType)
	return nil
}

func (r *codeRepository) DeleteCode(ctx context.Context, c *code.Code) error {
	var (
		query string
		arg   any
	)
	if c.GetID() != 0 {
		query = fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, r.tableName)
		arg = c.GetID()
	} else {
		query = fmt.Sprintf(`DELETE FROM %s WHERE user_id::text = $1`, r.tableName)
		arg = c.GetUserID()
	}
	if _, err := r.conn.Exec(ctx, query, arg); err != nil {
		return fmt.Errorf("delete code: %w", err)
	}
	return nil
}

func (r *codeRepository) QueryCode(ctx context.Context, query *models.QueryCodeDal) (*code.Code, error) {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)
	sb.WriteString(fmt.Sprintf(`
		SELECT id, user_id, code, link_token, generations_left, verifications_left, expires_at, created_at, updated_at
		FROM %s
		WHERE 1=1
	`, r.tableName))
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	appendEqual(&sb, "link_token", query.LinkToken, &args, &argPos)
	sb.WriteString(" LIMIT 1")
	if query.ForUpdate {
		sb.WriteString(" FOR UPDATE")
	}

	var res models.V1CodeDal
	err := r.conn.QueryRow(ctx, sb.String(), args...).Scan(
		&res.Id, &res.UserId, &res.Code, &res.LinkToken, &res.GenerationsLeft, &res.VerificationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query code: %w", err)
	}
	return res.ToDomain(r.codeType), nil
}
