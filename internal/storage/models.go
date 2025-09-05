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

// Audit ...
type Audit struct {
	ID      uint64
	Action  string
	Actor   string
	Payload []byte
	Ts      time.Time
}

// ValuesStoragePath values path like foo/bar/baz
type ValuesStoragePath string

// ValuesStorageKey is a key in values storage
type ValuesStorageKey string

// ValuesStorageValue is a value in values storage
type ValuesStorageValue []byte

// ValuesStorageKV is alias for map
type ValuesStorageKV map[ValuesStorageKey]ValuesStorageValue
