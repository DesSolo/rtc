package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/DesSolo/rtc/internal/storage"
)

// Environments ...
func (s *Storage) Environments(ctx context.Context, projectName string) ([]*storage.Environment, error) {
	query := `
		SELECT e.id, e.project_id, e.name FROM environments e
		JOIN projects p ON p.id = e.project_id
		WHERE p.name = $1
	`

	rows, err := s.manager.Conn(ctx).Query(ctx, query, projectName)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var environments []*storage.Environment
	for rows.Next() {
		var environment storage.Environment
		if err := rows.Scan(&environment.ID, &environment.ProjectID, &environment.Name); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		environments = append(environments, &environment)
	}

	return environments, nil
}

// Environment ...
func (s *Storage) Environment(ctx context.Context, projectID uint64, envName string) (*storage.Environment, error) {
	query := "SELECT id, project_id, name FROM environments WHERE project_id = $1 AND name = $2"

	var environment storage.Environment
	if err := s.manager.Conn(ctx).QueryRow(ctx, query, projectID, envName).Scan(&environment.ID, &environment.ProjectID, &environment.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return &environment, nil
}

// CreateEnvironment ...
func (s *Storage) CreateEnvironment(ctx context.Context, env *storage.Environment) error {
	query := "INSERT INTO environments (project_id, name) VALUES ($1, $2) RETURNING id"

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, env.ProjectID, env.Name).Scan(&env.ID); err != nil {
		if isAlreadyExistsError(err) {
			return storage.ErrAlreadyExists
		}

		return fmt.Errorf("row.Scan: %w", err)
	}

	return nil
}
