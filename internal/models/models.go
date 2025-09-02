package models

import "time"

type ValueType string

const (
	ValueTypeUnknown ValueType = ""
	ValueTypeString  ValueType = "string"
	ValueTypeBool    ValueType = "bool"
	ValueTypeInt     ValueType = "int"
	ValueTypeInt64   ValueType = "int64"
	ValueTypeUint    ValueType = "uint"
	ValueTypeUint64  ValueType = "uint64"
	ValueTypeFloat   ValueType = "float"
	ValueTypeFloat64 ValueType = "float64"
)

type ConfigMetadataView struct {
	Enum string
}

type ConfigMetadata struct {
	Group    string
	Usage    string
	Writable bool
	View     ConfigMetadataView
}

type Config struct {
	Key       string
	ValueType ValueType
	Value     []byte
	Metadata  ConfigMetadata

	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Project struct {
	Name        string
	Description string
	CreatedAt   time.Time
}

type Environment struct {
	Name string
}

type Release struct {
	Name      string
	CreatedAt time.Time
}
