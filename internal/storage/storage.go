package storage

import (
	"context"
)

// Storage ...
type Storage interface {
	Close() error

	WithTransaction(ctx context.Context, f func(ctx context.Context) error) error

	Projects(ctx context.Context, q string, limit, offset uint64) ([]*Project, uint64, error)
	ProjectByName(ctx context.Context, name string) (*Project, error)
	CreateProject(ctx context.Context, project *Project) error

	Environments(ctx context.Context, projectName string) ([]*Environment, error)
	Environment(ctx context.Context, projectID uint64, envName string) (*Environment, error)
	CreateEnvironment(ctx context.Context, env *Environment) error

	Releases(ctx context.Context, projectName, envName string) ([]*Release, error)
	Release(ctx context.Context, envID uint64, releaseName string) (*Release, error)
	CreateRelease(ctx context.Context, release *Release) error
	DeleteRelease(ctx context.Context, envID uint64, releaseName string) error

	Configs(ctx context.Context, projectName, envName, releaseName string) ([]*Config, error)
	ConfigsByKeys(ctx context.Context, projectName, envName, releaseName string, keys []string) ([]*Config, error)
	Config(ctx context.Context, projectName, envName, releaseName, key string) (*Config, error)
	UpsertConfigs(ctx context.Context, configs []*Config) error
	MarkConfigsUpdated(ctx context.Context, IDs []uint64) error
	DeleteConfigs(ctx context.Context, IDs []uint64) error

	AuditsSearch(ctx context.Context, filter AuditFilter) ([]*Audit, error)
	AddAuditRecord(ctx context.Context, audit *Audit) error

	Users(ctx context.Context, q string, limit, offset uint64) ([]*User, uint64, error)
	User(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, id uint64, user *User) error
}

// ValuesStorage ...
type ValuesStorage interface {
	Values(ctx context.Context, keys []ValuesStorageKey) (ValuesStorageKV, error)
	ValuesByPath(ctx context.Context, path ValuesStoragePath) (ValuesStorageKV, error)
	SetValues(ctx context.Context, values ValuesStorageKV) error
	DeleteValues(ctx context.Context, keys []ValuesStorageKey) error
	DeleteValuesByPath(ctx context.Context, path ValuesStoragePath) error
}
