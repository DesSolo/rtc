package server

import (
	"time"

	"github.com/DesSolo/rtc/internal/auth"
)

// OptionFunc ...
type OptionFunc func(s *Server)

// Noop ...
func Noop() OptionFunc {
	return func(_ *Server) {}
}

// WithAddress ...
func WithAddress(address string) OptionFunc {
	return func(s *Server) {
		s.address = address
	}
}

// WithReadHeaderTimeout ...
func WithReadHeaderTimeout(timeout time.Duration) OptionFunc {
	return func(s *Server) {
		s.readHeaderTimeout = timeout
	}
}

// WithAuth ...
func WithAuth(authenticators map[string]auth.Authenticator) OptionFunc {
	return func(s *Server) {
		s.auth = authenticators
	}
}

// WithAuthorizer ...
func WithAuthorizer(authorizer auth.Authorizer) OptionFunc {
	return func(s *Server) {
		s.authorizer = authorizer
	}
}
