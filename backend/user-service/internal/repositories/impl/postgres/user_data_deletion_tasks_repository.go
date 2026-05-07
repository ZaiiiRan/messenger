package postgresimpl

import (
	"context"
	"fmt"
	"strings"
	"time"

	outboxevent "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/outbox_event"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDataDeletionTasksRepository struct {
	conn *pgxpool.Conn
}

func NewUserDataDeletionTasksRepository(conn *pgxpool.Conn) interfaces.OutboxEventRepository {
	return &UserDataDeletionTasksRepository{
		conn: conn,
	}
}

func (r *UserDataDeletionTasksRepository) Create(ctx context.Context, events []*outboxevent.OutboxEvent) error {
	if len(events) == 0 {
		return nil
	}

	payloads := make([][]byte, len(events))
	statuses := make([]int16, len(events))
	attempts := make([]int16, len(events))
	createdAts := make([]time.Time, len(events))
	updatedAts := make([]time.Time, len(events))
	for i, e := range events {
		payloads[i] = e.GetPayload()
		statuses[i] = int16(e.GetStatus())
		attempts[i] = e.GetAttempts()
		createdAts[i] = e.GetCreatedAt()
		updatedAts[i] = e.GetUpdatedAt()
	}

	const sql = `
		INSERT INTO user_data_deletion_tasks_outbox (payload, status, attempts, created_at, updated_at)
		SELECT u.payload, u.status, u.attempts, u.created_at, u.updated_at
		FROM UNNEST($1::jsonb[], $2::smallint[], $3::smallint[], $4::timestamptz[], $5::timestamptz[])
			AS u(payload, status, attempts, created_at, updated_at)
		RETURNING id, payload, status, attempts, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, payloads, statuses, attempts, createdAts, updatedAts)
	if err != nil {
		return fmt.Errorf("insert outbox events: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var res models.V1OutboxEventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan inserted outbox event: %w", err)
		}
		events[i] = res.ToDomain()
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate inserted outbox events: %w", err)
	}

	return nil
}

func (r *UserDataDeletionTasksRepository) Update(ctx context.Context, events []*outboxevent.OutboxEvent) error {
	if len(events) == 0 {
		return nil
	}

	eventDals := make([]models.V1OutboxEventDal, len(events))
	for i, e := range events {
		eventDals[i] = models.V1OutboxEventFromDomain(e)
	}

	const sql = `
		UPDATE user_data_deletion_tasks_outbox AS t
		SET
			payload = e.payload,
			status = e.status,
			attempts = e.attempts,
			updated_at = e.updated_at
		FROM UNNEST($1::v1_outbox_event[]) AS e
		WHERE t.id = e.id
		RETURNING t.id, t.payload, t.status, t.attempts, t.created_at, t.updated_at
	`

	rows, err := r.conn.Query(ctx, sql, eventDals)
	if err != nil {
		return fmt.Errorf("update outbox events: %w", err)
	}
	defer rows.Close()

	eventById := make(map[string]models.V1OutboxEventDal)
	for rows.Next() {
		var res models.V1OutboxEventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan updated outbox event: %w", err)
		}
		eventById[res.Id] = res
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate updated outbox events: %w", err)
	}
	for _, e := range events {
		id := e.GetID()
		*e = *eventById[id].ToDomain()
	}

	return nil
}

func (r *UserDataDeletionTasksRepository) Delete(ctx context.Context, events []*outboxevent.OutboxEvent) error {
	if len(events) == 0 {
		return nil
	}

	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.GetID())
	}

	const sql = `
		DELETE FROM user_data_deletion_tasks_outbox o
		WHERE o.id::text = ANY($1)
	`

	if _, err := r.conn.Exec(ctx, sql, ids); err != nil {
		return fmt.Errorf("delete outbox events: %w", err)
	}

	return nil
}

func (r *UserDataDeletionTasksRepository) Query(ctx context.Context, query *models.QueryOutboxEventsDal) ([]*outboxevent.OutboxEvent, error) {
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
		FROM user_data_deletion_tasks_outbox e
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
		return nil, fmt.Errorf("query outbox events: %w", err)
	}
	defer rows.Close()

	var result []*outboxevent.OutboxEvent
	for rows.Next() {
		var res models.V1OutboxEventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox events: %w", err)
	}

	return result, nil
}

func (r *UserDataDeletionTasksRepository) QueryLocked(ctx context.Context, query *models.QueryOutboxEventsLockedDal) ([]*outboxevent.OutboxEvent, error) {
	if query == nil {
		return nil, nil
	}

	const sql = `
		SELECT e.id, e.payload, e.status, e.attempts, e.created_at, e.updated_at
		FROM user_data_deletion_tasks_outbox e
		WHERE e.status = $1
			OR (e.status = $2 AND e.updated_at <= $3 AND e.attempts <= $4)
		ORDER BY e.created_at ASC
		LIMIT $5
		FOR UPDATE SKIP LOCKED
	`

	rows, err := r.conn.Query(ctx, sql,
		int16(outboxevent.OutboxEventStatusPending),
		int16(outboxevent.OutboxEventStatusFailed),
		query.RetryAfter,
		int16(outboxevent.MaxAttempts-1),
		query.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query outbox events: %w", err)
	}
	defer rows.Close()

	var result []*outboxevent.OutboxEvent
	for rows.Next() {
		var res models.V1OutboxEventDal
		if err := rows.Scan(&res.Id, &res.Payload, &res.Status, &res.Attempts, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan outbox event: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox events: %w", err)
	}

	return result, nil
}
