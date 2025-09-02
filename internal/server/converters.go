package server

import "rtc/internal/models"

func convertModelsToProjects(projects []*models.Project) []project {
	result := make([]project, 0, len(projects))
	for _, proj := range projects {
		result = append(result, convertModelToProject(proj))
	}

	return result
}

func convertModelToProject(modelProject *models.Project) project {
	return project{
		Name:        modelProject.Name,
		Description: modelProject.Description,
		CreatedAt:   modelProject.CreatedAt,
	}
}

func convertModelsToEnvironments(environments []*models.Environment) []environment {
	result := make([]environment, 0, len(environments))
	for _, modelEnv := range environments {
		result = append(result, environment{
			Name: modelEnv.Name,
		})
	}

	return result
}

func convertModelsToReleases(releases []*models.Release) []release {
	result := make([]release, 0, len(releases))
	for _, modelRelease := range releases {
		result = append(result, release{
			Name:      modelRelease.Name,
			CreatedAt: modelRelease.CreatedAt,
		})
	}

	return result
}

func convertModelsToConfigs(configs []*models.Config) []config {
	result := make([]config, 0, len(configs))
	for _, modelConfig := range configs {
		result = append(result, config{
			Key:       modelConfig.Key,
			ValueType: string(modelConfig.ValueType),
			Value:     string(modelConfig.Value),
			Group:     modelConfig.Metadata.Group,
			Usage:     modelConfig.Metadata.Usage,
			Writable:  modelConfig.Metadata.Writable,
			View: configView{
				Enum: modelConfig.Metadata.View.Enum,
			},
			CreatedAt: modelConfig.CreatedAt,
			UpdatedAt: modelConfig.UpdatedAt,
		})
	}

	return result
}

func convertConfigsToModels(configs []config) []*models.Config {
	result := make([]*models.Config, 0, len(configs))
	for _, conf := range configs {
		result = append(result, &models.Config{
			Key:       conf.Key,
			ValueType: models.ConvertValueTypeToModel(conf.ValueType),
			Value:     []byte(conf.Value),
			Metadata: models.ConfigMetadata{
				Group:    conf.Group,
				Usage:    conf.Usage,
				Writable: conf.Writable,
				View: models.ConfigMetadataView{
					Enum: conf.View.Enum,
				},
			},
		})
	}

	return result
}
