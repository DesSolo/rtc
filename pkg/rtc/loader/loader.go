package loader

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/DesSolo/rtc/pkg/rtc"
	"github.com/DesSolo/rtc/pkg/rtc/internal"
	"github.com/DesSolo/rtc/pkg/rtc/providers/env"
)

var (
	defaultClient rtc.Client
	once          sync.Once
)

func init() {
	once.Do(func() {
		defaultClient = env.NewProvider(
			os.Getenv("SERVICE_NAME"),
		)
	})
}

// Default return a default client
func Default() rtc.Client {
	return defaultClient
}

// SetDefault set client as default
func SetDefault(c rtc.Client) {
	defaultClient = c
}

// Value return a value with error
func Value(ctx context.Context, key rtc.Key) (rtc.Value, error) {
	return defaultClient.Value(ctx, key) // nolint:wrapcheck
}

// Get return a value
func Get(ctx context.Context, key rtc.Key) rtc.Value {
	val, err := Value(ctx, key)
	if err != nil {
		slog.ErrorContext(ctx, "rtc: failed to get value", "err", err)
		return internal.NewValue(nil)
	}

	return val
}

// WatchValue watch value changes
func WatchValue(ctx context.Context, key rtc.Key, handler rtc.ValueChangeCallback) error {
	return defaultClient.WatchValue(ctx, key, handler) // nolint:wrapcheck
}
