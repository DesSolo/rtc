package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/DesSolo/rtc/internal/auth"
	"github.com/DesSolo/rtc/internal/models"
	"github.com/DesSolo/rtc/internal/storage"
)

// Projects ...
func (p *Provider) Projects(ctx context.Context, q string, limit, offset uint64) ([]*models.Project, uint64, error) {
	projects, total, err := p.storage.Projects(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("storage.Projects: %w", err)
	}

	return convertProjectsToModel(projects), total, nil
}

// CreateProject ...
func (p *Provider) CreateProject(ctx context.Context, name, description string) (*models.Project, error) {
	project := &storage.Project{
		Name:        name,
		Description: description,
	}

	actor := auth.UsernameFromContext(ctx)
	auditRecord, err := encodeAuditRecordProjectCreated(actor, name, description)
	if err != nil {
		return nil, fmt.Errorf("encodeAuditRecordProjectCreated: %w", err)
	}

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		if err := p.storage.CreateProject(ctx, project); err != nil {
			if errors.Is(err, storage.ErrAlreadyExists) {
				return ErrAlreadyExists
			}

			return fmt.Errorf("storage.CreateProject: %w", err)
		}

		if err := p.storage.AddAuditRecord(ctx, auditRecord); err != nil {
			return fmt.Errorf("storage.AddAuditRecord: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr // nolint:wrapcheck
	}

	return convertProjectToModel(project), nil
}

// UpdateProjectDescription ...
func (p *Provider) UpdateProjectDescription(ctx context.Context, name, newDescription string) error {
	project, err := p.storage.ProjectByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("storage.ProjectByName: %w", err)
	}

	actor := auth.UsernameFromContext(ctx)
	auditRecord, err := encodeAuditRecordProjectUpdated(actor, name, project.Description, newDescription)
	if err != nil {
		return fmt.Errorf("encodeAuditRecordProjectUpdated: %w", err)
	}

	project.Description = newDescription

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		if err := p.storage.UpdateProject(ctx, project); err != nil {
			return fmt.Errorf("storage.UpdateProject: %w", err)
		}
		if err := p.storage.AddAuditRecord(ctx, auditRecord); err != nil {
			return fmt.Errorf("storage.AddAuditRecord: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return txErr // nolint:wrapcheck
	}

	return nil
}

// DeleteProject ...
func (p *Provider) DeleteProject(ctx context.Context, name string) error {
	project, err := p.storage.ProjectByName(ctx, name)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("storage.ProjectByName: %w", err)
	}

	actor := auth.UsernameFromContext(ctx)
	auditRecord, err := encodeAuditRecordProjectDeled(actor, name)
	if err != nil {
		return fmt.Errorf("encodeAuditRecordProjectDeled: %w", err)
	}

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		if err := p.storage.DeleteProject(ctx, project.ID); err != nil {
			return fmt.Errorf("storage.DeleteProject: %w", err)
		}

		if err := p.storage.AddAuditRecord(ctx, auditRecord); err != nil {
			return fmt.Errorf("storage.AddAuditRecord: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return txErr // nolint:wrapcheck
	}

	return nil
}
