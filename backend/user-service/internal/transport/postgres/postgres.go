package postgres

import (
	"context"
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg settings.PostgresSettings) (*PostgresClient, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	pgCfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		names := []string{
			"v1_user", "_v1_user",
			"v1_profile", "_v1_profile",
			"v1_status", "_v1_status",
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
