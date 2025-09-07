package provider

import (
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

func convertAuditsToModels(audits []*storage.Audit) []*models.Audit {
	result := make([]*models.Audit, 0, len(audits))
	for _, audit := range audits {
		result = append(result, &models.Audit{
			Action:  models.ConvertAuditActionToModel(audit.Action),
			Actor:   audit.Actor,
			Payload: audit.Payload,
			Ts:      audit.Ts,
		})
	}

	return result
}

func convertUsersToModels(users []*storage.User) []*models.User {
	result := make([]*models.User, 0, len(users))
	for _, user := range users {
		result = append(result, convertUserToModel(user))
	}

	return result
}

func convertUserToModel(user *storage.User) *models.User {
	return &models.User{
		Username:  user.Username,
		IsEnabled: user.IsEnabled,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
	}
}

func convertModelToUser(user *models.User, passwordHash string) *storage.User {
	return &storage.User{
		Username:     user.Username,
		PasswordHash: passwordHash,
		IsEnabled:    user.IsEnabled,
		Roles:        user.Roles,
	}
}
