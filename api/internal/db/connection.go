package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pannaraiSIH/stock-paper-trading/internal/db/queries"
)

type Store struct {
	*queries.Queries
	Pool *pgxpool.Pool
}

func NewStore(connStr string) (*Store, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	store := &Store{
		Queries: queries.New(pool),
		Pool:    pool,
	}

	if err := store.Ping(ctx); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) Ping(ctx context.Context) error {
	return s.Pool.Ping(ctx)
}

func (s *Store) Close() {
	s.Pool.Close()
}
