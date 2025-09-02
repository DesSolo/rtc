package storage

import "time"

// Project ...
type Project struct {
	ID          uint64
	Name        string
	Description string
	CreatedAt   time.Time
}

// Environment ...
type Environment struct {
	ID        uint64
	ProjectID uint64
	Name      string
}

// Release ...
type Release struct {
	ID            uint64
	EnvironmentID uint64
	Name          string
	CreatedAt     time.Time
}

// Config ...
type Config struct {
	ID        uint64
	ReleaseID uint64
	Key       string
	ValueType string
	Metadata  []byte
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type ValuesStoragePath string
type ValuesStorageKey string

type ValuesStorageValue []byte
