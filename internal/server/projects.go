package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"rtc/internal/provider"
)

type project struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type listProjectsResponse struct {
	Projects []project `json:"projects"`
	Total    uint64    `json:"total"`
}

func (s *Server) handleListProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query().Get("q")
	limit := queryOr[uint64](r, "limit", 10)
	offset := queryOr[uint64](r, "offset", 0)

	projects, total, err := s.provider.Projects(ctx, q, limit, offset)
	if err != nil {
		slog.ErrorContext(ctx, "provider.Projects", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusOK, listProjectsResponse{
		Projects: convertModelsToProjects(projects),
		Total:    total,
	})
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r createProjectRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	if r.Description == "" {
		return fmt.Errorf("description is required")
	}

	return nil
}

type createProjectResponse struct {
	Project project `json:"project"`
}

func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req createProjectRequest
	if err := bindJSON(r, &req); err != nil {
		slog.ErrorContext(ctx, "bindJSON", "err", err)
		respondError(ctx, w, http.StatusBadRequest, err.Error())
		return
	}

	newProject, err := s.provider.CreateProject(ctx, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, provider.ErrAlreadyExists) {
			respondError(ctx, w, http.StatusConflict, err.Error())
			return
		}

		slog.ErrorContext(ctx, "provider.CreateProject", "err", err)
		respondError(ctx, w, http.StatusInternalServerError, err.Error())
		return
	}

	respondData(ctx, w, http.StatusCreated, createProjectResponse{
		Project: convertModelToProject(newProject),
	})
}
