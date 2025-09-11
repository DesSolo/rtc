package provider

import (
	"context"
	"fmt"

	"rtc/internal/models"
)

// AuditsSearch ...
func (p *Provider) AuditsSearch(ctx context.Context, filter models.AuditFilter) ([]*models.Audit, error) {
	audits, err := p.storage.AuditsSearch(ctx, convertModelToAuditFilter(filter))
	if err != nil {
		return nil, fmt.Errorf("storage.AuditsSearch: %w", err)
	}

	return convertAuditsToModels(audits), nil
}
