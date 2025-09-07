package auth

// Authenticator ...
type Authenticator interface {
	Authenticate(token string) (*Payload, error)
}
