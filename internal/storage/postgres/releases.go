package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/DesSolo/rtc/internal/storage"
)

// Releases ...
func (s *Storage) Releases(ctx context.Context, projectName, envName string) ([]*storage.Release, error) {
	query := `
		SELECT r.id, r.environment_id, r.name, r.created_at FROM releases r
		JOIN environments e ON e.id = r.environment_id
		JOIN projects p ON p.id = e.project_id
		WHERE p.name = $1 AND e.name = $2
	`

	rows, err := s.manager.Conn(ctx).Query(ctx, query, projectName, envName)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var releases []*storage.Release

	for rows.Next() {
		var release storage.Release
		if err := rows.Scan(&release.ID, &release.EnvironmentID, &release.Name, &release.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		releases = append(releases, &release)
	}

	return releases, nil
}

// Release ...
func (s *Storage) Release(ctx context.Context, envID uint64, releaseName string) (*storage.Release, error) {
	query := "SELECT id, environment_id, name, created_at FROM releases WHERE environment_id = $1 AND name = $2"

	var release storage.Release
	if err := s.manager.Conn(ctx).QueryRow(ctx, query, envID, releaseName).Scan(&release.ID, &release.EnvironmentID, &release.Name, &release.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("pool.QueryRow: %w", err)
	}

	return &release, nil
}

// CreateRelease ...
func (s *Storage) CreateRelease(ctx context.Context, release *storage.Release) error {
	query := "INSERT INTO releases (environment_id, name) VALUES ($1, $2) RETURNING id, created_at"

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, release.EnvironmentID, release.Name).Scan(&release.ID, &release.CreatedAt); err != nil {
		if isAlreadyExistsError(err) {
			return storage.ErrAlreadyExists
		}

		return fmt.Errorf("row.Scan: %w", err)
	}

	return nil
}

// DeleteRelease ...
func (s *Storage) DeleteRelease(ctx context.Context, envID uint64, releaseName string) error {
	query := "DELETE FROM releases WHERE environment_id = $1 AND name = $2"

	if _, err := s.manager.Conn(ctx).Exec(ctx, query, envID, releaseName); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
