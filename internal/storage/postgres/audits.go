package postgres

import (
	"context"
	"fmt"

	"rtc/internal/storage"
)

// AuditsByAction ...
func (s *Storage) AuditsByAction(ctx context.Context, action string, limit, offset int) ([]*storage.Audit, error) {
	query := "SELECT id, action, actor, payload, ts FROM audit_log WHERE action = $1 ORDER BY id DESC LIMIT $2 OFFSET $3"

	rows, err := s.manager.Conn(ctx).Query(ctx, query, action, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pool.Query: %w", err)
	}
	defer rows.Close()

	var audits []*storage.Audit

	for rows.Next() {
		var audit storage.Audit
		if err := rows.Scan(&audit.ID, &audit.Action, &audit.Actor, &audit.Payload, &audit.Ts); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		audits = append(audits, &audit)
	}

	return audits, nil
}

// AddAuditRecord ...
func (s *Storage) AddAuditRecord(ctx context.Context, audit *storage.Audit) error {
	query := "INSERT INTO audit_log (action, actor, payload) VALUES ($1, $2, $3) RETURNING id"

	if err := s.manager.Conn(ctx).QueryRow(ctx, query, audit.Action, audit.Actor, audit.Payload).Scan(&audit.ID); err != nil {
		return fmt.Errorf("pool.Query: %w", err)
	}

	return nil
}
