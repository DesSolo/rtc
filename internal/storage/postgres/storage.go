package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage ...
type Storage struct {
	manager *manager
}

// NewStorage ...
func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		manager: newManager(pool),
	}
}

// Close ...
func (s *Storage) Close() error {
	s.manager.Close()
	return nil
}

// WithTransaction ...
func (s *Storage) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	return s.manager.WithTransaction(ctx, f)
}
