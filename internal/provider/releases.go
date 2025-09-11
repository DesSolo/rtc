package provider

import (
	"context"
	"errors"
	"fmt"

	"rtc/internal/auth"
	"rtc/internal/models"
	"rtc/internal/storage"
)

// ListReleases ...
func (p *Provider) ListReleases(ctx context.Context, projectName, envName string) ([]*models.Release, error) {
	releases, err := p.storage.Releases(ctx, projectName, envName)
	if err != nil {
		return nil, fmt.Errorf("storage.Releases: %w", err)
	}

	return convertReleasesToModel(releases), nil
}

// DeleteRelease ...
func (p *Provider) DeleteRelease(ctx context.Context, projectName, envName, releaseName string) error {
	project, err := p.storage.ProjectByName(ctx, projectName)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("storage.ProjectByName: %w", err)
	}

	env, err := p.storage.Environment(ctx, project.ID, envName)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("storage.Environment: %w", err)
	}

	actor := auth.UsernameFromContext(ctx)
	auditRecord, err := encodeAuditRecordReleaseDeleted(actor, projectName, envName, releaseName)
	if err != nil {
		return fmt.Errorf("encodeAuditRecordReleaseDeleted: %w", err)
	}

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		if err := p.storage.DeleteRelease(ctx, env.ID, releaseName); err != nil {
			return fmt.Errorf("storage.DeleteRelease: %w", err)
		}

		if err := p.valuesStorage.DeleteValuesByPath(ctx, formatValuesStoragePath(projectName, envName, releaseName)); err != nil {
			return fmt.Errorf("storage.DeleteValues: %w", err)
		}

		if err := p.storage.AddAuditRecord(ctx, auditRecord); err != nil {
			return fmt.Errorf("storage.AddAuditRecord: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("storage.WithTransaction: %w", txErr)
	}

	return nil
}
