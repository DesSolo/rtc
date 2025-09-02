package provider

import (
	"context"
	"errors"
	"fmt"
	"path"

	"rtc/internal/models"
	"rtc/internal/storage"
)

// Configs ...
func (p *Provider) Configs(ctx context.Context, projectName, envName, releaseName string) ([]*models.Config, error) {
	// TODO: use sharded database

	configs, err := p.storage.Configs(ctx, projectName, envName, releaseName)
	if err != nil {
		return nil, fmt.Errorf("storage.Configs: %w", err)
	}

	if len(configs) == 0 {
		return nil, nil
	}

	values, err := p.valuesStorage.Values(ctx, formatValuesStoragePath(projectName, envName, releaseName))
	if err != nil {
		return nil, fmt.Errorf("valuesStorage.Values: %w", err)
	}

	return convertConfigsToModel(configs, values), nil
}

// SetConfigValue ...
func (p *Provider) SetConfigValue(ctx context.Context, projectName, envName, releaseName, key string, value []byte) error {
	configStorage, err := p.storage.Config(ctx, projectName, envName, releaseName, key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("storage.Config: %w", err)
	}

	valuesStorageKey := formatValuesStorageKey(projectName, envName, releaseName, key)

	actualValue, err := p.valuesStorage.Value(ctx, valuesStorageKey)
	if err != nil {
		return fmt.Errorf("valuesStorage.Value: %w", err)
	}

	config, err := decodeConfigToModel(configStorage, actualValue)
	if err != nil {
		return fmt.Errorf("decodeConfigToModel: %w", err)
	}

	if err := validateNewValue(config, value); err != nil {
		return fmt.Errorf("validateNewValue: %w", err)
	}

	if err := p.valuesStorage.SetValue(ctx, valuesStorageKey, value); err != nil {
		return fmt.Errorf("valuesStorage.SetValue: %w", err)
	}

	// TODO: update updated_at field and add to log

	return nil
}

// UpsertConfigs ...
func (p *Provider) UpsertConfigs(ctx context.Context, projectName, envName, releaseName string, configs []*models.Config) error {
	project, err := p.storage.ProjectByName(ctx, projectName)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}

		return fmt.Errorf("storage.ProjectByName: %w", err)
	}

	for _, config := range configs {
		if err := validateUpsert(config); err != nil {
			return fmt.Errorf("key %q validateUpsert: %w", config.Key, err)
		}
	}

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		environment, err := storage.GetOrCreateEnvironment(ctx, p.storage, project.ID, envName)
		if err != nil {
			return fmt.Errorf("storage.GetOrCreateEnvironment: %w", err)
		}

		release, err := storage.GetOrCreateRelease(ctx, p.storage, environment.ID, releaseName)
		if err != nil {
			return fmt.Errorf("storage.GetOrCreateRelease: %w", err)
		}

		newKeys, err := p.getNewValuesStorageKeys(ctx, configs, projectName, envName, releaseName)
		if err != nil {
			return fmt.Errorf("getNewKeys: %w", err)
		}

		if len(newKeys) != 0 {
			if err := p.valuesStorage.SetValues(ctx, newKeys); err != nil {
				return fmt.Errorf("valuesStorage.SetValues: %w", err)
			}
		}

		if err := p.storage.UpsertConfigs(ctx, convertModelsToConfigs(configs, release.ID)); err != nil {
			return fmt.Errorf("storage.UpsertConfigs: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("storage.WithTransaction: %w", txErr)
	}

	return nil
}

func (p *Provider) getNewValuesStorageKeys(ctx context.Context, configs []*models.Config, projectName, envName, releaseName string) (map[storage.ValuesStorageKey]storage.ValuesStorageValue, error) {
	newKeys := make(map[storage.ValuesStorageKey]storage.ValuesStorageValue, len(configs))

	for _, config := range configs {
		key := formatValuesStorageKey(projectName, envName, releaseName, config.Key)

		if _, err := p.valuesStorage.Value(ctx, key); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				newKeys[key] = config.Value
				continue
			}

			return nil, fmt.Errorf("valuesStorage.Value: %w", err)
		}
	}

	return newKeys, nil
}

func formatValuesStoragePath(projectName, envName, releaseName string) storage.ValuesStoragePath {
	return storage.ValuesStoragePath(
		path.Join(projectName, envName, releaseName),
	)
}

func formatValuesStorageKey(projectName, envName, releaseName, key string) storage.ValuesStorageKey {
	return storage.ValuesStorageKey(
		path.Join(projectName, envName, releaseName, key),
	)
}
