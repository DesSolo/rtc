package server

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	assets "rtc"
	"rtc/internal/auth"
	"rtc/internal/provider"
	"rtc/internal/server/middlewares"
)

const (
	defaultAddress           = ":8080"
	defaultReadHeaderTimeout = time.Second * 3

	shutdownTimeout = time.Second * 5
)

// Server HTTP server
type Server struct {
	provider *provider.Provider
	jwt      *auth.JWT
	auth     map[string]auth.Authenticator

	address           string
	readHeaderTimeout time.Duration
	mux               chi.Router
}

// NewServer ...
func NewServer(provider *provider.Provider, jwt *auth.JWT, options ...OptionFunc) *Server {
	s := &Server{
		provider: provider,
		jwt:      jwt,
		auth: map[string]auth.Authenticator{
			"jwt": jwt,
		},
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

	if err := s.initUI(); err != nil {
		slog.Warn("initUI", "err", err)
	}

	s.mux.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", s.handleLogin)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Authenticate(s.auth))

			r.Get("/projects", s.handleListProjects)
			r.Post("/projects", s.handleCreateProject)

			r.Get("/projects/{projectName}/envs", s.handleListEnvironments)

			r.Get("/projects/{projectName}/envs/{envName}/releases", s.handleListReleases)

			r.Delete("/projects/{projectName}/envs/{envName}/releases/{releaseName}", s.handleDeleteRelease)

			r.Get("/projects/{projectName}/envs/{envName}/releases/{releaseName}/configs", s.handleListConfigs)
			r.Put("/projects/{projectName}/envs/{envName}/releases/{releaseName}/configs", s.handleSetConfigValues)
			r.Post("/projects/{projectName}/envs/{envName}/releases/{releaseName}/configs", s.handleUpsertConfigs)

			r.Get("/audits", s.handleListAudits)

			r.Get("/users", s.handleListUsers)
			r.Post("/users", s.handleCreateUser)
			r.Patch("/users/{username}", s.handleUpdateUser)
		})

		r.Get("/health", s.handleHealth)
	})
}

func (s *Server) initUI() error {
	sub, err := fs.Sub(assets.FS, "frontend/ui/dist")
	if err != nil {
		return fmt.Errorf("fs.Sub: %w", err)
	}

	fileServer := http.FileServer(http.FS(sub))

	s.mux.Handle("/ui/*", http.StripPrefix("/ui", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, ".") {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})))

	return nil
}
