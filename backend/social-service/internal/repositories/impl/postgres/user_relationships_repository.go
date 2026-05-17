package postgresimpl

import (
	"context"
	"fmt"
	"strings"
	"time"

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

func (r *UserRelationshipsRepository) CreateUserRelationships(ctx context.Context, urs []*userrelationship.UserRelationship) error {
	if len(urs) == 0 {
		return nil
	}

	user1Ids := make([]string, len(urs))
	user2Ids := make([]string, len(urs))
	statuses := make([]int16, len(urs))
	createdAts := make([]time.Time, len(urs))
	updatedAts := make([]time.Time, len(urs))
	for i, ur := range urs {
		user1Ids[i] = ur.GetUserID1()
		user2Ids[i] = ur.GetUserID2()
		statuses[i] = int16(ur.GetStatus())
		createdAts[i] = ur.GetCreatedAt()
		updatedAts[i] = ur.GetUpdatedAt()
	}

	const sql = `
		INSERT INTO user_relationships (user_id_1, user_id_2, status, created_at, updated_at)
		SELECT u.user_id_1, u.user_id_2, u.status, u.created_at, u.updated_at
		FROM UNNEST($1::text[], $2::text[], $3::smallint[], $4::timestamptz[], $5::timestamptz[])
			AS u(user_id_1, user_id_2, status, created_at, updated_at)
		RETURNING user_id_1, user_id_2, status, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, user1Ids, user2Ids, statuses, createdAts, updatedAts)
	if err != nil {
		return fmt.Errorf("insert user relationships: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var res models.V1UserRelationshipDal
		if err := rows.Scan(&res.UserId1, &res.UserId2, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan inserted user relationship: %w", err)
		}
		urs[i] = res.ToDomain()
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate inserted user relationships: %w", err)
	}

	return nil
}

func (r *UserRelationshipsRepository) UpdateUserRelationships(ctx context.Context, urs []*userrelationship.UserRelationship) error {
	if len(urs) == 0 {
		return nil
	}

	urDals := make([]models.V1UserRelationshipDal, len(urs))
	for i, ur := range urs {
		urDals[i] = models.V1UserRelationshipDalFromDomain(ur)
	}

	const sql = `
		UPDATE user_relationships AS t
		SET
			status     = ur.status,
			updated_at = ur.updated_at
		FROM UNNEST($1::v1_user_relationship[]) AS ur
		WHERE t.user_id_1 = ur.user_id_1 AND t.user_id_2 = ur.user_id_2
		RETURNING t.user_id_1, t.user_id_2, t.status, t.created_at, t.updated_at
	`

	rows, err := r.conn.Query(ctx, sql, urDals)
	if err != nil {
		return fmt.Errorf("update user relationships: %w", err)
	}
	defer rows.Close()

	urById := make(map[string]models.V1UserRelationshipDal, len(urs))
	for rows.Next() {
		var res models.V1UserRelationshipDal
		if err := rows.Scan(&res.UserId1, &res.UserId2, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan updated user relationship: %w", err)
		}
		urById[res.UserId1+res.UserId2] = res
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate updated user relationships: %w", err)
	}
	for _, ur := range urs {
		id := ur.GetUserID1() + ur.GetUserID2()
		if res, ok := urById[id]; ok {
			*ur = *res.ToDomain()
		}
	}

	return nil
}

func (r *UserRelationshipsRepository) DeleteUserRelationships(ctx context.Context, urs []*userrelationship.UserRelationship) error {
	if len(urs) == 0 {
		return nil
	}

	user1Ids := make([]string, len(urs))
	user2Ids := make([]string, len(urs))
	for i, ur := range urs {
		user1Ids[i] = ur.GetUserID1()
		user2Ids[i] = ur.GetUserID2()
	}

	const sql = `
		DELETE FROM user_relationships
		WHERE (user_id_1::text, user_id_2::text) IN (
			SELECT u1, u2 FROM UNNEST($1::text[], $2::text[]) AS t(u1, u2)
		)
	`

	if _, err := r.conn.Exec(ctx, sql, user1Ids, user2Ids); err != nil {
		return fmt.Errorf("delete user relationships: %w", err)
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
