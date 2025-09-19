package provider

import (
	"testing"

	"github.com/DesSolo/rtc/internal/storage/mocks"
)

type mk struct {
	storage       *mocks.MockStorage
	valuesStorage *mocks.MockValuesStorage
	provider      *Provider
}

func newMk(t *testing.T) *mk {
	t.Helper()

	storage := mocks.NewMockStorage(t)
	valuesStorage := mocks.NewMockValuesStorage(t)

	return &mk{
		storage:       storage,
		valuesStorage: valuesStorage,
		provider:      NewProvider(storage, valuesStorage),
	}
}
