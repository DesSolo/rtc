package auth

import "context"

type ctxKeyPayload struct{}

// ToContext ...
func ToContext(ctx context.Context, payload *Payload) context.Context {
	return context.WithValue(ctx, ctxKeyPayload{}, payload)
}

// FromContext ...
func FromContext(ctx context.Context) *Payload {
	payload, ok := ctx.Value(ctxKeyPayload{}).(*Payload)
	if !ok {
		return nil
	}

	return payload
}

// UsernameFromContext returns a username from context
func UsernameFromContext(ctx context.Context) string {
	payload := FromContext(ctx)
	if payload == nil {
		return ""
	}

	return payload.Username
}
