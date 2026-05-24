package postgresimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	filevariant "github.com/ZaiiiRan/messenger/backend/file-service/internal/domain/file_variant"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileVariantRepository struct {
	conn *pgxpool.Conn
}

func NewFileVariantRepository(conn *pgxpool.Conn) interfaces.FileVariantRepository {
	return &FileVariantRepository{
		conn: conn,
	}
}

func (r *FileVariantRepository) CreateFileVariants(ctx context.Context, fileVariants []*filevariant.FileVariant) error {
	if len(fileVariants) == 0 {
		return nil
	}

	fileIds := make([]string, len(fileVariants))
	variants := make([]int16, len(fileVariants))
	buckets := make([]int16, len(fileVariants))
	s3Keys := make([]string, len(fileVariants))
	mimeTypes := make([]string, len(fileVariants))
	sizes := make([]int64, len(fileVariants))
	metas := make([]json.RawMessage, len(fileVariants))
	createdAts := make([]time.Time, len(fileVariants))
	updatedAts := make([]time.Time, len(fileVariants))
	for i, fileVariant := range fileVariants {
		fileIds[i] = fileVariant.GetFileID()
		variants[i] = int16(fileVariant.GetVariant())
		buckets[i] = int16(fileVariant.GetBucket())
		s3Keys[i] = fileVariant.GetS3Key()
		mimeTypes[i] = fileVariant.GetMimeType()
		sizes[i] = fileVariant.GetSize()

		metaBytes, err := json.Marshal(fileVariant.GetMeta())
		if err != nil {
			return fmt.Errorf("marshal meta: %w", err)
		}
		metas[i] = metaBytes

		createdAts[i] = fileVariant.GetCreatedAt()
		updatedAts[i] = fileVariant.GetUpdatedAt()
	}

	const sql = `
		INSERT INTO file_variants (file_id, variant, bucket, s3_key, mime_type, size, meta, created_at, updated_at)
		SELECT f.file_id, f.variant, f.bucket, f.s3_key, f.mime_type, f.size, f.meta, f.created_at, f.updated_at
		FROM UNNEST($1::uuid[], $2::smallint[], $3::smallint[], $4::text[], $5::text[], $6::bigint[], $7::jsonb[], $8::timestamptz[], $9::timestamptz[]) 
			AS f(file_id, variant, bucket, s3_key, mime_type, size, meta, created_at, updated_at)
		RETURNING id, file_id, variant, bucket, s3_key, mime_type, size, meta, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, fileIds, variants, buckets, s3Keys, mimeTypes, sizes, metas, createdAts, updatedAts)
	if err != nil {
		return fmt.Errorf("insert file variants: %w", err)
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var res models.V1FileVariantDal
		if err := rows.Scan(&res.Id, &res.FileId, &res.Variant, &res.Bucket, &res.S3Key, &res.MimeType, &res.Size, &res.Meta, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan inserted file variant: %w", err)
		}
		resDomain, err := res.ToDomain()
		if err != nil {
			return fmt.Errorf("convert inserted file variant to domain: %w", err)
		}
		fileVariants[i] = resDomain
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate inserted file variants: %w", err)
	}
	return nil
}

func (r *FileVariantRepository) UpdateFileVariants(ctx context.Context, fileVariants []*filevariant.FileVariant) error {
	if len(fileVariants) == 0 {
		return nil
	}

	fileVariantDals := make([]models.V1FileVariantDal, len(fileVariants))
	for i, fileVariant := range fileVariants {
		fileVariantDal, err := models.V1FileVariantFromDomain(fileVariant)
		if err != nil {
			return fmt.Errorf("convert file variant to dal: %w", err)
		}
		fileVariantDals[i] = fileVariantDal
	}

	const sql = `
		UPDATE file_variants t
		SET
			variant = f.variant,
			bucket = f.bucket,
			s3_key = f.s3_key,
			mime_type = f.mime_type,
			size = f.size,
			meta = f.meta,
			updated_at = f.updated_at
		FROM UNNEST($1::v1_file_variant[]) AS f
		RETURNING t.id, t.file_id, t.variant, t.bucket, t.s3_key, t.mime_type, t.size, t.meta, t.created_at, t.updated_at
	`

	rows, err := r.conn.Query(ctx, sql, fileVariantDals)
	if err != nil {
		return fmt.Errorf("update file variants: %w", err)
	}
	defer rows.Close()

	fileVariantById := make(map[string]models.V1FileVariantDal)
	for rows.Next() {
		var res models.V1FileVariantDal
		if err := rows.Scan(&res.Id, &res.FileId, &res.Variant, &res.Bucket, &res.S3Key, &res.MimeType, &res.Size, &res.Meta, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return fmt.Errorf("scan updated file variant: %w", err)
		}
		fileVariantById[res.Id] = res
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate updated file variants: %w", err)
	}
	for _, fileVariant := range fileVariants {
		id := fileVariant.GetID()
		fileVariantDal := fileVariantById[id]
		fileVariantDomain, err := fileVariantDal.ToDomain()
		if err != nil {
			return fmt.Errorf("convert updated file variant to domain: %w", err)
		}
		*fileVariant = *fileVariantDomain
	}
	return nil
}

func (r *FileVariantRepository) DeleteFileVariants(ctx context.Context, fileVariants []*filevariant.FileVariant) error {
	if len(fileVariants) == 0 {
		return nil
	}

	ids := make([]string, 0, len(fileVariants))
	for _, fileVariant := range fileVariants {
		ids = append(ids, fileVariant.GetID())
	}

	const sql = `
		DELETE FROM file_variants t
		WHERE t.id = ANY($1)
	`

	if _, err := r.conn.Exec(ctx, sql, ids); err != nil {
		return fmt.Errorf("delete file variants: %w", err)
	}
	return nil
}

func (r *FileVariantRepository) QueryFileVariants(ctx context.Context, query *models.QueryFileVariantsDal, locked bool) ([]*filevariant.FileVariant, error) {
	if query == nil {
		return nil, nil
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT t.id, t.file_id, t.variant, t.bucket, t.s3_key, t.mime_type, t.size, t.meta, t.created_at, t.updated_at
		FROM file_variants t
		WHERE 1=1
	`)

	appendAnyEqual(&sb, "t.id", query.Filter.Ids, &args, &argPos)
	appendAnyEqual(&sb, "t.file_id", query.Filter.FileIds, &args, &argPos)
	appendAnyEqual(&sb, "t.variant", query.Filter.Variants, &args, &argPos)
	appendAnyEqual(&sb, "t.s3_key", query.Filter.S3Keys, &args, &argPos)
	appendRange(&sb, "t.created_at", query.Filter.CreatedFrom, query.Filter.CreatedTo, &args, &argPos)
	appendRange(&sb, "t.updated_at", query.Filter.UpdatedFrom, query.Filter.UpdatedTo, &args, &argPos)

	if query.Limit != nil {
		fmt.Fprintf(&sb, " LIMIT $%d", argPos)
		argPos++
		args = append(args, *query.Limit)
		if query.Offset != nil {
			fmt.Fprintf(&sb, " OFFSET $%d", argPos)
			argPos++
			args = append(args, *query.Offset)
		}
	}

	if locked {
		sb.WriteString(` FOR UPDATE`)
	}

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query file variants: %w", err)
	}
	defer rows.Close()

	var result []*filevariant.FileVariant
	for rows.Next() {
		var res models.V1FileVariantDal
		if err := rows.Scan(&res.Id, &res.FileId, &res.Variant, &res.Bucket, &res.S3Key, &res.MimeType, &res.Size, &res.Meta, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan file variant: %w", err)
		}
		fileVariantDomain, err := res.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("convert file variant to domain: %w", err)
		}
		result = append(result, fileVariantDomain)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate file variants: %w", err)
	}
	return result, nil
}
