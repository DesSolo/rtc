package provider

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DesSolo/rtc/internal/models"
	"github.com/DesSolo/rtc/internal/storage"
)

func Test_Configs_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	createdAt := time.Now().Add(-time.Hour)

	m.storage.EXPECT().
		Configs(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
		).
		Return([]*storage.Config{
			{
				ID:        10,
				ReleaseID: 11,
				Key:       "test_key",
				ValueType: "string",
				Metadata:  []byte(`{"version":"v1"}`),
				CreatedAt: createdAt,
			},
		}, nil)

	m.valuesStorage.EXPECT().
		ValuesByPath(
			mock.AnythingOfType("context.backgroundCtx"),
			storage.ValuesStoragePath("test_project/test_env/test_release"),
		).Return(storage.ValuesStorageKV{
		"test_key": []byte("test_value"),
	}, nil)

	got, err := m.provider.Configs(context.Background(), "test_project", "test_env", "test_release")
	require.NoError(t, err)
	require.Equal(t, []*models.Config{
		{
			Key:       "test_key",
			ValueType: "string",
			Value:     []byte("test_value"),
			CreatedAt: createdAt,
		},
	}, got)
}

func Test_Configs_NoConfigs_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		Configs(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
		).
		Return(nil, nil)

	got, err := m.provider.Configs(context.Background(), "test_project", "test_env", "test_release")
	require.NoError(t, err)
	require.Empty(t, got)
}

func Test_Configs_StorageError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		Configs(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
		).
		Return(nil, errors.New("test error"))

	got, err := m.provider.Configs(context.Background(), "test_project", "test_env", "test_release")
	require.Empty(t, got)
	require.EqualError(t, err, "storage.Configs: test error")
}

func Test_Configs_ValuesStorageError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		Configs(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
		).
		Return([]*storage.Config{
			{
				ID:        10,
				ReleaseID: 11,
				Key:       "test_key",
				ValueType: "string",
			},
		}, nil)

	m.valuesStorage.EXPECT().
		ValuesByPath(
			mock.AnythingOfType("context.backgroundCtx"),
			storage.ValuesStoragePath("test_project/test_env/test_release"),
		).Return(nil, errors.New("test error"))

	got, err := m.provider.Configs(context.Background(), "test_project", "test_env", "test_release")
	require.Empty(t, got)
	require.EqualError(t, err, "valuesStorage.ValuesByPath: test error")
}

func Test_SetConfigValues_ExpectOk(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		ConfigsByKeys(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
			[]string{"test_key"},
		).
		Return([]*storage.Config{
			{
				ID:        10,
				ReleaseID: 11,
				Key:       "test_key",
				ValueType: "string",
				Metadata:  []byte(`{"version":"v1", "writable": true}`),
			},
		}, nil)

	m.valuesStorage.EXPECT().
		Values(
			mock.AnythingOfType("context.backgroundCtx"),
			[]storage.ValuesStorageKey{"test_project/test_env/test_release/test_key"},
		).
		Return(storage.ValuesStorageKV{
			"test_key": []byte("old_value"),
		}, nil)

	m.valuesStorage.EXPECT().
		SetValues(
			mock.AnythingOfType("context.backgroundCtx"),
			storage.ValuesStorageKV{
				"test_project/test_env/test_release/test_key": []byte("new_value"),
			},
		).
		Return(nil)

	m.storage.EXPECT().
		MarkConfigsUpdated(
			mock.AnythingOfType("context.backgroundCtx"),
			[]uint64{10},
		).
		Return(nil)

	m.storage.EXPECT().
		AddAuditRecord(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.Anything,
		).
		Return(nil)

	m.storage.EXPECT().
		WithTransaction(
			mock.AnythingOfType("context.backgroundCtx"),
			mock.Anything,
		).
		Run(func(ctx context.Context, f func(ctx context.Context) error) {
			require.NoError(t, f(ctx))
		}).
		Return(nil)

	err := m.provider.SetConfigValues(context.Background(), "test_project", "test_env", "test_release", models.KV{
		"test_key": []byte("new_value"),
	})
	require.NoError(t, err)
}

func Test_SetConfigValues_ConfigsByKeysError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		ConfigsByKeys(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
			[]string{"test_key"},
		).
		Return(nil, errors.New("storage error"))

	err := m.provider.SetConfigValues(context.Background(), "test_project", "test_env", "test_release", models.KV{
		"test_key": []byte("new_value"),
	})
	require.EqualError(t, err, "storage.ConfigsByKeys: storage error")
}

func Test_SetConfigValues_ValuesError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		ConfigsByKeys(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
			"test_env",
			"test_release",
			[]string{"test_key"},
		).
		Return([]*storage.Config{
			{
				ID:        10,
				ReleaseID: 11,
				Key:       "test_key",
				ValueType: "string",
			},
		}, nil)

	m.valuesStorage.EXPECT().
		Values(
			mock.AnythingOfType("context.backgroundCtx"),
			[]storage.ValuesStorageKey{"test_project/test_env/test_release/test_key"},
		).
		Return(nil, errors.New("values error"))

	err := m.provider.SetConfigValues(context.Background(), "test_project", "test_env", "test_release", models.KV{
		"test_key": []byte("new_value"),
	})
	require.EqualError(t, err, "valuesStorage.Values: values error")
}

func Test_UpsertConfigs_ProjectNotFound_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		ProjectByName(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
		).
		Return(nil, storage.ErrNotFound)

	configs := []*models.Config{
		{
			Key:       "test_key",
			ValueType: "string",
			Value:     []byte("value"),
		},
	}
	err := m.provider.UpsertConfigs(context.Background(), "test_project", "test_env", "test_release", configs)
	require.ErrorIs(t, err, ErrNotFound)
}

func Test_UpsertConfigs_ProjectByNameError_ExpectErr(t *testing.T) {
	t.Parallel()

	m := newMk(t)

	m.storage.EXPECT().
		ProjectByName(
			mock.AnythingOfType("context.backgroundCtx"),
			"test_project",
		).
		Return(nil, errors.New("project error"))

	configs := []*models.Config{
		{
			Key:       "test_key",
			ValueType: "string",
			Value:     []byte("value"),
		},
	}
	err := m.provider.UpsertConfigs(context.Background(), "test_project", "test_env", "test_release", configs)
	require.EqualError(t, err, "storage.ProjectByName: project error")
}
