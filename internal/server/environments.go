package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type environment struct {
	Name string `json:"name"`
}

type listEnvironmentsResponse struct {
	Environments []environment `json:"environments"`
}

func (s *Server) handleListEnvironments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectName := chi.URLParam(r, "projectName")

	environments, err := s.provider.ListEnvironments(ctx, projectName)
	if err != nil {
		slog.ErrorContext(ctx, "provider.ListEnvironments", "projectName", projectName, "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listEnvironmentsResponse{
		Environments: convertModelsToEnvironments(environments),
	})
}
