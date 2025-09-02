package provider

import "rtc/internal/storage"

// Provider ...
type Provider struct {
	storage       storage.Storage
	valuesStorage storage.ValuesStorage
}

// NewProvider ...
func NewProvider(storage storage.Storage, valuesStorage storage.ValuesStorage) *Provider {
	return &Provider{
		storage:       storage,
		valuesStorage: valuesStorage,
	}
}
