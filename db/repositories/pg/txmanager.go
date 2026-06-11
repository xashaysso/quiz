package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgTxManager struct {
	Pool *pgxpool.Pool
}

func NewPgTxManager(pool *pgxpool.Pool) *PgTxManager {
	return &PgTxManager{Pool: pool}
}

func (m *PgTxManager) WithinTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
