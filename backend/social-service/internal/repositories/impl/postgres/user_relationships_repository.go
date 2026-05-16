package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	userrelationship "github.com/ZaiiiRan/messenger/backend/social-service/internal/domain/user_relationship"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRelationshipsRepository struct {
	conn *pgxpool.Conn
}

func NewUserRelationshipsRepository(conn *pgxpool.Conn) interfaces.UserRelationshipsRepository {
	return &UserRelationshipsRepository{conn: conn}
}

func (r *UserRelationshipsRepository) CreateUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error {
	const sql = `
		INSERT INTO user_relationships (user_id_1, user_id_2, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id_1, user_id_2, status, created_at, updated_at
	`
	var res models.V1UserRelationshipDal
	row := r.conn.QueryRow(ctx, sql,
		ur.GetUserID1(), ur.GetUserID2(),
		int16(ur.GetStatus()),
		ur.GetCreatedAt(), ur.GetUpdatedAt(),
	)
	if err := row.Scan(&res.UserId1, &res.UserId2, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
		return fmt.Errorf("insert user relationship: %w", err)
	}
	*ur = *res.ToDomain()
	return nil
}

func (r *UserRelationshipsRepository) UpdateUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error {
	const sql = `
		UPDATE user_relationships AS t
		SET
			status     = ur.status,
			updated_at = ur.updated_at
		FROM UNNEST($1::v1_user_relationship[]) AS ur
		WHERE t.user_id_1 = ur.user_id_1 AND t.user_id_2 = ur.user_id_2
		RETURNING t.user_id_1, t.user_id_2, t.status, t.created_at, t.updated_at
	`
	dal := models.V1UserRelationshipDalFromDomain(ur)

	rows, err := r.conn.Query(ctx, sql, []models.V1UserRelationshipDal{dal})
	if err != nil {
		return fmt.Errorf("update user relationship: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var res models.V1UserRelationshipDal
		if err := rows.Scan(&res.UserId1, &res.UserId2, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan updated user relationship: %w", err)
		}
		*ur = *res.ToDomain()
	}
	return rows.Err()
}

func (r *UserRelationshipsRepository) DeleteUserRelationship(ctx context.Context, ur *userrelationship.UserRelationship) error {
	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`DELETE FROM user_relationships WHERE 1=1`)
	uid1, uid2 := ur.GetUserID1(), ur.GetUserID2()
	appendEqual(&sb, "user_id_1::text", &uid1, &args, &argPos)
	appendEqual(&sb, "user_id_2::text", &uid2, &args, &argPos)

	if _, err := r.conn.Exec(ctx, sb.String(), args...); err != nil {
		return fmt.Errorf("delete user relationship: %w", err)
	}
	return nil
}

func (r *UserRelationshipsRepository) QueryUserRelationships(
	ctx context.Context,
	query *models.QueryUserRelationshipsDal,
	forUpdate bool,
) ([]*userrelationship.UserRelationship, error) {
	if query == nil {
		return nil, nil
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT ur.user_id_1, ur.user_id_2, ur.status, ur.created_at, ur.updated_at
		FROM user_relationships ur
		WHERE 1=1
	`)

	if query.FirstUserId != nil {
		fmt.Fprintf(&sb, " AND (%s = $%d OR %s = $%d)", "ur.user_id_1::text", argPos, "ur.user_id_2::text", argPos)
		args = append(args, *query.FirstUserId)
		argPos++
	}
	if len(query.SecondUserIds) > 0 {
		fmt.Fprintf(&sb, " AND (%s = ANY($%d) OR %s = ANY($%d))", "ur.user_id_1::text", argPos, "ur.user_id_2::text", argPos)
		args = append(args, query.SecondUserIds)
		argPos++
	}

	appendAnyEqual(&sb, "ur.status", query.Statuses, &args, &argPos)
	if query.OrderByUpdatedAtDesc {
		appendOrder(&sb, "ur.updated_at", false)
	}
	appendLimitOffset(&sb, query.Limit, query.Offset, &args, &argPos)

	if forUpdate {
		sb.WriteString(" FOR UPDATE")
	}

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query user relationships: %w", err)
	}
	defer rows.Close()

	var result []*userrelationship.UserRelationship
	for rows.Next() {
		var res models.V1UserRelationshipDal
		if err := rows.Scan(&res.UserId1, &res.UserId2, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user relationship: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user relationships: %w", err)
	}

	return result, nil
}
