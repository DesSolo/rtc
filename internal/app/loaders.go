package app

import (
	"log/slog"
	"os"
)

func configureLogger(di *container) {
	options := di.Config().Logging

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(options.Level),
	})))
}
