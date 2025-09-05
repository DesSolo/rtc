package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"rtc/internal/storage"
)

// Configs ...
func (s *Storage) Configs(ctx context.Context, projectName, envName, releaseName string) ([]*storage.Config, error) {
	query := `
		SELECT c.id, c.release_id, c.key, c.value_type, c.metadata, c.created_at, c.updated_at FROM configs c
		JOIN releases r ON r.id = c.release_id
		JOIN environments e ON e.id = r.environment_id
		JOIN projects p ON p.id = e.project_id
		WHERE p.name = $1 AND e.name = $2 AND r.name = $3
	`

	rows, err := s.manager.Conn(ctx).Query(ctx, query, projectName, envName, releaseName)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var configs []*storage.Config

	for rows.Next() {
		var config storage.Config
		if err := rows.Scan(&config.ID, &config.ReleaseID, &config.Key, &config.ValueType, &config.Metadata, &config.CreatedAt, &config.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

// ConfigsByKeys ...
func (s *Storage) ConfigsByKeys(ctx context.Context, projectName, envName, releaseName string, keys []string) ([]*storage.Config, error) {
	query := `
		SELECT c.id, c.release_id, c.key, c.value_type, c.metadata, c.created_at, c.updated_at FROM configs c
		JOIN releases r ON r.id = c.release_id
		JOIN environments e ON e.id = r.environment_id
		JOIN projects p ON p.id = e.project_id
		WHERE p.name = $1 AND e.name = $2 AND r.name = $3 AND c.key = ANY ($4)
	`

	rows, err := s.manager.Conn(ctx).Query(ctx, query, projectName, envName, releaseName, keys)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var configs []*storage.Config

	for rows.Next() {
		var config storage.Config
		if err := rows.Scan(&config.ID, &config.ReleaseID, &config.Key, &config.ValueType, &config.Metadata, &config.CreatedAt, &config.UpdatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		configs = append(configs, &config)
	}

	return configs, nil
}

// Config ...
func (s *Storage) Config(ctx context.Context, projectName, envName, releaseName, key string) (*storage.Config, error) {
	query := `
		SELECT c.id, c.release_id, c.key, c.value_type, c.metadata, c.created_at, c.updated_at FROM configs c
		JOIN releases r ON r.id = c.release_id
		JOIN environments e ON e.id = r.environment_id
		JOIN projects p ON p.id = e.project_id
		WHERE p.name = $1 AND e.name = $2 AND r.name = $3 AND c.key = $4
	`

	var config storage.Config

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, projectName, envName, releaseName, key).Scan(&config.ID, &config.ReleaseID, &config.Key, &config.ValueType, &config.Metadata, &config.CreatedAt, &config.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("pool.Query: %w", err)
	}

	return &config, nil
}

// UpsertConfigs ...
func (s *Storage) UpsertConfigs(ctx context.Context, configs []*storage.Config) error {
	query := `
		INSERT INTO configs (release_id, key, value_type, metadata) VALUES ($1, $2, $3, $4)
		ON CONFLICT (release_id, key) DO UPDATE SET value_type = EXCLUDED.value_type, metadata = EXCLUDED.metadata, updated_at = NOW()
	`

	var batch pgx.Batch

	for _, config := range configs {
		batch.Queue(query, config.ReleaseID, config.Key, config.ValueType, config.Metadata)
	}

	results := s.manager.Conn(ctx).SendBatch(ctx, &batch)
	defer results.Close()

	for range configs {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("results.Exec: %w", err)
		}
	}

	return nil
}

// MarkConfigsUpdated ...
func (s *Storage) MarkConfigsUpdated(ctx context.Context, IDs []uint64) error {
	query := "UPDATE configs SET updated_at = NOW() WHERE id = ANY($1)"

	if _, err := s.manager.Conn(ctx).Exec(ctx, query, IDs); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}

// DeleteConfigs ...
func (s *Storage) DeleteConfigs(ctx context.Context, IDs []uint64) error {
	query := "DELETE FROM configs WHERE id = ANY ($1)"

	if _, err := s.manager.Conn(ctx).Exec(ctx, query, IDs); err != nil {
		return fmt.Errorf("pool.Exec: %w", err)
	}

	return nil
}
