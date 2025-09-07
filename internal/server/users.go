package server

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rtc/internal/models"
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

type user struct {
	Username  string    `json:"username"`
	IsEnabled bool      `json:"is_enabled"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
}

type listUsersResponse struct {
	Users []user `json:"users"`
	Total uint64 `json:"total"`
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query().Get("q")
	limit := queryOr[uint64](r, "limit", 10)
	offset := queryOr[uint64](r, "offset", 0)

	users, total, err := s.provider.ListUsers(ctx, q, limit, offset)
	if err != nil {
		slog.ErrorContext(ctx, "provider.ListUsers", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listUsersResponse{
		Users: convertModelsToUsers(users),
		Total: total,
	})
}

type createUserRequest struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	IsEnabled bool     `json:"is_enabled"`
	Roles     []string `json:"roles"`
}

func (r *createUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}

	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createUserRequest
	if err := bindJSON(r, &req); err != nil {
		slog.ErrorContext(ctx, "bindJSON", "err", err)
		respondError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	modelUser := &models.User{
		Username:  req.Username,
		IsEnabled: req.IsEnabled,
		Roles:     req.Roles,
	}

	if err := s.provider.CreateUser(ctx, modelUser, req.Password); err != nil {
		if errors.Is(err, provider.ErrAlreadyExists) {
			respondError(ctx, w, http.StatusConflict, "user already exists")
			return
		}

		slog.ErrorContext(ctx, "provider.CreateUser", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondStatus(w, http.StatusCreated)
}

type updateUserRequest struct {
	IsEnabled *bool    `json:"is_enabled,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")

	var req updateUserRequest
	if err := bindJSON(r, &req); err != nil {
		slog.ErrorContext(ctx, "bindJSON", "err", err)
		respondError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.provider.UpdateUser(ctx, username, &provider.UpdateUserFields{
		IsEnabled: req.IsEnabled,
		Roles:     req.Roles,
	}); err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			respondError(ctx, w, http.StatusNotFound, err.Error())
			return
		}

		slog.ErrorContext(ctx, "provider.ChangeUserEnabled", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondStatus(w, http.StatusOK)
}
