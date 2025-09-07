package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"rtc/internal/storage"
)

// Projects ...
func (s *Storage) Projects(ctx context.Context, q string, limit, offset uint64) ([]*storage.Project, uint64, error) {
	query := queryBuilder().
		Select("id, name, description, created_at, COUNT(*) OVER() AS total").
		From("projects").
		Limit(limit).
		Offset(offset).
		OrderBy("id DESC")

	if q != "" {
		query = queryLike(query, "name", q)
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

	return scanProjectsWithTotal(rows) // nolint:errcheck
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

func scanProjectsWithTotal(rows pgx.Rows) ([]*storage.Project, uint64, error) {
	var (
		projects []*storage.Project
		total    uint64
	)

	for rows.Next() {
		var project storage.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &total); err != nil {
			return nil, 0, fmt.Errorf("rows.Scan: %w", err)
		}

		projects = append(projects, &project)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows.Err: %w", err)
	}

	return projects, total, nil
}
