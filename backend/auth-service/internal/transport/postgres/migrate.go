package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Migrate(ctx context.Context, cfg settings.PostgresSettings) error {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		cfg.User,
		cfg.Password,
		cfg.Address,
		cfg.Database,
	)
	if cfg.Options != "" {
		connectionString = fmt.Sprintf("%s?%s", connectionString, cfg.Options)
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	goose.SetTableName("db_version")
	if err := goose.Up(db, "./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
