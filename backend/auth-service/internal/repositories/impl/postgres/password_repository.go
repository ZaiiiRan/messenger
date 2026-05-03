package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/password"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordRepository struct {
	conn *pgxpool.Conn
}

func NewPasswordRepository(conn *pgxpool.Conn) interfaces.PasswordRepository {
	return &PasswordRepository{conn: conn}
}

func (r *PasswordRepository) CreatePassword(ctx context.Context, p *password.Password) error {
	dal := models.V1PasswordDalFromDomain(p)
	sql := `
		INSERT INTO passwords (user_id, password_hash)
		SELECT (i).user_id, (i).password_hash
		FROM UNNEST($1::v1_password[]) i
		RETURNING id, user_id, password_hash, created_at, updated_at
	`
	var res models.V1PasswordDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1PasswordDal{dal}).Scan(
		&res.Id, &res.UserId, &res.PasswordHash, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create password: %w", err)
	}
	*p = *res.ToDomain()
	return nil
}

func (r *PasswordRepository) UpdatePassword(ctx context.Context, p *password.Password) error {
	dal := models.V1PasswordDalFromDomain(p)
	sql := `
		UPDATE passwords AS t
		SET
			password_hash = u.password_hash,
			updated_at    = u.updated_at
		FROM UNNEST($1::v1_password[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.user_id, t.password_hash, t.created_at, t.updated_at
	`
	var res models.V1PasswordDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1PasswordDal{dal}).Scan(
		&res.Id, &res.UserId, &res.PasswordHash, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	*p = *res.ToDomain()
	return nil
}

func (r *PasswordRepository) DeletePassword(ctx context.Context, p *password.Password) error {
	var (
		query string
		arg   any
	)
	if p.GetID() != 0 {
		query = `DELETE FROM passwords WHERE id = $1`
		arg = p.GetID()
	} else {
		query = `DELETE FROM passwords WHERE user_id::text = $1`
		arg = p.GetUserID()
	}
	if _, err := r.conn.Exec(ctx, query, arg); err != nil {
		return fmt.Errorf("delete password: %w", err)
	}
	return nil
}

func (r *PasswordRepository) QueryPassword(ctx context.Context, query *models.QueryPasswordDal) (*password.Password, error) {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)
	sb.WriteString(`
		SELECT id, user_id, password_hash, created_at, updated_at
		FROM passwords
		WHERE 1=1
	`)
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	sb.WriteString(" LIMIT 1")

	var res models.V1PasswordDal
	err := r.conn.QueryRow(ctx, sb.String(), args...).Scan(
		&res.Id, &res.UserId, &res.PasswordHash, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query password: %w", err)
	}
	return res.ToDomain(), nil
}
