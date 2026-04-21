package postgresunitofwork

import (
	"context"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/transport/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	pgClient *postgres.PostgresClient
	conn     *pgxpool.Conn
	tx       pgx.Tx
	closed   bool
}

func New(pgClient *postgres.PostgresClient) *UnitOfWork {
	return &UnitOfWork{pgClient: pgClient}
}

func (u *UnitOfWork) GetConn(ctx context.Context) (*pgxpool.Conn, error) {
	if u.closed {
		return nil, pgx.ErrTxClosed
	}
	if u.conn != nil {
		return u.conn, nil
	}
	c, err := u.pgClient.GetConn(ctx)
	if err != nil {
		return nil, err
	}
	u.conn = c
	return u.conn, nil
}

func (u *UnitOfWork) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	if u.closed {
		return nil, pgx.ErrTxClosed
	}
	if u.tx != nil {
		return u.tx, nil
	}

	if u.conn == nil {
		c, err := u.pgClient.GetConn(ctx)
		if err != nil {
			return nil, err
		}
		u.conn = c
	}

	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	u.tx = tx
	return u.tx, nil
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Commit(ctx)
	u.tx = nil
	return err
}

func (u *UnitOfWork) Rollback(ctx context.Context) error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Rollback(ctx)
	u.tx = nil
	return err
}

func (u *UnitOfWork) Close() {
	if u.tx != nil {
		u.tx.Rollback(context.Background())
		u.tx = nil
	}
	if u.conn != nil {
		u.conn.Release()
		u.conn = nil
	}
	u.closed = true
}
