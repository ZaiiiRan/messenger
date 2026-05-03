package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	userversion "github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/user_version"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserVersionRepository struct {
	conn *pgxpool.Conn
}

func NewUserVersionRepository(conn *pgxpool.Conn) interfaces.UserVersionRepository {
	return &UserVersionRepository{conn: conn}
}

func (r *UserVersionRepository) CreateUserVersion(ctx context.Context, uv *userversion.UserVersion) error {
	dal := models.V1UserVersionDalFromDomain(uv)
	sql := `
		INSERT INTO user_versions (user_id, version)
		SELECT (i).user_id, (i).version
		FROM UNNEST($1::v1_user_version[]) i
		RETURNING id, user_id, version, created_at, updated_at
	`
	var res models.V1UserVersionDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1UserVersionDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Version, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("create user version: %w", err)
	}
	*uv = *res.ToDomain()
	return nil
}

func (r *UserVersionRepository) UpdateUserVersion(ctx context.Context, uv *userversion.UserVersion) error {
	dal := models.V1UserVersionDalFromDomain(uv)
	sql := `
		UPDATE user_versions AS t
		SET
			version    = u.version,
			updated_at = u.updated_at
		FROM UNNEST($1::v1_user_version[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.user_id, t.version, t.created_at, t.updated_at
	`
	var res models.V1UserVersionDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1UserVersionDal{dal}).Scan(
		&res.Id, &res.UserId, &res.Version, &res.CreatedAt, &res.UpdatedAt,
	); err != nil {
		return fmt.Errorf("update user version: %w", err)
	}
	*uv = *res.ToDomain()
	return nil
}

func (r *UserVersionRepository) DeleteUserVersion(ctx context.Context, uv *userversion.UserVersion) error {
	var (
		query string
		arg   any
	)
	if uv.GetID() != 0 {
		query = `DELETE FROM user_versions WHERE id = $1`
		arg = uv.GetID()
	} else {
		query = `DELETE FROM user_versions WHERE user_id::text = $1`
		arg = uv.GetUserID()
	}
	if _, err := r.conn.Exec(ctx, query, arg); err != nil {
		return fmt.Errorf("delete user version: %w", err)
	}
	return nil
}

func (r *UserVersionRepository) QueryUserVersion(ctx context.Context, query *models.QueryUserVersionDal) (*userversion.UserVersion, error) {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)
	sb.WriteString(`
		SELECT id, user_id, version, created_at, updated_at
		FROM user_versions
		WHERE 1=1
	`)
	appendEqual(&sb, "id", query.Id, &args, &argPos)
	appendEqual(&sb, "user_id::text", query.UserId, &args, &argPos)
	sb.WriteString(" LIMIT 1")

	var res models.V1UserVersionDal
	err := r.conn.QueryRow(ctx, sb.String(), args...).Scan(
		&res.Id, &res.UserId, &res.Version, &res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query user version: %w", err)
	}
	return res.ToDomain(), nil
}
