package postgresimpl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/domain/event"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDataDeletionTasksInboxRepository struct {
	conn *pgxpool.Conn
}

func NewUserDataDeletionTasksInboxRepository(conn *pgxpool.Conn) interfaces.EventRepository {
	return &UserDataDeletionTasksInboxRepository{
		conn: conn,
	}
}

func (r *UserDataDeletionTasksInboxRepository) CreateEvents(ctx context.Context, events []*event.Event) error {
	if len(events) == 0 {
		return nil
	}

	ids := make([]string, len(events))
	payloads := make([][]byte, len(events))
	statuses := make([]int16, len(events))
	attempts := make([]int16, len(events))
	createdAts := make([]time.Time, len(events))
	updatedAts := make([]time.Time, len(events))
	for i, e := range events {
		ids[i] = e.GetID()
		payloads[i] = e.GetPayload()
		statuses[i] = int16(e.GetStatus())
		attempts[i] = e.GetAttempts()
		createdAts[i] = e.GetCreatedAt()
		updatedAts[i] = e.GetUpdatedAt()
	}

	const sql = `
		INSERT INTO user_data_deletion_tasks_inbox (id, payload, status, attempts, created_at, updated_at)
		SELECT u.id, u.payload, u.status, u.attempts, u.created_at, u.updated_at
		FROM UNNEST($1::uuid[], $2::jsonb[], $3::smallint[], $4::smallint[], $5::timestamptz[], $6::timestamptz[])
			AS u(id, payload, status, attempts, created_at, updated_at)
		ON CONFLICT (id) DO UPDATE
		SET
			payload = EXCLUDED.payload,
			status = EXCLUDED.status,
			attempts = EXCLUDED.attempts,
			updated_at = EXCLUDED.updated_at
		RETURNING id, payload, status, attempts, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, ids, payloads, statuses, attempts, createdAts, updatedAts)
	if err != nil {
		return fmt.Errorf("insert inbox events: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var res models.V1EventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan inserted inbox event: %w", err)
		}
		events[i] = res.ToDomain()
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate inserted inbox events: %w", err)
	}

	return nil
}

func (r *UserDataDeletionTasksInboxRepository) UpdateEvents(ctx context.Context, events []*event.Event) error {
	if len(events) == 0 {
		return nil
	}

	eventDals := make([]models.V1EventDal, len(events))
	for i, e := range events {
		eventDals[i] = models.V1EventFromDomain(e)
	}

	const sql = `
		UPDATE user_data_deletion_tasks_inbox AS t
		SET
			payload = u.payload,
			status = u.status,
			attempts = u.attempts,
			updated_at = u.updated_at
		FROM UNNEST($1::v1_inbox_outbox_event[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.payload, t.status, t.attempts, t.created_at, t.updated_at
	`

	rows, err := r.conn.Query(ctx, sql, eventDals)
	if err != nil {
		return fmt.Errorf("update inbox events: %w", err)
	}
	defer rows.Close()

	eventById := make(map[string]models.V1EventDal)
	for rows.Next() {
		var res models.V1EventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan updated inbox event: %w", err)
		}
		eventById[res.Id] = res
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate updated inbox events: %w", err)
	}
	for _, e := range events {
		id := e.GetID()
		*e = *eventById[id].ToDomain()
	}

	return nil
}

func (r *UserDataDeletionTasksInboxRepository) DeleteEvents(ctx context.Context, events []*event.Event) error {
	if len(events) == 0 {
		return nil
	}

	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.GetID())
	}

	const sql = `
		DELETE FROM user_data_deletion_tasks_inbox i
		WHERE i.id::text = ANY($1)
	`

	if _, err := r.conn.Exec(ctx, sql, ids); err != nil {
		return fmt.Errorf("delete inbox events: %w", err)
	}

	return nil
}

func (r *UserDataDeletionTasksInboxRepository) QueryEvents(ctx context.Context, query *models.QueryEventsDal) ([]*event.Event, error) {
	if query == nil {
		return nil, nil
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT e.id, e.payload, e.status, e.attempts, e.created_at, e.updated_at
		FROM user_data_deletion_tasks_inbox e
		WHERE 1=1
	`)

	appendAnyEqual(&sb, "e.id", query.Filter.Ids, &args, &argPos)
	appendAnyEqual(&sb, "e.status", query.Filter.Statuses, &args, &argPos)
	appendRange(&sb, "e.attempts", query.Filter.AttemptsFrom, query.Filter.AttemptsTo, &args, &argPos)
	appendRange(&sb, "e.created_at", query.Filter.CreatedFrom, query.Filter.CreatedTo, &args, &argPos)
	appendRange(&sb, "e.updated_at", query.Filter.UpdatedFrom, query.Filter.UpdatedTo, &args, &argPos)
	appendOrder(&sb, "e.id", true)
	appendLimitOffset(&sb, query.Limit, query.Offset, &args, &argPos)

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query inbox events: %w", err)
	}
	defer rows.Close()

	var result []*event.Event
	for rows.Next() {
		var res models.V1EventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan inbox event: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inbox events: %w", err)
	}

	return result, nil
}

func (r *UserDataDeletionTasksInboxRepository) QueryEventsLocked(ctx context.Context, query *models.QueryEventsLockedDal) ([]*event.Event, error) {
	if query == nil {
		return nil, nil
	}

	const sql = `
		SELECT e.id, e.payload, e.status, e.attempts, e.created_at, e.updated_at
		FROM user_data_deletion_tasks_inbox e
		WHERE
			(
				e.status = $1
				OR (e.status = $2 AND e.updated_at < $3 AND e.attempts <= $4)
			)
			AND ($5::timestamptz IS NULL OR e.created_at >= $5)
		ORDER BY e.created_at ASC
		LIMIT $6
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.conn.Query(ctx, sql,
		int16(event.EventStatusPending),
		int16(event.EventStatusFailed),
		query.RetryAfter,
		int16(event.MaxAttempts-1),
		query.CreatedAfter,
		query.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query inbox events: %w", err)
	}
	defer rows.Close()

	var result []*event.Event
	for rows.Next() {
		var res models.V1EventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan inbox event: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inbox events: %w", err)
	}

	return result, nil
}
