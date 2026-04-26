package postgresimpl

import (
	"context"
	"fmt"
	"strings"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/domain/user"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/interfaces"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/repositories/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	conn *pgxpool.Conn
}

func NewUserRepository(conn *pgxpool.Conn) interfaces.UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	userDal := models.V1UserDalFromDomain(u)

	insertedUser, err := r.insertUser(ctx, userDal)
	if err != nil {
		return err
	}

	profileDal := models.V1ProfileDalFromDomain(insertedUser.Id, u.GetProfile())
	insertedProfile, err := r.insertProfile(ctx, profileDal)
	if err != nil {
		return err
	}

	statusDal := models.V1StatusDalFromDomain(insertedUser.Id, u.GetStatus())
	insertedStatus, err := r.insertStatus(ctx, statusDal)
	if err != nil {
		return err
	}

	*u = *insertedUser.ToDomain(insertedProfile.ToDomain(), insertedStatus.ToDomain())
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	userDal := models.V1UserDalFromDomain(u)

	updatedUser, err := r.updateUser(ctx, userDal)
	if err != nil {
		return err
	}

	profileDal := models.V1ProfileDalFromDomain(updatedUser.Id, u.GetProfile())
	updatedProfile, err := r.updateProfile(ctx, profileDal)
	if err != nil {
		return err
	}

	statusDal := models.V1StatusDalFromDomain(updatedUser.Id, u.GetStatus())
	updatedStatus, err := r.updateStatus(ctx, statusDal)
	if err != nil {
		return err
	}

	*u = *updatedUser.ToDomain(updatedProfile.ToDomain(), updatedStatus.ToDomain())
	return nil
}

func (r *UserRepository) Query(ctx context.Context, query *models.QueryUsersDal) ([]*user.User, error) {
	if query == nil {
		query = &models.QueryUsersDal{}
	}

	var (
		sb     strings.Builder
		args   []any
		argPos = 1
	)

	sb.WriteString(`
		SELECT
			u.id, u.username, u.email, u.created_at, u.updated_at,
			p.first_name, p.last_name, p.phone, p.birthdate, p.bio,
			s.is_confirmed, s.is_permanently_banned, s.banned_until, s.is_deleted, s.deleted_at
		FROM users u
		JOIN profile p ON p.user_id = u.id
		JOIN status  s ON s.user_id = u.id
		WHERE 1=1
	`)

	appendAnyEqual(&sb, "u.id::text", query.Filter.Ids, &args, &argPos)
	appendAnyEqual(&sb, "u.username", query.Filter.Usernames, &args, &argPos)
	appendIPrefix(&sb, "u.username", query.Filter.PartialUsernames, &args, &argPos)
	appendAnyEqual(&sb, "u.email", query.Filter.Emails, &args, &argPos)
	appendIPrefix(&sb, "u.email", query.Filter.PartialEmails, &args, &argPos)
	appendAnyEqual(&sb, "p.phone", query.Filter.PhoneNumbers, &args, &argPos)
	appendPartialNames(&sb, "p.first_name", "p.last_name", query.Filter.PartialNames, &args, &argPos)
	appendBool(&sb, "s.is_confirmed", query.Filter.IsConfirmed, &args, &argPos)
	appendBool(&sb, "s.is_deleted", query.Filter.IsDeleted, &args, &argPos)
	appendBool(&sb, "s.is_permanently_banned", query.Filter.IsPermanentlyBanned, &args, &argPos)
	appendIsNotNull(&sb, "s.banned_until", query.Filter.IsTemporarilyBanned)
	appendRange(&sb, "s.deleted_at", query.Filter.DeletedFrom, query.Filter.DeletedTo, &args, &argPos)
	appendRange(&sb, "u.created_at", query.Filter.CreatedFrom, query.Filter.CreatedTo, &args, &argPos)
	appendRange(&sb, "u.updated_at", query.Filter.UpdatedFrom, query.Filter.UpdatedTo, &args, &argPos)
	appendOrder(&sb, "u.id", true)
	appendLimitOffset(&sb, query.Limit, query.Offset, &args, &argPos)

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var result []*user.User
	for rows.Next() {
		var (
			userDal    models.V1UserDal
			profileDal models.V1ProfileDal
			statusDal  models.V1StatusDal
		)
		if err := rows.Scan(
			&userDal.Id, &userDal.Username, &userDal.Email, &userDal.CreatedAt, &userDal.UpdatedAt,
			&profileDal.FirstName, &profileDal.LastName, &profileDal.Phone, &profileDal.Birthdate, &profileDal.Bio,
			&statusDal.IsConfirmed, &statusDal.IsPermanentlyBanned, &statusDal.BannedUntil, &statusDal.IsDeleted, &statusDal.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		profileDal.UserId = userDal.Id
		statusDal.UserId = userDal.Id
		result = append(result, userDal.ToDomain(profileDal.ToDomain(), statusDal.ToDomain()))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return result, nil
}

func (r *UserRepository) insertUser(ctx context.Context, dal models.V1UserDal) (models.V1UserDal, error) {
	sql := `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id, username, email, created_at, updated_at
	`
	var res models.V1UserDal
	if err := r.conn.QueryRow(ctx, sql, dal.Username, dal.Email).Scan(
		&res.Id,
		&res.Username,
		&res.Email,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return models.V1UserDal{}, fmt.Errorf("insert user: %w", err)
	}
	return res, nil
}

func (r *UserRepository) insertProfile(ctx context.Context, dal models.V1ProfileDal) (models.V1ProfileDal, error) {
	sql := `
		INSERT INTO profile (user_id, first_name, last_name, phone, birthdate, bio)
		SELECT (i).user_id, (i).first_name, (i).last_name, (i).phone, (i).birthdate, (i).bio
		FROM UNNEST($1::v1_profile[]) i
		RETURNING id, user_id, first_name, last_name, phone, birthdate, bio
	`
	var res models.V1ProfileDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1ProfileDal{dal}).Scan(
		&res.Id,
		&res.UserId,
		&res.FirstName,
		&res.LastName,
		&res.Phone,
		&res.Birthdate,
		&res.Bio,
	); err != nil {
		return models.V1ProfileDal{}, fmt.Errorf("insert profile: %w", err)
	}
	return res, nil
}

func (r *UserRepository) insertStatus(ctx context.Context, dal models.V1StatusDal) (models.V1StatusDal, error) {
	sql := `
		INSERT INTO status (user_id, is_confirmed, is_permanently_banned, banned_until, is_deleted, deleted_at)
		SELECT (i).user_id, (i).is_confirmed, (i).is_permanently_banned, (i).banned_until, (i).is_deleted, (i).deleted_at
		FROM UNNEST($1::v1_status[]) i
		RETURNING id, user_id, is_confirmed, is_permanently_banned, banned_until, is_deleted, deleted_at
	`
	var res models.V1StatusDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1StatusDal{dal}).Scan(
		&res.Id,
		&res.UserId,
		&res.IsConfirmed,
		&res.IsPermanentlyBanned,
		&res.BannedUntil,
		&res.IsDeleted,
		&res.DeletedAt,
	); err != nil {
		return models.V1StatusDal{}, fmt.Errorf("insert status: %w", err)
	}
	return res, nil
}

func (r *UserRepository) updateUser(ctx context.Context, dal models.V1UserDal) (models.V1UserDal, error) {
	sql := `
		UPDATE users AS t
		SET
			username   = u.username,
			email      = u.email,
			updated_at = u.updated_at
		FROM UNNEST($1::v1_user[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.username, t.email, t.created_at, t.updated_at
	`
	var res models.V1UserDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1UserDal{dal}).Scan(
		&res.Id,
		&res.Username,
		&res.Email,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return models.V1UserDal{}, fmt.Errorf("update user: %w", err)
	}
	return res, nil
}

func (r *UserRepository) updateProfile(ctx context.Context, dal models.V1ProfileDal) (models.V1ProfileDal, error) {
	sql := `
		UPDATE profile AS t
		SET
			first_name = u.first_name,
			last_name  = u.last_name,
			phone      = u.phone,
			birthdate  = u.birthdate,
			bio        = u.bio
		FROM UNNEST($1::v1_profile[]) AS u
		WHERE t.user_id = u.user_id
		RETURNING t.id, t.user_id, t.first_name, t.last_name, t.phone, t.birthdate, t.bio
	`
	var res models.V1ProfileDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1ProfileDal{dal}).Scan(
		&res.Id,
		&res.UserId,
		&res.FirstName,
		&res.LastName,
		&res.Phone,
		&res.Birthdate,
		&res.Bio,
	); err != nil {
		return models.V1ProfileDal{}, fmt.Errorf("update profile: %w", err)
	}
	return res, nil
}

func (r *UserRepository) updateStatus(ctx context.Context, dal models.V1StatusDal) (models.V1StatusDal, error) {
	sql := `
		UPDATE status AS t
		SET
			is_confirmed          = u.is_confirmed,
			is_permanently_banned = u.is_permanently_banned,
			banned_until          = u.banned_until,
			is_deleted            = u.is_deleted,
			deleted_at            = u.deleted_at
		FROM UNNEST($1::v1_status[]) AS u
		WHERE t.user_id = u.user_id
		RETURNING t.id, t.user_id, t.is_confirmed, t.is_permanently_banned, t.banned_until, t.is_deleted, t.deleted_at
	`
	var res models.V1StatusDal
	if err := r.conn.QueryRow(ctx, sql, []models.V1StatusDal{dal}).Scan(
		&res.Id,
		&res.UserId,
		&res.IsConfirmed,
		&res.IsPermanentlyBanned,
		&res.BannedUntil,
		&res.IsDeleted,
		&res.DeletedAt,
	); err != nil {
		return models.V1StatusDal{}, fmt.Errorf("update status: %w", err)
	}
	return res, nil
}
