package etcd

import (
	"context"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"

	"rtc/internal/storage"
)

const (
	defaultPath = "rtc"
)

// ValuesStorage ...
type ValuesStorage struct {
	client *clientv3.Client
	path   string
}

// NewValuesStorage ...
func NewValuesStorage(client *clientv3.Client) *ValuesStorage {
	return &ValuesStorage{client: client, path: defaultPath} // TODO: mode to options
}

// Values ...
func (s *ValuesStorage) Values(ctx context.Context, path storage.ValuesStoragePath) (map[storage.ValuesStorageKey]storage.ValuesStorageValue, error) {
	resp, err := s.client.Get(ctx, s.formatPath(string(path)), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", err)
	}

	values := make(map[storage.ValuesStorageKey]storage.ValuesStorageValue, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		values[decodeStorageKey(kv.Key)] = kv.Value
	}

	return values, nil
}

// Value ...
func (s *ValuesStorage) Value(ctx context.Context, key storage.ValuesStorageKey) (storage.ValuesStorageValue, error) {
	resp, err := s.client.Get(ctx, s.formatPath(string(key)))
	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", err)
	}

	if len(resp.Kvs) != 1 {
		return nil, storage.ErrNotFound
	}

	return resp.Kvs[0].Value, nil
}

// SetValue ...
func (s *ValuesStorage) SetValue(ctx context.Context, key storage.ValuesStorageKey, value storage.ValuesStorageValue) error {
	if _, err := s.client.Put(ctx, s.formatPath(string(key)), string(value)); err != nil {
		return fmt.Errorf("client.Put: %w", err)
	}

	return nil
}

// SetValues ...
func (s *ValuesStorage) SetValues(ctx context.Context, values map[storage.ValuesStorageKey]storage.ValuesStorageValue) error {
	ops := make([]clientv3.Op, 0, len(values))
	for key, value := range values {
		ops = append(ops, clientv3.OpPut(s.formatPath(string(key)), string(value)))
	}

	txn := s.client.Txn(ctx)
	if _, err := txn.Then(ops...).Commit(); err != nil {
		return fmt.Errorf("txn.Commit: %w", err)
	}

	return nil
}

// DeleteValues ...
func (s *ValuesStorage) DeleteValues(ctx context.Context, path storage.ValuesStoragePath) error {
	if _, err := s.client.Delete(ctx, s.formatPath(string(path)), clientv3.WithPrefix()); err != nil {
		return fmt.Errorf("client.Delete: %w", err)
	}

	return nil
}

func (s *ValuesStorage) formatPath(key string) string {
	return path.Join(s.path, key)
}

func decodeStorageKey(key []byte) storage.ValuesStorageKey {
	return storage.ValuesStorageKey(path.Base(string(key)))
}
