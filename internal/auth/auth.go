package auth

import "context"

// Authenticator ...
type Authenticator interface {
	Authenticate(token string) (*Payload, error)
}

// Authorizer ...
type Authorizer interface {
	Authorize(ctx context.Context, input map[string]any) error
}
