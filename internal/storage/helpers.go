package storage

import (
	"context"
	"errors"
	"fmt"
)

// GetOrCreateEnvironment ...
func GetOrCreateEnvironment(ctx context.Context, s Storage, projectID uint64, name string) (*Environment, error) {
	environment, err := s.Environment(ctx, projectID, name)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("s.Environment: %w", err)
		}

		environment = &Environment{
			ProjectID: projectID,
			Name:      name,
		}
		if err := s.CreateEnvironment(ctx, environment); err != nil {
			return nil, fmt.Errorf("s.CreateEnvironment: %w", err)
		}

		return environment, nil
	}

	return environment, nil
}

// GetOrCreateRelease ...
func GetOrCreateRelease(ctx context.Context, s Storage, envID uint64, name string) (*Release, error) {
	release, err := s.Release(ctx, envID, name)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("s.Release: %w", err)
		}

		release = &Release{
			EnvironmentID: envID,
			Name:          name,
		}
		if err := s.CreateRelease(ctx, release); err != nil {
			return nil, fmt.Errorf("s.CreateRelease: %w", err)
		}

		return release, nil
	}

	return release, nil
}
