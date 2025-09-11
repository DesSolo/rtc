package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"rtc/internal/storage"
)

// AuditsSearch ...
func (s *Storage) AuditsSearch(ctx context.Context, filter storage.AuditFilter) ([]*storage.Audit, error) {
	query := queryBuilder().Select("id, action, actor, payload, ts").
		From("audit_log").
		Where(squirrel.GtOrEq{"ts": filter.FromDate}).
		Where(squirrel.LtOrEq{"ts": filter.ToDate})

	if filter.Action != "" {
		query = query.Where(squirrel.Eq{"action": filter.Action})
	}

	if filter.Actor != "" {
		query = query.Where(squirrel.Eq{"actor": filter.Actor})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql: %w", err)
	}

	var audits []*storage.Audit

	rows, err := s.manager.Conn(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var audit storage.Audit
		if err := rows.Scan(&audit.ID, &audit.Action, &audit.Actor, &audit.Payload, &audit.Ts); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		audits = append(audits, &audit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows: %w", err)
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
