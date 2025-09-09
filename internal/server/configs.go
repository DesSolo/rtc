package server

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"rtc/internal/provider"
)

type configView struct {
	Enum string `json:"enum,omitempty"`
}

type config struct {
	Key       string     `json:"key"`
	ValueType string     `json:"value_type"`
	Value     string     `json:"value"`
	Group     string     `json:"group"`
	Usage     string     `json:"usage"`
	Writable  bool       `json:"writable"`
	View      configView `json:"view,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type listConfigsResponse struct {
	Configs []config `json:"configs"`
}

func (s *Server) handleListConfigs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectName := chi.URLParam(r, "projectName")
	envName := chi.URLParam(r, "envName")
	releaseName := chi.URLParam(r, "releaseName")

	configs, err := s.provider.Configs(ctx, projectName, envName, releaseName)
	if err != nil {
		slog.ErrorContext(ctx, "provider.Configs", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listConfigsResponse{
		Configs: convertModelsToConfigs(configs),
	})
}

type setConfigValuesRequest map[string]string

func (r setConfigValuesRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("empty values")
	}

	return nil
}

func (s *Server) handleSetConfigValues(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectName := chi.URLParam(r, "projectName")
	envName := chi.URLParam(r, "envName")
	releaseName := chi.URLParam(r, "releaseName")

	var req setConfigValuesRequest
	if err := bindJSON(r, &req); err != nil {
		slog.ErrorContext(ctx, "bindJSON", "err", err)
		respondError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.provider.SetConfigValues(ctx, projectName, envName, releaseName, convertValuesToModels(req)); err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			respondError(ctx, w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, provider.ErrNotValid) {
			respondError(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(ctx, "provider.SetConfigValues", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondStatus(w, http.StatusCreated)
}

type upsertConfigRequest []config

func (u upsertConfigRequest) Validate() error {
	for _, conf := range u {
		if conf.Key == "" {
			return errors.New("key is required")
		}

		if conf.ValueType == "" {
			return errors.New("value type is required")
		}

		if conf.Value == "" {
			return errors.New("value is required")
		}

		if conf.Usage == "" {
			return errors.New("usage is required")
		}
	}

	return nil
}

func (s *Server) handleUpsertConfigs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectName := chi.URLParam(r, "projectName")
	envName := chi.URLParam(r, "envName")
	releaseName := chi.URLParam(r, "releaseName")

	var req upsertConfigRequest
	if err := bindJSON(r, &req); err != nil {
		slog.ErrorContext(ctx, "bindJSON", "err", err)
		respondError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.provider.UpsertConfigs(ctx, projectName, envName, releaseName, convertConfigsToModels(req)); err != nil {
		if errors.Is(err, provider.ErrNotFound) {
			respondError(ctx, w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, provider.ErrNotValid) {
			respondError(ctx, w, http.StatusBadRequest, err.Error())
			return
		}

		slog.ErrorContext(ctx, "provider.UpsertConfigs", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondStatus(w, http.StatusCreated)
}
