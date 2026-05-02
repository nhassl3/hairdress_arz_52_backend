package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(pool),
		pool:    pool,
	}
}

func (s *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed when acquiring db connection: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("failed when acquiring db transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	if err = fn(s.WithTx(tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err : %v, rbErr: %v", err, rbErr)
		}
		return fmt.Errorf("failed when executing db transaction: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) Pool() *pgxpool.Pool {
	return s.Pool()
}

// BeginTx starts a transaction and returns a Queries scoped to it.
// The caller is responsible for commiting or rolling back the transaction
func (s *Store) BeginTx(ctx context.Context) (*Queries, pgx.Tx, error) {
	tx, err := s.Pool().Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	return s.WithTx(tx), tx, nil
}
