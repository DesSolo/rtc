package postgres

import (
	"context"
	"fmt"

	"rtc/internal/storage"
)

// Projects ...
func (s *Storage) Projects(ctx context.Context, limit, offset int) ([]*storage.Project, error) {
	query := "SELECT id, name, description, created_at FROM projects ORDER BY id DESC LIMIT $1 OFFSET $2"

	rows, err := s.manager.Conn(ctx).Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var projects []*storage.Project

	for rows.Next() {
		var project storage.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return projects, nil
}

// SearchProjects ...
func (s *Storage) SearchProjects(ctx context.Context, q string, limit int) ([]*storage.Project, error) {
	query := "SELECT id, name, description FROM projects WHERE name LIKE '%' || $1 || '%' LIMIT $2"

	var projects []*storage.Project
	rows, err := s.manager.Conn(ctx).Query(ctx, query, q, limit)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project storage.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	return projects, nil
}

// ProjectByName ...
func (s *Storage) ProjectByName(ctx context.Context, name string) (*storage.Project, error) {
	query := "SELECT id, name, description, created_at FROM projects WHERE name = $1"

	var project storage.Project
	if err := s.manager.Conn(ctx).QueryRow(ctx, query, name).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt); err != nil {
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}

	return &project, nil
}

// CreateProject ...
func (s *Storage) CreateProject(ctx context.Context, project *storage.Project) error {
	query := "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id, created_at"

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, project.Name, project.Description).Scan(&project.ID, &project.CreatedAt); err != nil {
		if isAlreadyExistsError(err) {
			return storage.ErrAlreadyExists
		}

		return fmt.Errorf("row.Scan: %w", err)
	}

	return nil
}
