package env

import (
	"context"
	"os"
	"strings"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/internal"
)

// Provider ...
type Provider struct {
	envPrefix string
}

// NewProvider ...
func NewProvider(serviceName string) *Provider {
	return &Provider{
		envPrefix: strings.ToUpper(serviceName) + "_",
	}
}

// Value ...
func (c *Provider) Value(_ context.Context, key rtc.Key) (rtc.Value, error) {
	val := os.Getenv(c.encodeKeyPath(key))
	return internal.NewValue([]byte(val)), nil
}

// WatchValue ...
func (c *Provider) WatchValue(_ context.Context, _ rtc.Key, _ rtc.ValueChangeCallback) error {
	return rtc.ErrNotImplemented
}

// Close ...
func (c *Provider) Close() error {
	return nil
}

func (c *Provider) encodeKeyPath(key rtc.Key) string {
	return c.envPrefix + strings.ToUpper(string(key))
}
