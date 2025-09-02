package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type release struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type listReleasesResponse struct {
	Releases []release `json:"releases"`
}

func (s *Server) handleListReleases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectName := chi.URLParam(r, "projectName")
	envName := chi.URLParam(r, "envName")

	releases, err := s.provider.ListReleases(ctx, projectName, envName)
	if err != nil {
		slog.ErrorContext(ctx, "provider.ListReleases",
			"projectName", projectName,
			"envName", envName,
			"err", err,
		)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listReleasesResponse{
		Releases: convertModelsToReleases(releases),
	})
}

func (s *Server) handleDeleteRelease(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectName := chi.URLParam(r, "projectName")
	envName := chi.URLParam(r, "envName")
	releaseName := chi.URLParam(r, "releaseName")

	if err := s.provider.DeleteRelease(ctx, projectName, envName, releaseName); err != nil {
		slog.ErrorContext(ctx, "provider.DeleteRelease", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondStatus(w, http.StatusNoContent)
}
