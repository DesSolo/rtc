package models

import "time"

// ValueType type of config value
type ValueType string

const (
	// ValueTypeUnknown ...
	ValueTypeUnknown ValueType = ""
	// ValueTypeString ...
	ValueTypeString ValueType = "string"
	// ValueTypeBool ...
	ValueTypeBool ValueType = "bool"
	// ValueTypeInt ...
	ValueTypeInt ValueType = "int"
	// ValueTypeInt64 ...
	ValueTypeInt64 ValueType = "int64"
	// ValueTypeUint ...
	ValueTypeUint ValueType = "uint"
	// ValueTypeUint64 ...
	ValueTypeUint64 ValueType = "uint64"
	// ValueTypeFloat ...
	ValueTypeFloat ValueType = "float"
	// ValueTypeFloat64 ...
	ValueTypeFloat64 ValueType = "float64"
)

// ConfigMetadataView specific view options
type ConfigMetadataView struct {
	Enum string
}

// ConfigMetadata advanced config metadata (can expand future)
type ConfigMetadata struct {
	Group    string
	Usage    string
	Writable bool
	View     ConfigMetadataView
}

// Config configuration
type Config struct {
	Key       string
	ValueType ValueType
	Value     []byte
	Metadata  ConfigMetadata

	CreatedAt time.Time
	UpdatedAt *time.Time
}

// Project some client for use config
type Project struct {
	Name        string
	Description string
	CreatedAt   time.Time
}

// Environment like dev stage prod
type Environment struct {
	Name string
}

// Release specific project release
type Release struct {
	Name      string
	CreatedAt time.Time
}

// AuditAction ...
type AuditAction string

const (
	// AuditActionUnknown ...
	AuditActionUnknown AuditAction = ""
	// AuditActionConfigUpdated ...
	AuditActionConfigUpdated AuditAction = "config_updated"
)

// Audit log record for history
type Audit struct {
	Action  AuditAction
	Actor   string
	Payload []byte
	Ts      time.Time
}
