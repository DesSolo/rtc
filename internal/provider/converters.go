package provider

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"rtc/internal/models"
	"rtc/internal/storage"
)

func convertProjectsToModel(projects []*storage.Project) []*models.Project {
	result := make([]*models.Project, 0, len(projects))
	for _, project := range projects {
		result = append(result, convertProjectToModel(project))
	}

	return result
}

func convertProjectToModel(project *storage.Project) *models.Project {
	return &models.Project{
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
	}
}

func convertEnvironmentsToModel(envs []*storage.Environment) []*models.Environment {
	result := make([]*models.Environment, 0, len(envs))
	for _, env := range envs {
		result = append(result, &models.Environment{
			Name: env.Name,
		})
	}

	return result
}

func convertReleasesToModel(releases []*storage.Release) []*models.Release {
	result := make([]*models.Release, 0, len(releases))
	for _, release := range releases {
		result = append(result, &models.Release{
			Name:      release.Name,
			CreatedAt: release.CreatedAt,
		})
	}

	return result
}

func convertConfigsToModel(configs []*storage.Config, values map[storage.ValuesStorageKey]storage.ValuesStorageValue) []*models.Config {
	result := make([]*models.Config, 0, len(configs))
	for _, config := range configs {
		value, ok := values[storage.ValuesStorageKey(config.Key)]
		if !ok {
			slog.Warn("not found value", "key", config.Key)
			continue
		}

		configModel, err := decodeConfigToModel(config, value)
		if err != nil {
			slog.Warn("failed to decode config", "key", config.Key, "value", value, "err", err)
			continue
		}

		result = append(result, configModel)
	}

	return result
}

func decodeConfigToModel(config *storage.Config, value storage.ValuesStorageValue) (*models.Config, error) {
	metadata, err := decodeMetadata(config.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	return &models.Config{
		Key:       config.Key,
		ValueType: models.ConvertValueTypeToModel(config.ValueType),
		Value:     value,
		Metadata:  metadata,
		CreatedAt: config.CreatedAt,
		UpdatedAt: config.UpdatedAt,
	}, nil
}

func convertModelsToConfigs(configs []*models.Config, releaseID uint64) []*storage.Config {
	storageConfigs := make([]*storage.Config, 0, len(configs))
	for _, config := range configs {
		meta, err := encodeMetadata(config.Metadata)
		if err != nil {
			slog.Warn("failed to encode metadata", "err", err)
			continue
		}

		storageConfigs = append(storageConfigs, &storage.Config{
			ReleaseID: releaseID,
			Key:       config.Key,
			ValueType: string(config.ValueType),
			Metadata:  meta,
		})
	}

	return storageConfigs
}

type metadataV1 struct {
	Version  string `json:"version"`
	Group    string `json:"group"`
	Usage    string `json:"usage"`
	Writable bool   `json:"writable"`
	View     struct {
		Enum string `json:"enum"`
	} `json:"view"`
}

func decodeMetadata(metadata []byte) (models.ConfigMetadata, error) {
	var meta metadataV1
	if err := json.Unmarshal(metadata, &meta); err != nil {
		return models.ConfigMetadata{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return models.ConfigMetadata{
		Group:    meta.Group,
		Usage:    meta.Usage,
		Writable: meta.Writable,
		View: models.ConfigMetadataView{
			Enum: meta.View.Enum,
		},
	}, nil
}

func encodeMetadata(metadata models.ConfigMetadata) ([]byte, error) {
	meta := metadataV1{
		Version:  "v1",
		Group:    metadata.Group,
		Usage:    metadata.Usage,
		Writable: metadata.Writable,
	}

	meta.View.Enum = metadata.View.Enum

	data, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return data, nil
}
