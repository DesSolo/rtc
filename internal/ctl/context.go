package ctl

import (
	"context"

	"rtc/internal/ctl/client"
)

type contextKeyClient struct{}

func clientToContext(ctx context.Context, client *client.Client) context.Context {
	return context.WithValue(ctx, contextKeyClient{}, client)
}
func clientFromContext(ctx context.Context) *client.Client {
	return ctx.Value(contextKeyClient{}).(*client.Client)
}
