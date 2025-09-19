package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/DesSolo/rtc/internal/app"
	"github.com/DesSolo/rtc/pkg/closer"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()

		slog.Info("shutdown signal received")

		if err := closer.Close(); err != nil {
			log.Fatalf("failed to close: %s", err.Error())
		}
	}()

	application := app.New()

	if err := application.Run(ctx); err != nil {
		log.Fatalf("failed to start application: %s", err.Error())
	}
}
