package provider

import (
	"context"
	"fmt"

	"rtc/internal/models"
)

// Audits ...
func (p *Provider) Audits(ctx context.Context, action models.AuditAction, limit, offset int) ([]*models.Audit, error) {
	audits, err := p.storage.AuditsByAction(ctx, string(action), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("storage.AuditsByAction: %w", err)
	}

	return convertAuditsToModels(audits), nil
}
