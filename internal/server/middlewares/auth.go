package middlewares

import (
	"net/http"
	"strings"

	"rtc/internal/auth"
)

func JWTAuth(jwt *auth.JWT) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtToken := parseJWTToken(r)
			if jwtToken == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			p, err := jwt.Decode(jwtToken)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(
				auth.ToContext(r.Context(), p),
			))
		})
	}
}

func parseJWTToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const prefix = "jwt "
	if !strings.HasPrefix(authHeader, prefix) {
		return ""
	}

	token := authHeader[len(prefix):]
	if len(token) == 0 {
		return ""
	}

	return token
}
