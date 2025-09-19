package etcd

import (
	"context"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/DesSolo/rtc/internal/storage"
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
func NewValuesStorage(client *clientv3.Client, options ...OptionFunc) *ValuesStorage {
	s := &ValuesStorage{client: client, path: defaultPath}

	for _, option := range options {
		option(s)
	}

	return s
}

// Values ...
func (s *ValuesStorage) Values(ctx context.Context, keys []storage.ValuesStorageKey) (storage.ValuesStorageKV, error) {
	ops := make([]clientv3.Op, 0, len(keys))
	for _, key := range keys {
		ops = append(ops, clientv3.OpGet(s.formatPath(string(key))))
	}

	txnResp, err := s.client.Txn(ctx).Then(ops...).Commit()
	if err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	response := make(storage.ValuesStorageKV, len(keys))

	for _, kv := range txnResp.Responses {
		getResp := kv.GetResponseRange()
		if len(getResp.Kvs) != 1 {
			continue
		}

		response[decodeStorageKey(getResp.Kvs[0].Key)] = getResp.Kvs[0].Value
	}

	return response, nil
}

// ValuesByPath ...
func (s *ValuesStorage) ValuesByPath(ctx context.Context, path storage.ValuesStoragePath) (storage.ValuesStorageKV, error) {
	resp, err := s.client.Get(ctx, s.formatPath(string(path)), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", err)
	}

	values := make(storage.ValuesStorageKV, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		values[decodeStorageKey(kv.Key)] = kv.Value
	}

	return values, nil
}

// SetValues ...
func (s *ValuesStorage) SetValues(ctx context.Context, values storage.ValuesStorageKV) error {
	ops := make([]clientv3.Op, 0, len(values))
	for key, value := range values {
		ops = append(ops, clientv3.OpPut(s.formatPath(string(key)), string(value)))
	}

	if _, err := s.client.Txn(ctx).Then(ops...).Commit(); err != nil {
		return fmt.Errorf("txn.Commit: %w", err)
	}

	return nil
}

// DeleteValues ...
func (s *ValuesStorage) DeleteValues(ctx context.Context, keys []storage.ValuesStorageKey) error {
	ops := make([]clientv3.Op, 0, len(keys))
	for _, key := range keys {
		ops = append(ops, clientv3.OpDelete(s.formatPath(string(key))))
	}

	if _, err := s.client.Txn(ctx).Then(ops...).Commit(); err != nil {
		return fmt.Errorf("txn.Commit: %w", err)
	}
	return nil
}

// DeleteValuesByPath ...
func (s *ValuesStorage) DeleteValuesByPath(ctx context.Context, path storage.ValuesStoragePath) error {
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
