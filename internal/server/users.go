package server

import (
	"errors"
	"log/slog"
	"net/http"

	"rtc/internal/provider"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req loginRequest
	if err := bindJSON(r, &req); err != nil {
		slog.ErrorContext(ctx, "bindJSON", "err", err)
		respondError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := s.provider.AuthenticateUser(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			respondError(ctx, w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		slog.ErrorContext(ctx, "provider.AuthenticateUser", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	token, err := s.auth.Encode(convertModelToAuth(user))
	if err != nil {
		slog.ErrorContext(ctx, "auth.Encode", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	respondData(ctx, w, http.StatusOK, loginResponse{
		Token: token,
	})
}
