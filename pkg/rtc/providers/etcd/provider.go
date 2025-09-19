package etcd

import (
	"context"
	"fmt"
	"path"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/internal"
)

const (
	defaultClientDialTimeout = 1 * time.Second
)

// Provider ...
type Provider struct {
	client *clientv3.Client
	path   string
}

// NewProvider ...
func NewProvider(endpoints []string, project, env, release string, options ...OptionFunc) (*Provider, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: defaultClientDialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("clientv3.New: %w", err)
	}

	c := &Provider{
		client: client,
		path:   path.Join("rtc", project, env, release),
	}

	for _, option := range options {
		option(c)
	}

	return c, nil
}

// Value ...
func (p *Provider) Value(ctx context.Context, key rtc.Key) (rtc.Value, error) {
	resp, err := p.client.Get(ctx, p.encodeKeyPath(key))
	if err != nil {
		return nil, fmt.Errorf("client.Get: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("key not found")
	}

	return internal.NewValue(resp.Kvs[0].Value), nil
}

// WatchValue ...
func (p *Provider) WatchValue(ctx context.Context, key rtc.Key, handler rtc.ValueChangeCallback) error {
	ch := p.client.Watch(ctx, p.encodeKeyPath(key), clientv3.WithPrevKV())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-ch:
				for _, event := range resp.Events {
					handler(convertPrevKeyToValue(event), internal.NewValue(event.Kv.Value))
				}
			}
		}
	}()

	return nil
}

// Close ...
func (p *Provider) Close() error {
	return p.client.Close() // nolint:wrapcheck
}

func (p *Provider) encodeKeyPath(key rtc.Key) string {
	return path.Join(p.path, string(key))
}

func convertPrevKeyToValue(event *clientv3.Event) rtc.Value {
	if event.PrevKv == nil {
		return internal.NewValue([]byte{})
	}

	return internal.NewValue(event.PrevKv.Value)
}
