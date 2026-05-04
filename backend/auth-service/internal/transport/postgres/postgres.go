package postgres

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg settings.PostgresSettings) (*PostgresClient, error) {
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
	
	pgCfg, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	pgCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		names := []string{
			"v1_password", "_v1_password",
			"v1_refresh_token", "_v1_refresh_token",
			"v1_user_version", "_v1_user_version",
			"v1_confirmation_code", "_v1_confirmation_code",
			"v1_password_reset_token", "_v1_password_reset_token",
		}

		types, err := conn.LoadTypes(ctx, names)
		if err != nil {
			return fmt.Errorf("load types: %w", err)
		}
		conn.TypeMap().RegisterTypes(types)
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &PostgresClient{pool: pool}, nil
}

func (p *PostgresClient) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	return p.pool.Acquire(ctx)
}

func (p *PostgresClient) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}
