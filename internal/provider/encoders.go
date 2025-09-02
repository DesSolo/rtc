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

func encodeAuditRecordConfigUpdated(actor, key, oldValue, newValue string) (*storage.Audit, error) {
	type payloadV1 struct {
		Version  string `json:"version"`
		Key      string `json:"key"`
		OldValue string `json:"old_value"`
		NewValue string `json:"new_value"`
	}

	data, err := json.Marshal(payloadV1{
		Version:  "v1",
		Key:      key,
		OldValue: oldValue,
		NewValue: newValue,
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
