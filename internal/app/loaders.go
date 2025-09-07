package app

import (
	"log/slog"
	"os"

	"rtc/internal/auth"
	"rtc/internal/server"
)

func configureLogger(di *container) {
	options := di.Config().Logging

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(options.Level),
	})))
}

func loadServerAuth(di *container) server.OptionFunc {
	options := di.Config().Server.Auth

	if len(options.Tokens) == 0 {
		return server.Noop()
	}

	tokenItems := make(map[string]*auth.Payload, len(options.Tokens))

	for token, payload := range options.Tokens {
		tokenItems[token] = &auth.Payload{
			Username: payload.Username,
			Roles:    payload.Roles,
		}
	}

	return server.WithAuth(map[string]auth.Authenticator{
		"jwt":   di.JWTAuth(),
		"token": auth.NewToken(tokenItems),
	})
}
