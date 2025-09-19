package middlewares

import (
	"net/http"

	"github.com/DesSolo/rtc/internal/auth"
)

// Authorize ...
func Authorize(authorizer auth.Authorizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			payload := auth.FromContext(ctx)

			input := map[string]any{
				"method": r.Method,
				"path":   r.URL.Path,
				"user": map[string]any{
					"username": payload.Username,
					"roles":    payload.Roles,
				},
			}

			if err := authorizer.Authorize(ctx, input); err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
