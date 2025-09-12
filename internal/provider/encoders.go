package provider

import (
	"encoding/json"
	"fmt"

	"rtc/internal/models"
	"rtc/internal/storage"
)

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

type auditRecordConfigUpdated struct {
	ProjectName     string
	EnvironmentName string
	ReleaseName     string
	Items           []*auditRecordConfigUpdatedItems
}

type auditRecordConfigUpdatedItems struct {
	Key      string
	OldValue string
	NewValue string
}

func encodeAuditRecordConfigUpdated(actor string, record *auditRecordConfigUpdated) (*storage.Audit, error) {
	type payloadV1Items struct {
		Key      string `json:"key"`
		OldValue string `json:"old_value"`
		NewValue string `json:"new_value"`
	}

	type payloadV1 struct {
		Version     string           `json:"version"`
		Project     string           `json:"project"`
		Environment string           `json:"environment"`
		Release     string           `json:"release"`
		Items       []payloadV1Items `json:"items"`
	}

	items := make([]payloadV1Items, 0, len(record.Items))
	for _, item := range record.Items {
		items = append(items, payloadV1Items{
			Key:      item.Key,
			OldValue: item.OldValue,
			NewValue: item.NewValue,
		})
	}

	data, err := json.Marshal(payloadV1{
		Version:     "v1",
		Project:     record.ProjectName,
		Environment: record.EnvironmentName,
		Release:     record.ReleaseName,
		Items:       items,
	})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return &storage.Audit{
		Action:  string(models.AuditActionConfigUpdated),
		Actor:   actor,
		Payload: data,
	}, nil
}

func encodeAuditRecordProjectCreated(actor string, projectName, description string) (*storage.Audit, error) {
	type payloadV1 struct {
		Version     string `json:"version"`
		Project     string `json:"project"`
		Description string `json:"description"`
	}

	data, err := json.Marshal(payloadV1{
		Version:     "v1",
		Project:     projectName,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	return &storage.Audit{
		Action:  string(models.AuditActionProjectCreated),
		Actor:   actor,
		Payload: data,
	}, nil
}

func encodeAuditRecordProjectUpdated(actor string, projectName, oldDescription, newDescription string) (*storage.Audit, error) {
	type payloadV1 struct {
		Version        string `json:"version"`
		Project        string `json:"project"`
		OldDescription string `json:"old_description"`
		NewDescription string `json:"new_description"`
	}

	data, err := json.Marshal(payloadV1{
		Version:        "v1",
		Project:        projectName,
		OldDescription: oldDescription,
		NewDescription: newDescription,
	})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	return &storage.Audit{
		Action:  string(models.AuditActionProjectUpdated),
		Actor:   actor,
		Payload: data,
	}, nil
}

func encodeAuditRecordProjectDeled(actor string, projectName string) (*storage.Audit, error) {
	type payloadV1 struct {
		Version string `json:"version"`
		Project string `json:"project"`
	}

	data, err := json.Marshal(payloadV1{
		Version: "v1",
		Project: projectName,
	})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	return &storage.Audit{
		Action:  string(models.AuditActionProjectDeleted),
		Actor:   actor,
		Payload: data,
	}, nil
}

func encodeAuditRecordReleaseDeleted(actor string, projectName, envName, releaseName string) (*storage.Audit, error) {
	type payloadV1 struct {
		Version string `json:"version"`
		Project string `json:"project"`
		Env     string `json:"env"`
		Release string `json:"release"`
	}

	data, err := json.Marshal(payloadV1{
		Version: "v1",
		Project: projectName,
		Env:     envName,
		Release: releaseName,
	})
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return &storage.Audit{
		Action:  string(models.AuditActionReleaseDeleted),
		Actor:   actor,
		Payload: data,
	}, nil
}
