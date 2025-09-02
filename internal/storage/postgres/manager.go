package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type manager struct {
	pool *pgxpool.Pool
}

func newManager(pool *pgxpool.Pool) *manager {
	return &manager{pool: pool}
}

func (m *manager) Conn(ctx context.Context) conn {
	tx := txFromContext(ctx)
	if tx != nil {
		return tx
	}

	return m.pool
}

func (m *manager) Close() {
	m.pool.Close()
}

func (m *manager) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("pool.BeginTx: %w", err)
	}

	if err := f(txToContext(ctx, tx)); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			return fmt.Errorf("tx.Rollback: %w", errRollback)
		}

		return fmt.Errorf("f: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}

	return nil
}

type conn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type txContextKey struct{}

func txToContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

func txFromContext(ctx context.Context) pgx.Tx {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	if !ok {
		return nil
	}

	return tx
}
