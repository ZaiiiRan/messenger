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

func (r *UserRepository) Create(ctx context.Context, users []*user.User) error {
	if len(users) == 0 {
		return nil
	}

	userDals := make([]models.V1UserDal, len(users))
	for i, u := range users {
		userDals[i] = models.V1UserDalFromDomain(u)
	}

	insertedUsers, err := r.insertUsers(ctx, userDals)
	if err != nil {
		return err
	}

	userByUsername := make(map[string]models.V1UserDal, len(insertedUsers))
	for _, u := range insertedUsers {
		userByUsername[u.Username] = u
	}

	profileDals := make([]models.V1ProfileDal, len(users))
	statusDals := make([]models.V1StatusDal, len(users))
	privacySettingsDals := make([]models.V1PrivacySettingsDal, len(users))
	for i, u := range users {
		inserted := userByUsername[u.GetUsername()]
		profileDals[i] = models.V1ProfileDalFromDomain(inserted.Id, u.GetProfile())
		statusDals[i] = models.V1StatusDalFromDomain(inserted.Id, u.GetStatus())
		privacySettingsDal, err := models.V1PrivacySettingsDalFromDomain(inserted.Id, u.GetPrivacySettings())
		if err != nil {
			return err
		}
		privacySettingsDals[i] = privacySettingsDal
	}

	insertedProfiles, err := r.insertProfiles(ctx, profileDals)
	if err != nil {
		return err
	}

	insertedStatuses, err := r.insertStatuses(ctx, statusDals)
	if err != nil {
		return err
	}

	insertedPrivacySettings, err := r.insertPrivacySettings(ctx, privacySettingsDals)
	if err != nil {
		return err
	}

	profileByUserId := make(map[string]models.V1ProfileDal, len(insertedProfiles))
	for _, p := range insertedProfiles {
		profileByUserId[p.UserId] = p
	}
	statusByUserId := make(map[string]models.V1StatusDal, len(insertedStatuses))
	for _, s := range insertedStatuses {
		statusByUserId[s.UserId] = s
	}
	privacySettingsByUserId := make(map[string]models.V1PrivacySettingsDal, len(insertedPrivacySettings))
	for _, ps := range insertedPrivacySettings {
		privacySettingsByUserId[ps.UserId] = ps
	}

	for _, u := range users {
		inserted := userByUsername[u.GetUsername()]
		p := profileByUserId[inserted.Id]
		s := statusByUserId[inserted.Id]
		ps := privacySettingsByUserId[inserted.Id]
		domainPs, err := ps.ToDomain()
		if err != nil {
			return err
		}
		*u = *inserted.ToDomain(p.ToDomain(), domainPs, s)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, users []*user.User) error {
	if len(users) == 0 {
		return nil
	}

	userDals := make([]models.V1UserDal, len(users))
	profileDals := make([]models.V1ProfileDal, len(users))
	statusDals := make([]models.V1StatusDal, len(users))
	privacySettingsDals := make([]models.V1PrivacySettingsDal, len(users))
	for i, u := range users {
		userDals[i] = models.V1UserDalFromDomain(u)
		profileDals[i] = models.V1ProfileDalFromDomain(u.GetID(), u.GetProfile())
		statusDals[i] = models.V1StatusDalFromDomain(u.GetID(), u.GetStatus())
		privacySettingsDal, err := models.V1PrivacySettingsDalFromDomain(u.GetID(), u.GetPrivacySettings())
		if err != nil {
			return err
		}
		privacySettingsDals[i] = privacySettingsDal
	}

	updatedUsers, err := r.updateUsers(ctx, userDals)
	if err != nil {
		return err
	}

	updatedProfiles, err := r.updateProfiles(ctx, profileDals)
	if err != nil {
		return err
	}

	updatedStatuses, err := r.updateStatuses(ctx, statusDals)
	if err != nil {
		return err
	}

	updatedPrivacySettings, err := r.updatePrivacySettings(ctx, privacySettingsDals)
	if err != nil {
		return err
	}

	userById := make(map[string]models.V1UserDal, len(updatedUsers))
	for _, u := range updatedUsers {
		userById[u.Id] = u
	}
	profileByUserId := make(map[string]models.V1ProfileDal, len(updatedProfiles))
	for _, p := range updatedProfiles {
		profileByUserId[p.UserId] = p
	}
	statusByUserId := make(map[string]models.V1StatusDal, len(updatedStatuses))
	for _, s := range updatedStatuses {
		statusByUserId[s.UserId] = s
	}
	privacySettingsByUserId := make(map[string]models.V1PrivacySettingsDal, len(updatedPrivacySettings))
	for _, ps := range updatedPrivacySettings {
		privacySettingsByUserId[ps.UserId] = ps
	}

	for _, u := range users {
		id := u.GetID()
		ps := privacySettingsByUserId[id]
		domainPs, err := ps.ToDomain()
		if err != nil {
			return err
		}
		*u = *userById[id].ToDomain(profileByUserId[id].ToDomain(), domainPs, statusByUserId[id])
	}

	return nil
}

func (r *UserRepository) Query(ctx context.Context, query *models.QueryUsersDal) ([]*user.User, error) {
	if query == nil {
		return nil, nil
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
			s.is_confirmed, s.is_permanently_banned, s.banned_until, s.is_deleted, s.deleted_at, s.is_permanently_deleted, s.old_email, s.email_updated_at,
			ps.settings
		FROM users u
		JOIN profile p ON p.user_id = u.id
		JOIN status  s ON s.user_id = u.id
		JOIN privacy_settings ps ON ps.user_id = u.id
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
	appendBool(&sb, "s.is_permanently_deleted", query.Filter.IsPermanentlyDeleted, &args, &argPos)
	appendRange(&sb, "s.banned_until", query.Filter.BannedUntilFrom, query.Filter.BannedUntilTo, &args, &argPos)
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
			userDal            models.V1UserDal
			profileDal         models.V1ProfileDal
			statusDal          models.V1StatusDal
			privacySettingsDal models.V1PrivacySettingsDal
		)
		if err := rows.Scan(
			&userDal.Id, &userDal.Username, &userDal.Email, &userDal.CreatedAt, &userDal.UpdatedAt,
			&profileDal.FirstName, &profileDal.LastName, &profileDal.Phone, &profileDal.Birthdate, &profileDal.Bio,
			&statusDal.IsConfirmed, &statusDal.IsPermanentlyBanned, &statusDal.BannedUntil, &statusDal.IsDeleted, &statusDal.DeletedAt,
			&statusDal.IsPermanentlyDeleted, &statusDal.OldEmail, &statusDal.EmailUpdatedAt,
			&privacySettingsDal.Settings,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		profileDal.UserId = userDal.Id
		statusDal.UserId = userDal.Id
		privacySettingsDal.UserId = userDal.Id

		domainPs, err := privacySettingsDal.ToDomain()
		if err != nil {
			return nil, err
		}

		result = append(result, userDal.ToDomain(profileDal.ToDomain(), domainPs, statusDal))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return result, nil
}

func (r *UserRepository) QueryLocked(ctx context.Context, query *models.QueryUsersDal) ([]*user.User, error) {
	if query == nil {
		return nil, nil
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
			s.is_confirmed, s.is_permanently_banned, s.banned_until, s.is_deleted, s.deleted_at, s.is_permanently_deleted, s.old_email, s.email_updated_at,
			ps.settings
		FROM users u
		JOIN profile p ON p.user_id = u.id
		JOIN status  s ON s.user_id = u.id
		JOIN privacy_settings ps ON ps.user_id = u.id
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
	appendBool(&sb, "s.is_permanently_deleted", query.Filter.IsPermanentlyDeleted, &args, &argPos)
	appendRange(&sb, "s.banned_until", query.Filter.BannedUntilFrom, query.Filter.BannedUntilTo, &args, &argPos)
	appendRange(&sb, "u.created_at", query.Filter.CreatedFrom, query.Filter.CreatedTo, &args, &argPos)
	appendRange(&sb, "u.updated_at", query.Filter.UpdatedFrom, query.Filter.UpdatedTo, &args, &argPos)
	appendOrder(&sb, "u.updated_at", true)
	appendLimitOffset(&sb, query.Limit, query.Offset, &args, &argPos)

	sb.WriteString(" FOR UPDATE OF u SKIP LOCKED")

	rows, err := r.conn.Query(ctx, sb.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var result []*user.User
	for rows.Next() {
		var (
			userDal            models.V1UserDal
			profileDal         models.V1ProfileDal
			statusDal          models.V1StatusDal
			privacySettingsDal models.V1PrivacySettingsDal
		)
		if err := rows.Scan(
			&userDal.Id, &userDal.Username, &userDal.Email, &userDal.CreatedAt, &userDal.UpdatedAt,
			&profileDal.FirstName, &profileDal.LastName, &profileDal.Phone, &profileDal.Birthdate, &profileDal.Bio,
			&statusDal.IsConfirmed, &statusDal.IsPermanentlyBanned, &statusDal.BannedUntil, &statusDal.IsDeleted, &statusDal.DeletedAt,
			&statusDal.IsPermanentlyDeleted, &statusDal.OldEmail, &statusDal.EmailUpdatedAt,
			&privacySettingsDal.Settings,
		); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		profileDal.UserId = userDal.Id
		statusDal.UserId = userDal.Id
		privacySettingsDal.UserId = userDal.Id

		domainPs, err := privacySettingsDal.ToDomain()
		if err != nil {
			return nil, err
		}

		result = append(result, userDal.ToDomain(profileDal.ToDomain(), domainPs, statusDal))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return result, nil
}

func (r *UserRepository) Delete(ctx context.Context, users []*user.User) error {
	if len(users) == 0 {
		return nil
	}

	ids := make([]string, 0, len(users))
	for _, u := range users {
		if u == nil {
			continue
		}
		ids = append(ids, u.GetID())
	}

	if len(ids) == 0 {
		return nil
	}

	const sql = `
		WITH deleted_profile AS (
			DELETE FROM profile
			WHERE user_id::text = ANY($1)
		),
		deleted_status AS (
			DELETE FROM status
			WHERE user_id::text = ANY($1)
		),
		deleted_privacy_settings AS (
			DELETE FROM privacy_settings
			WHERE user_id::text = ANY($1)
		)
		DELETE FROM users
		WHERE id::text = ANY($1)
	`

	if _, err := r.conn.Exec(ctx, sql, ids); err != nil {
		return fmt.Errorf("delete users: %w", err)
	}

	return nil
}

func (r *UserRepository) insertUsers(ctx context.Context, dals []models.V1UserDal) ([]models.V1UserDal, error) {
	usernames := make([]string, len(dals))
	emails := make([]string, len(dals))
	for i, d := range dals {
		usernames[i] = d.Username
		emails[i] = d.Email
	}

	const sql = `
		INSERT INTO users (username, email)
		SELECT u.username, u.email
		FROM UNNEST($1::text[], $2::text[]) AS u(username, email)
		RETURNING id, username, email, created_at, updated_at
	`

	rows, err := r.conn.Query(ctx, sql, usernames, emails)
	if err != nil {
		return nil, fmt.Errorf("insert users: %w", err)
	}
	defer rows.Close()

	var result []models.V1UserDal
	for rows.Next() {
		var res models.V1UserDal
		if err := rows.Scan(&res.Id, &res.Username, &res.Email, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan inserted user: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inserted users: %w", err)
	}

	return result, nil
}

func (r *UserRepository) insertProfiles(ctx context.Context, dals []models.V1ProfileDal) ([]models.V1ProfileDal, error) {
	const sql = `
		INSERT INTO profile (user_id, first_name, last_name, phone, birthdate, bio)
		SELECT (i).user_id, (i).first_name, (i).last_name, (i).phone, (i).birthdate, (i).bio
		FROM UNNEST($1::v1_profile[]) i
		RETURNING id, user_id, first_name, last_name, phone, birthdate, bio
	`

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("insert profiles: %w", err)
	}
	defer rows.Close()

	var result []models.V1ProfileDal
	for rows.Next() {
		var res models.V1ProfileDal
		if err := rows.Scan(&res.Id, &res.UserId, &res.FirstName, &res.LastName, &res.Phone, &res.Birthdate, &res.Bio); err != nil {
			return nil, fmt.Errorf("scan inserted profile: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inserted profiles: %w", err)
	}

	return result, nil
}

func (r *UserRepository) insertStatuses(ctx context.Context, dals []models.V1StatusDal) ([]models.V1StatusDal, error) {
	const sql = `
		INSERT INTO status (user_id, is_confirmed, is_permanently_banned, banned_until, is_deleted, deleted_at, is_permanently_deleted, old_email, email_updated_at)
		SELECT (i).user_id, (i).is_confirmed, (i).is_permanently_banned, (i).banned_until, (i).is_deleted, (i).deleted_at, (i).is_permanently_deleted, (i).old_email, (i).email_updated_at
		FROM UNNEST($1::v1_status[]) i
		RETURNING id, user_id, is_confirmed, is_permanently_banned, banned_until, is_deleted, deleted_at, is_permanently_deleted, old_email, email_updated_at
	`

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("insert statuses: %w", err)
	}
	defer rows.Close()

	var result []models.V1StatusDal
	for rows.Next() {
		var res models.V1StatusDal
		if err := rows.Scan(
			&res.Id, &res.UserId, &res.IsConfirmed, &res.IsPermanentlyBanned, &res.BannedUntil, &res.IsDeleted, &res.DeletedAt, &res.IsPermanentlyDeleted, &res.OldEmail, &res.EmailUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan inserted status: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inserted statuses: %w", err)
	}

	return result, nil
}

func (r *UserRepository) insertPrivacySettings(ctx context.Context, dals []models.V1PrivacySettingsDal) ([]models.V1PrivacySettingsDal, error) {
	const sql = `
		INSERT INTO privacy_settings (user_id, settings)
		SELECT (i).user_id, (i).settings
		FROM UNNEST($1::v1_privacy_settings[]) i
		RETURNING id, user_id, settings
	`

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("insert privacy settings: %w", err)
	}
	defer rows.Close()

	var result []models.V1PrivacySettingsDal
	for rows.Next() {
		var res models.V1PrivacySettingsDal
		if err := rows.Scan(&res.Id, &res.UserId, &res.Settings); err != nil {
			return nil, fmt.Errorf("scan inserted privacy settings: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate inserted privacy settings: %w", err)
	}

	return result, nil
}

func (r *UserRepository) updateUsers(ctx context.Context, dals []models.V1UserDal) ([]models.V1UserDal, error) {
	const sql = `
		UPDATE users AS t
		SET
			username   = u.username,
			email      = u.email,
			updated_at = u.updated_at
		FROM UNNEST($1::v1_user[]) AS u
		WHERE t.id = u.id
		RETURNING t.id, t.username, t.email, t.created_at, t.updated_at
	`

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("update users: %w", err)
	}
	defer rows.Close()

	var result []models.V1UserDal
	for rows.Next() {
		var res models.V1UserDal
		if err := rows.Scan(&res.Id, &res.Username, &res.Email, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan updated user: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate updated users: %w", err)
	}

	return result, nil
}

func (r *UserRepository) updateProfiles(ctx context.Context, dals []models.V1ProfileDal) ([]models.V1ProfileDal, error) {
	const sql = `
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

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("update profiles: %w", err)
	}
	defer rows.Close()

	var result []models.V1ProfileDal
	for rows.Next() {
		var res models.V1ProfileDal
		if err := rows.Scan(&res.Id, &res.UserId, &res.FirstName, &res.LastName, &res.Phone, &res.Birthdate, &res.Bio); err != nil {
			return nil, fmt.Errorf("scan updated profile: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate updated profiles: %w", err)
	}

	return result, nil
}

func (r *UserRepository) updateStatuses(ctx context.Context, dals []models.V1StatusDal) ([]models.V1StatusDal, error) {
	const sql = `
		UPDATE status AS t
		SET
			is_confirmed           = u.is_confirmed,
			is_permanently_banned  = u.is_permanently_banned,
			banned_until           = u.banned_until,
			is_deleted             = u.is_deleted,
			deleted_at             = u.deleted_at,
			is_permanently_deleted = u.is_permanently_deleted,
			old_email              = u.old_email,
			email_updated_at       = u.email_updated_at
		FROM UNNEST($1::v1_status[]) AS u
		WHERE t.user_id = u.user_id
		RETURNING t.id, t.user_id, t.is_confirmed, t.is_permanently_banned, t.banned_until, t.is_deleted, t.deleted_at, t.is_permanently_deleted, t.old_email, t.email_updated_at
	`

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("update statuses: %w", err)
	}
	defer rows.Close()

	var result []models.V1StatusDal
	for rows.Next() {
		var res models.V1StatusDal
		if err := rows.Scan(
			&res.Id, &res.UserId, &res.IsConfirmed, &res.IsPermanentlyBanned, &res.BannedUntil, &res.IsDeleted, &res.DeletedAt, &res.IsPermanentlyDeleted, &res.OldEmail, &res.EmailUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan updated status: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate updated statuses: %w", err)
	}

	return result, nil
}

func (r *UserRepository) updatePrivacySettings(ctx context.Context, dals []models.V1PrivacySettingsDal) ([]models.V1PrivacySettingsDal, error) {
	const sql = `
		UPDATE privacy_settings AS t
		SET
			settings = u.settings
		FROM UNNEST($1::v1_privacy_settings[]) AS u
		WHERE t.user_id = u.user_id
		RETURNING t.id, t.user_id, t.settings
	`

	rows, err := r.conn.Query(ctx, sql, dals)
	if err != nil {
		return nil, fmt.Errorf("update privacy settings: %w", err)
	}
	defer rows.Close()

	var result []models.V1PrivacySettingsDal
	for rows.Next() {
		var res models.V1PrivacySettingsDal
		if err := rows.Scan(&res.Id, &res.UserId, &res.Settings); err != nil {
			return nil, fmt.Errorf("scan updated privacy settings: %w", err)
		}
		result = append(result, res)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate updated privacy settings: %w", err)
	}

	return result, nil
}
