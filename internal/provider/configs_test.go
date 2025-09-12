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
