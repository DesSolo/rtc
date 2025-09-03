package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"rtc/internal/provider"
)

const (
	defaultAddress           = ":8080"
	defaultReadHeaderTimeout = time.Second * 3

	shutdownTimeout = time.Second * 5
)

// Server HTTP server
type Server struct {
	provider *provider.Provider

	address           string
	readHeaderTimeout time.Duration
	mux               chi.Router
}

// NewServer ...
func NewServer(provider *provider.Provider, options ...OptionFunc) *Server {
	s := &Server{
		provider:          provider,
		address:           defaultAddress,
		readHeaderTimeout: defaultReadHeaderTimeout,
		mux:               chi.NewRouter(),
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// Run start HTTP server
func (s *Server) Run(ctx context.Context) error {
	s.initRoutes()

	srv := &http.Server{
		Addr:        s.address,
		ReadTimeout: s.readHeaderTimeout,
		Handler:     s.mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
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

		r.Get("/audits", s.handleListAudits)
	})
}
