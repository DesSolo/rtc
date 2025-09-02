package rtc

import (
	"context"
)

// Key ...
type Key string

// Value ...
type Value interface {
	String() string
	MaybeString() (string, error)

	Float64() float64
	MaybeFloat64() (float64, error)

	Bool() bool
	MaybeBool() (bool, error)

	Int() int
	MaybeInt() (int, error)
}

// ValueChangeCallback ...
type ValueChangeCallback func(oldValue, newValue Value)

// Client ...
type Client interface {
	Value(context.Context, Key) (Value, error)
	WatchValue(ctx context.Context, key Key, handler ValueChangeCallback) error
	Close() error
}
