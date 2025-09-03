package server

import "time"

// OptionFunc ...
type OptionFunc func(s *Server)

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
