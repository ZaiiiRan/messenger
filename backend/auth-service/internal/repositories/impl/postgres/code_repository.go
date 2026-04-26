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

type CodeRepository struct {
	conn *pgxpool.Conn
}

func NewCodeRepository(conn *pgxpool.Conn) interfaces.CodeRepository {
	return &CodeRepository{conn: conn}
}

func (r *CodeRepository) CreateCode(ctx context.Context, c *code.Code) error {
	dal := models.V1CodeDalFromDomain(c)
	sql := `
		INSERT INTO confirmation_codes (user_id, code, generations_left, expires_at)
		SELECT (i).user_id, (i).code, (i).generations_left, (i).expires_at
		FROM UNNEST($1::v1_confirmation_code[]) i
		RETURNING id, user_id, code, generations_left, expires_at, created_at, updated_at
	`
	var res models.V1CodeDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1CodeDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Code, &res.GenerationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create code: %w", err)
	}
	*c = *res.ToDomain()
	return nil
}

func (r *CodeRepository) UpdateCode(ctx context.Context, c *code.Code) error {
	dal := models.V1CodeDalFromDomain(c)
	sql := `
		UPDATE confirmation_codes AS t
		SET
			code             = u.code,
			generations_left = u.generations_left,
			expires_at       = u.expires_at,
			updated_at       = u.updated_at
		FROM UNNEST($1::v1_confirmation_code[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.user_id, t.code, t.generations_left, t.expires_at, t.created_at, t.updated_at
	`
	var res models.V1CodeDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1CodeDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Code, &res.GenerationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update code: %w", err)
	}
	*c = *res.ToDomain()
	return nil
}

func (r *CodeRepository) DeleteCode(ctx context.Context, c *code.Code) error {
	var (
		query string
		arg   any
	)
	if c.GetID() != 0 {
		query = `DELETE FROM confirmation_codes WHERE id = $1`
		arg = c.GetID()
	} else {
		query = `DELETE FROM confirmation_codes WHERE user_id::text = $1`
		arg = c.GetUserID()
	}
	if _, err := r.conn.Exec(ctx, query, arg); err != nil {
		return fmt.Errorf("delete code: %w", err)
	}
	return nil
}

func (r *CodeRepository) QueryCode(ctx context.Context, query *models.QueryCodeDal) (*code.Code, error) {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)
	sb.WriteString(`
		SELECT id, user_id, code, generations_left, expires_at, created_at, updated_at
		FROM confirmation_codes
		WHERE 1=1
	`)
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	sb.WriteString(" LIMIT 1")
	if query.ForUpdate {
		sb.WriteString(" FOR UPDATE")
	}

	var res models.V1CodeDal
	err := r.conn.QueryRow(ctx, sb.String(), args...).Scan(
		&res.Id, &res.UserId, &res.Code, &res.GenerationsLeft, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query code: %w", err)
	}
	return res.ToDomain(), nil
}
