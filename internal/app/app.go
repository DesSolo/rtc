package app

import (
	"context"
	"fmt"
)

// App ...
type App struct{}

// New ...
func New() *App {
	return &App{}
}

// Run ...
func (app *App) Run(ctx context.Context) error {
	di := newContainer()

	configureLogger(di)

	// nolint:contextcheck
	if err := di.Server().Run(ctx); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
