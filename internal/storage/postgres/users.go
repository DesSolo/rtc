package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"rtc/internal/storage"
)

// Users ...
func (s *Storage) Users(ctx context.Context, q string, limit, offset uint64) ([]*storage.User, uint64, error) {
	query := queryBuilder().
		Select("id, username, password_hash, is_enabled, roles, created_at, COUNT(*) OVER()").
		From("users").
		Limit(limit).
		Offset(offset).
		OrderBy("id DESC")

	if q != "" {
		query = queryLike(query, "username", q)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("query: %w", err)
	}

	rows, err := s.manager.Conn(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var (
		users []*storage.User
		total uint64
	)

	for rows.Next() {
		var user storage.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.IsEnabled, &user.Roles, &user.CreatedAt, &total); err != nil {
			return nil, 0, fmt.Errorf("rows.Scan: %w", err)
		}

		users = append(users, &user)
	}

	return users, total, nil
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
	query := "INSERT INTO users (username, password_hash, is_enabled, roles) VALUES ($1, $2, $3, $4) RETURNING id, created_at"

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, user.Username, user.PasswordHash, user.IsEnabled, user.Roles).Scan(&user.ID, &user.CreatedAt); err != nil {
		if isAlreadyExistsError(err) {
			return storage.ErrAlreadyExists
		}

		return fmt.Errorf("row.Scan: %w", err)
	}

	return nil
}

// UpdateUser ...
func (s *Storage) UpdateUser(ctx context.Context, id uint64, user *storage.User) error {
	query := "UPDATE users SET password_hash=$1, is_enabled=$2, roles=$3 WHERE id=$4"

	if _, err := s.manager.Conn(ctx).Exec(ctx, query, user.PasswordHash, user.IsEnabled, user.Roles, id); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
