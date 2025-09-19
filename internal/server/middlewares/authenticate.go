package middlewares

import (
	"net/http"
	"strings"

	"github.com/DesSolo/rtc/internal/auth"
)

// Authenticate ...
func Authenticate(authenticators map[string]auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			kind, token := parseAuthorization(r)
			if kind == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			authAlgo, ok := authenticators[kind]
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			payload, err := authAlgo.Authenticate(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				auth.ToContext(r.Context(), payload),
			))
		})
	}
}

func parseAuthorization(r *http.Request) (string, string) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}
