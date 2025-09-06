package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"rtc/internal/storage"
)

// Users ...
func (s *Storage) Users(ctx context.Context, limit, offset int) ([]*storage.User, error) {
	query := "SELECT id, username, password_hash, is_enabled, roles, created_at FROM users LIMIT $1 OFFSET $2"

	rows, err := s.manager.Conn(ctx).Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var users []*storage.User

	for rows.Next() {
		var user storage.User
		if err := rows.Scan(&user.ID, &user.PasswordHash, &user.IsEnabled, &user.Roles, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil

}

// User ...
func (s *Storage) User(ctx context.Context, username string) (*storage.User, error) {
	query := "SELECT id, username, password_hash, is_enabled, roles, created_at FROM users WHERE username = $1"

	var user storage.User

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.IsEnabled, &user.Roles, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("pool.Query: %w", err)
	}

	return &user, nil
}

// CreateUser ...
func (s *Storage) CreateUser(ctx context.Context, user *storage.User) error {
	query := "INSERT INTO users (username, password_hash, is_enabled, roles) VALUES ($1, $2, $3) RETURNING id, created_at"

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, user.Username, user.PasswordHash, user.IsEnabled, user.Roles).Scan(&user.ID, &user.CreatedAt); err != nil {
		if isAlreadyExistsError(err) {
			return storage.ErrAlreadyExists
		}

		return fmt.Errorf("row.Scan: %w", err)
	}

	return nil
}
