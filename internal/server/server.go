package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"rtc/internal/provider"
)

type Server struct {
	provider *provider.Provider

	address string
	mux     chi.Router
}

func NewServer(provider *provider.Provider, address string) *Server {
	return &Server{
		provider: provider,
		address:  address,
		mux:      chi.NewRouter(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.initRoutes()

	srv := &http.Server{
		Addr:    s.address,
		Handler: s.mux,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown server: %v", err)
		}
	}()

	slog.InfoContext(ctx, "http server running", "address", s.address)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}

	return nil
}

func (s *Server) initRoutes() {
	s.mux.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	s.mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/projects", s.handleListProjects)
		r.Post("/projects", s.handleCreateProject)

		r.Get("/projects/{projectName}/envs", s.handleListEnvironments)

		r.Get("/projects/{projectName}/envs/{envName}/releases", s.handleListReleases)

		r.Delete("/projects/{projectName}/envs/{envName}/releases/{releaseName}", s.handleDeleteRelease)

		r.Get("/projects/{projectName}/envs/{envName}/releases/{releaseName}/configs", s.handleListConfigs)
		r.Put("/projects/{projectName}/envs/{envName}/releases/{releaseName}/configs", s.handleSetConfigValue)
		r.Post("/projects/{projectName}/envs/{envName}/releases/{releaseName}/configs", s.handleUpsertConfigs)
	})
}
