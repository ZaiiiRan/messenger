package postgresimpl

import (
	"context"
	"fmt"

	userphoto "github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user_photo"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userPhotoRepository struct {
	conn *pgxpool.Conn
}

func NewUserPhotoRepository(conn *pgxpool.Conn) interfaces.UserPhotoRepository {
	return &userPhotoRepository{conn: conn}
}

func (r *userPhotoRepository) Create(ctx context.Context, photo *userphoto.UserPhoto) (*userphoto.UserPhoto, error) {
	row := r.conn.QueryRow(ctx, `
		INSERT INTO user_photos (user_id, file_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, file_id, created_at, updated_at
	`, photo.GetUserID(), photo.GetFileID(), photo.GetCreatedAt(), photo.GetUpdatedAt())

	var d models.V1UserPhotoDal
	if err := row.Scan(&d.Id, &d.UserId, &d.FileId, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, err
	}
	return d.ToDomain(), nil
}

func (r *userPhotoRepository) GetByUserID(ctx context.Context, userID string, limit, offset int32) ([]*userphoto.UserPhoto, int64, error) {
	var total int64
	if err := r.conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_photos WHERE user_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.conn.Query(ctx, `
		SELECT id, user_id, file_id, created_at, updated_at
		FROM user_photos
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var photos []*userphoto.UserPhoto
	for rows.Next() {
		var d models.V1UserPhotoDal
		if err := rows.Scan(&d.Id, &d.UserId, &d.FileId, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		photos = append(photos, d.ToDomain())
	}
	return photos, total, rows.Err()
}

func (r *userPhotoRepository) GetLatestByUserID(ctx context.Context, userID string) (*userphoto.UserPhoto, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT id, user_id, file_id, created_at, updated_at
		FROM user_photos
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, userID)

	var d models.V1UserPhotoDal
	if err := row.Scan(&d.Id, &d.UserId, &d.FileId, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, nil
	}
	return d.ToDomain(), nil
}

func (r *userPhotoRepository) GetLatestByUserIDs(ctx context.Context, userIDs []string) (map[string]*userphoto.UserPhoto, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	rows, err := r.conn.Query(ctx, `
		SELECT DISTINCT ON (user_id) id, user_id, file_id, created_at, updated_at
		FROM user_photos
		WHERE user_id = ANY($1)
		ORDER BY user_id, created_at DESC
	`, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*userphoto.UserPhoto, len(userIDs))
	for rows.Next() {
		var d models.V1UserPhotoDal
		if err := rows.Scan(&d.Id, &d.UserId, &d.FileId, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		result[d.UserId] = d.ToDomain()
	}
	return result, rows.Err()
}

func (r *userPhotoRepository) GetFileIDsByUserIDs(ctx context.Context, userIDs []string) (map[string][]string, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	rows, err := r.conn.Query(ctx, `
		SELECT user_id, file_id
		FROM user_photos
		WHERE user_id = ANY($1)
		ORDER BY user_id, created_at DESC
	`, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]string)
	for rows.Next() {
		var userID, fileID string
		if err := rows.Scan(&userID, &fileID); err != nil {
			return nil, err
		}
		result[userID] = append(result[userID], fileID)
	}
	return result, rows.Err()
}

func (r *userPhotoRepository) GetByID(ctx context.Context, id string) (*userphoto.UserPhoto, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT id, user_id, file_id, created_at, updated_at
		FROM user_photos WHERE id = $1
	`, id)

	var d models.V1UserPhotoDal
	if err := row.Scan(&d.Id, &d.UserId, &d.FileId, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return nil, nil
	}
	return d.ToDomain(), nil
}

func (r *userPhotoRepository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.conn.Exec(ctx, `DELETE FROM user_photos WHERE id = $1`, id)
	return err
}

func (r *userPhotoRepository) DeleteByUserIDs(ctx context.Context, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	_, err := r.conn.Exec(ctx,
		fmt.Sprintf(`DELETE FROM user_photos WHERE user_id = ANY($1)`),
		userIDs,
	)
	return err
}
