package provider

import (
	"context"
	"errors"
	"fmt"

	"rtc/internal/models"
	"rtc/internal/storage"
)

// Projects ...
func (p *Provider) Projects(ctx context.Context, q string, limit, offset int) ([]*models.Project, error) {
	if len(q) != 0 {
		projects, err := p.storage.SearchProjects(ctx, q, limit)
		if err != nil {
			return nil, fmt.Errorf("storage.SearchProjects: %w", err)
		}

		return convertProjectsToModel(projects), nil
	}

	projects, err := p.storage.Projects(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("storage.Projects: %w", err)
	}

	return convertProjectsToModel(projects), nil
}

// CreateProject ...
func (p *Provider) CreateProject(ctx context.Context, name, description string) (*models.Project, error) {
	project := &storage.Project{
		Name:        name,
		Description: description,
	}

	if err := p.storage.CreateProject(ctx, project); err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, ErrAlreadyExists
		}

		return nil, fmt.Errorf("storage.CreateProject: %w", err)
	}

	return convertProjectToModel(project), nil
}
