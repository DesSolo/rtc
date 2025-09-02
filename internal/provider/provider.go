package provider

import "rtc/internal/storage"

type Provider struct {
	storage       storage.Storage
	valuesStorage storage.ValuesStorage
}

func NewProvider(storage storage.Storage, valuesStorage storage.ValuesStorage) *Provider {
	return &Provider{
		storage:       storage,
		valuesStorage: valuesStorage,
	}
}
