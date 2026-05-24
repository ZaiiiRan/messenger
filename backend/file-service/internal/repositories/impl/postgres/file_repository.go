package postgresimpl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	conn *pgxpool.Conn
}

func NewFileRepository(conn *pgxpool.Conn) interfaces.FileRepository {
	return &FileRepository{
		conn: conn,
	}
}

func (r *FileRepository) CreateFiles(ctx context.Context, files []*file.File) error {
	if len(files) == 0 {
		return nil
	}

	ownerIds := make([]string, len(files))
	statuses := make([]int16, len(files))
	createadAts := make([]time.Time, len(files))
	updatesAts := make([]time.Time, len(files))
	for i, file := range files {
		ownerIds[i] = file.GetOwnerID()
		statuses[i] = int16(file.GetStatus())
		createadAts[i] = file.GetCreatedAt()
		updatesAts[i] = file.GetUpdatedAt()
	}

	const sql = `
		INSERT INTO files (owner_id, status, created_at, updated_at)
		SELECT f.owner_id, f.status, f.created_at, f.updated_at
		FROM UNNEST($1::uuid[], $2::smallint[], $3::timestamptz[], $4::timestamptz[])
			AS f(owner_id, status, created_at, updated_at)
		RETURNING id, owner_id, status, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, ownerIds, statuses, createadAts, updatesAts)
	if err != nil {
		return fmt.Errorf("insert files: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var res models.V1FileDal
		if err := rows.Scan(&res.Id, &res.OwnerId, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan inserted file: %w", err)
		}
		files[i] = res.ToDomain()
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate inserted files: %w", err)
	}
	return nil
}

func (r *FileRepository) UpdateFiles(ctx context.Context, files []*file.File) error {
	if len(files) == 0 {
		return nil
	}

	fileDals := make([]models.V1FileDal, len(files))
	for i, file := range files {
		fileDals[i] = models.V1FileFromDomain(file)
	}

	const sql = `
		UPDATE files AS t
		SET
			status = f.status,
			updated_at = f.updated_at
		FROM UNNEST($1::v1_file[]) AS f
		WHERE t.id = f.id
		RETURNING t.id, t.owner_id, t.status, t.created_at, t.updated_at
	`

	rows, err := r.conn.Query(ctx, sql, fileDals)
	if err != nil {
		return fmt.Errorf("update files: %w", err)
	}
	defer rows.Close()

	fileById := make(map[string]models.V1FileDal)
	for rows.Next() {
		var res models.V1FileDal
		if err := rows.Scan(&res.Id, &res.OwnerId, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan updated file: %w", err)
		}
		fileById[res.Id] = res
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate updated files: %w", err)
	}
	for _, file := range files {
		id := file.GetID()
		*file = *fileById[id].ToDomain()
	}
	return nil
}

func (r *FileRepository) DeleteFiles(ctx context.Context, files []*file.File) error {
	if len(files) == 0 {
		return nil
	}

	ids := make([]string, 0, len(files))
	for _, file := range files {
		ids = append(ids, file.GetID())
	}

	const sql = `
		DELETE FROM files t
		WHERE t.id = ANY($1)
	`

	if _, err := r.conn.Exec(ctx, sql, ids); err != nil {
		return fmt.Errorf("delete files: %w", err)
	}
	return nil
}

func (r *FileRepository) QueryFiles(ctx context.Context, query *models.QueryFilesDal, locked bool, skipLocked bool) ([]*file.File, error) {
	if query == nil {
		return nil, nil
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT t.id, t.owner_id, t.status, t.created_at, t.updated_at
		FROM files t
		WHERE 1=1
	`)

	appendAnyEqual(&sb, "t.id", query.Filter.Ids, &args, &argPos)
	appendAnyEqual(&sb, "t.owner_id", query.Filter.OwnerIds, &args, &argPos)
	appendAnyEqual(&sb, "t.status", query.Filter.Statuses, &args, &argPos)
	appendRange(&sb, "t.created_at", query.Filter.CreatedFrom, query.Filter.CreatedTo, &args, &argPos)
	appendRange(&sb, "t.updated_at", query.Filter.UpdatedFrom, query.Filter.UpdatedTo, &args, &argPos)
	appendLimitOffset(&sb, query.Limit, query.Offset, &args, &argPos)

	if locked {
		sb.WriteString(` FOR UPDATE`)
		if skipLocked {
			sb.WriteString(` SKIP LOCKED`)
		}
	}

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query files: %w", err)
	}
	defer rows.Close()

	var result []*file.File
	for rows.Next() {
		var res models.V1FileDal
		if err := rows.Scan(&res.Id, &res.OwnerId, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan file: %w", err)
		}
		result = append(result, res.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate files: %w", err)
	}
	return result, nil
}
