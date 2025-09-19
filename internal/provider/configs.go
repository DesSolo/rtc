package provider

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/samber/lo"

	"github.com/DesSolo/rtc/internal/auth"
	"github.com/DesSolo/rtc/internal/models"
	"github.com/DesSolo/rtc/internal/storage"
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

	values, err := p.valuesStorage.ValuesByPath(ctx, formatValuesStoragePath(projectName, envName, releaseName))
	if err != nil {
		return nil, fmt.Errorf("valuesStorage.ValuesByPath: %w", err)
	}

	return convertConfigsToModel(configs, values), nil
}

// SetConfigValues ...
func (p *Provider) SetConfigValues(ctx context.Context, projectName, envName, releaseName string, kv models.KV) error {
	storageKeys := lo.MapToSlice(kv, func(key string, _ []byte) string {
		return key
	})

	configsFromStorage, err := p.storage.ConfigsByKeys(ctx, projectName, envName, releaseName, storageKeys)
	if err != nil {
		return fmt.Errorf("storage.ConfigsByKeys: %w", err)
	}

	valuesStorageKeys := lo.Map(configsFromStorage, func(item *storage.Config, _ int) storage.ValuesStorageKey {
		return formatValuesStorageKey(projectName, envName, releaseName, item.Key)
	})

	actualValues, err := p.valuesStorage.Values(ctx, valuesStorageKeys)
	if err != nil {
		return fmt.Errorf("valuesStorage.Values: %w", err)
	}

	newValuesStorageItems := make(storage.ValuesStorageKV, len(actualValues))
	auditLogItems := make([]*auditRecordConfigUpdatedItems, 0, len(actualValues))

	for _, config := range convertConfigsToModel(configsFromStorage, actualValues) {
		newValue, ok := kv[config.Key]
		if !ok {
			return fmt.Errorf("missing key %s", config.Key)
		}

		if err := validateNewValue(config, newValue); err != nil {
			return fmt.Errorf("validateNewValue key: %q err: %w", config.Key, err)
		}

		valuesStorageKey := formatValuesStorageKey(projectName, envName, releaseName, config.Key)

		newValuesStorageItems[valuesStorageKey] = newValue

		auditLogItems = append(auditLogItems, &auditRecordConfigUpdatedItems{
			Key:      config.Key,
			OldValue: string(config.Value),
			NewValue: string(newValue),
		})
	}

	updatedStorageConfigIDs := make([]uint64, 0, len(configsFromStorage))
	for _, config := range configsFromStorage {
		updatedStorageConfigIDs = append(updatedStorageConfigIDs, config.ID)
	}

	actor := auth.UsernameFromContext(ctx)
	auditRecord, err := encodeAuditRecordConfigUpdated(actor, &auditRecordConfigUpdated{
		ProjectName:     projectName,
		EnvironmentName: envName,
		ReleaseName:     releaseName,
		Items:           auditLogItems,
	})
	if err != nil {
		return fmt.Errorf("encodeAuditRecordConfigUpdated: %w", err)
	}

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		if err := p.valuesStorage.SetValues(ctx, newValuesStorageItems); err != nil {
			return fmt.Errorf("valuesStorage.SetValue: %w", err)
		}

		if err := p.storage.MarkConfigsUpdated(ctx, updatedStorageConfigIDs); err != nil {
			return fmt.Errorf("storage.MarkConfigUpdated: %w", err)
		}

		if err := p.storage.AddAuditRecord(ctx, auditRecord); err != nil {
			return fmt.Errorf("storage.AddAuditRecord: %w", err)
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("storage.WithTransaction: %w", txErr)
	}

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

	if err := validateUpsert(configs); err != nil {
		return fmt.Errorf("validateUpsert: %w", err)
	}

	environment, err := storage.GetOrCreateEnvironment(ctx, p.storage, project.ID, envName)
	if err != nil {
		return fmt.Errorf("storage.GetOrCreateEnvironment: %w", err)
	}

	release, err := storage.GetOrCreateRelease(ctx, p.storage, environment.ID, releaseName)
	if err != nil {
		return fmt.Errorf("storage.GetOrCreateRelease: %w", err)
	}

	newValuesStorageItems, err := p.resolveNewValuesStorageItems(ctx, configs, projectName, envName, releaseName)
	if err != nil {
		return fmt.Errorf("resolveNewValuesStorageItems: %w", err)
	}

	toDelete, err := p.resolveToDeleteConfigs(ctx, configs, projectName, envName, releaseName)
	if err != nil {
		return fmt.Errorf("resolveToDeleteConfigsIDs: %w", err)
	}

	txErr := p.storage.WithTransaction(ctx, func(ctx context.Context) error {
		if len(newValuesStorageItems) != 0 {
			if err := p.valuesStorage.SetValues(ctx, newValuesStorageItems); err != nil {
				return fmt.Errorf("valuesStorage.SetValues: %w", err)
			}
		}

		if err := p.storage.UpsertConfigs(ctx, convertModelsToConfigs(configs, release.ID)); err != nil {
			return fmt.Errorf("storage.UpsertConfigs: %w", err)
		}

		if len(toDelete.valuesStorageKeys) != 0 {
			if err := p.valuesStorage.DeleteValues(ctx, toDelete.valuesStorageKeys); err != nil {
				return fmt.Errorf("valuesStorage.DeleteValues: %w", err)
			}
		}

		if len(toDelete.ids) != 0 {
			if err := p.storage.DeleteConfigs(ctx, toDelete.ids); err != nil {
				return fmt.Errorf("storage.DeleteConfigs: %w", err)
			}
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("storage.WithTransaction: %w", txErr)
	}

	return nil
}

func (p *Provider) resolveNewValuesStorageItems(ctx context.Context, configs []*models.Config, projectName, envName, releaseName string) (map[storage.ValuesStorageKey]storage.ValuesStorageValue, error) {
	var (
		keys          = make([]storage.ValuesStorageKey, 0, len(configs))
		keyToNewValue = make(map[storage.ValuesStorageKey][]byte, len(configs))
	)

	for _, config := range configs {
		key := formatValuesStorageKey(projectName, envName, releaseName, config.Key)
		keys = append(keys, key)
		keyToNewValue[key] = config.Value
	}

	newKeys := make(map[storage.ValuesStorageKey]storage.ValuesStorageValue, len(configs))

	actualKeys, err := p.valuesStorage.Values(ctx, keys)
	if err != nil {
		return nil, fmt.Errorf("valuesStorage.Values: %w", err)
	}

	for _, key := range keys {
		if _, ok := actualKeys[key]; ok {
			continue
		}

		newKeys[key] = keyToNewValue[key]
	}

	return newKeys, nil
}

type toDeleteConfigs struct {
	ids               []uint64
	valuesStorageKeys []storage.ValuesStorageKey
}

func (p *Provider) resolveToDeleteConfigs(ctx context.Context, configs []*models.Config, projectName, envName, releaseName string) (*toDeleteConfigs, error) {
	newKeys := make(map[string]struct{}, len(configs))
	for _, config := range configs {
		newKeys[config.Key] = struct{}{}
	}

	allActualConfigs, err := p.storage.Configs(ctx, projectName, envName, releaseName)
	if err != nil {
		return nil, fmt.Errorf("storage.Configs: %w", err)
	}

	var toDelete toDeleteConfigs

	for _, actualConfig := range allActualConfigs {
		if _, ok := newKeys[actualConfig.Key]; !ok {
			toDelete.ids = append(toDelete.ids, actualConfig.ID)
			toDelete.valuesStorageKeys = append(toDelete.valuesStorageKeys, formatValuesStorageKey(projectName, envName, releaseName, actualConfig.Key))
		}
	}

	return &toDelete, nil
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
