package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"rtc/internal/models"
	"rtc/internal/storage"
)

func Test_AuditSearch_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	now := time.Now()
	fromDate := now.AddDate(0, 0, -1)
	toDate := now

	m.storage.EXPECT().AuditsSearch(
		mock.AnythingOfType("context.backgroundCtx"),
		storage.AuditFilter{
			Action:   "project_updated",
			Actor:    "test_actor",
			FromDate: fromDate,
			ToDate:   toDate,
		},
	).Return([]*storage.Audit{
		{
			ID:     10,
			Action: "project_updated",
			Actor:  "test_actor",
			Ts:     now,
		},
	}, nil)

	got, err := m.provider.AuditsSearch(context.Background(), models.AuditFilter{
		Action:   "project_updated",
		Actor:    "test_actor",
		FromDate: fromDate,
		ToDate:   toDate,
	})

	require.NoError(t, err)
	require.Equal(t, got, []*models.Audit{
		{
			Action: models.AuditActionProjectUpdated,
			Actor:  "test_actor",
			Ts:     now,
		},
	})
}

func Test_AuditSearch_StorageError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	now := time.Now()
	fromDate := now.AddDate(0, 0, -1)
	toDate := now

	m.storage.EXPECT().AuditsSearch(
		mock.AnythingOfType("context.backgroundCtx"),
		storage.AuditFilter{
			Action:   "project_updated",
			Actor:    "test_actor",
			FromDate: fromDate,
			ToDate:   toDate,
		},
	).Return(nil, errors.New("some error"))

	got, err := m.provider.AuditsSearch(context.Background(), models.AuditFilter{
		Action:   "project_updated",
		Actor:    "test_actor",
		FromDate: fromDate,
		ToDate:   toDate,
	})

	require.Empty(t, got)
	require.EqualError(t, err, "storage.AuditsSearch: some error")
}

func Test_AuditActions_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	got, err := m.provider.AuditActions(context.Background())

	require.NoError(t, err)
	require.Equal(t, got, []models.AuditAction{
		models.AuditActionConfigUpdated,
		models.AuditActionProjectCreated,
		models.AuditActionProjectUpdated,
		models.AuditActionProjectDeleted,
		models.AuditActionReleaseDeleted,
	})
}
