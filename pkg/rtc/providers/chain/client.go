package chain

import (
	"context"
	"errors"

	"rtc/pkg/rtc"
)

// Provider ...
type Provider struct {
	clients []rtc.Client
}

// NewProvider ...
func NewProvider(clients ...rtc.Client) *Provider {
	return &Provider{clients: clients}
}

// Value ...
func (c *Provider) Value(ctx context.Context, key rtc.Key) (rtc.Value, error) {
	for _, client := range c.clients {
		val, err := client.Value(ctx, key)
		if err != nil {
			continue
		}

		return val, nil
	}

	return nil, rtc.ErrNotPresent
}

// WatchValue ...
func (c *Provider) WatchValue(ctx context.Context, key rtc.Key, handler rtc.ValueChangeCallback) error {
	for _, client := range c.clients {
		if err := client.WatchValue(ctx, key, handler); err != nil {
			if errors.Is(err, rtc.ErrNotPresent) {
				continue
			}

			return err
		}
	}

	return nil
}

// Close ...
func (c *Provider) Close() error {
	return nil
}
