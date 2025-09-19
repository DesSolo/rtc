package provider

import (
	"context"
	"fmt"

	"github.com/DesSolo/rtc/internal/models"
)

// ListEnvironments ...
func (p *Provider) ListEnvironments(ctx context.Context, projectName string) ([]*models.Environment, error) {
	environments, err := p.storage.Environments(ctx, projectName)
	if err != nil {
		return nil, fmt.Errorf("storage.ListEnvironments: %w", err)
	}

	return convertEnvironmentsToModel(environments), nil
}
