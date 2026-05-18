package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
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
			"v1_user", "_v1_user",
			"v1_profile", "_v1_profile",
			"v1_status", "_v1_status",
			"v1_inbox_outbox_event", "_v1_inbox_outbox_event",
			"v1_privacy_settings", "_v1_privacy_settings",
		}

		types, err := conn.LoadTypes(ctx, names)
		if err != nil {
			return fmt.Errorf("load types: %w", err)
		}
		conn.TypeMap().RegisterTypes(types)
		return nil
	}

	pgCfg.MinConns = int32(cfg.MinConns)
	pgCfg.MaxConns = int32(cfg.MaxConns)
	pgCfg.MinIdleConns = int32(cfg.MinIdleConns)
	pgCfg.MaxConnIdleTime = time.Duration(cfg.MaxConnIdleTime) * time.Second
	pgCfg.MaxConnLifetime = time.Duration(cfg.MaxConnLifetime) * time.Second
	pgCfg.PingTimeout = time.Duration(cfg.PingTimeout) * time.Second
	pgCfg.HealthCheckPeriod = time.Duration(cfg.HealthCheckPeriod) * time.Second

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
