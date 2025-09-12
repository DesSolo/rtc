package auth

import (
	"context"
)

// Noop ...
type Noop struct{}

// NewNoop ...
func NewNoop() *Noop {
	return &Noop{}
}

// Authorize ...
func (n *Noop) Authorize(_ context.Context, _ map[string]any) error {
	return nil
}
