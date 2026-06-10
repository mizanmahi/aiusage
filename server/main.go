package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/mizanmahi/aiusage/server/internal/config"
	"github.com/mizanmahi/aiusage/server/internal/handler"
	"github.com/mizanmahi/aiusage/server/internal/middleware"
	"github.com/mizanmahi/aiusage/server/internal/repository"
	"github.com/mizanmahi/aiusage/server/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.Env)
	slog.SetDefault(logger)

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("db open failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Error("db ping failed", "error", err)
		os.Exit(1)
	}

	events := repository.NewEventRepo(db)
	users := repository.NewUserRepo(db)
	projects := repository.NewProjectRepo(db)
	ingest := handler.NewIngestHandler(service.NewIngestService(events))
	admin := handler.NewAdminHandler(service.NewAdminService(users, projects))
	staticHandler, err := handler.NewStaticHandler(cfg.StaticDir)
	if err != nil {
		logger.Warn("static UI disabled", "dir", cfg.StaticDir, "error", err)
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: newRouter(logger, ingest, admin, users, staticHandler, cfg),
	}

	logger.Info("server starting", "addr", server.Addr, "env", cfg.Env)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func newRouter(
	logger *slog.Logger,
	ingest *handler.IngestHandler,
	admin *handler.AdminHandler,
	users repository.UserRepository,
	staticHandler *handler.StaticHandler,
	cfg *config.Config,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS(cfg.CORSOrigins))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "ok")
	})

	if ingest != nil && users != nil {
		router.With(
			middleware.Auth(users),
			minCLIVersionHeader(cfg.MinCLIVersion),
		).Post("/ingest", ingest.Create)
	}
	if admin != nil && users != nil {
		router.Route("/admin", func(r chi.Router) {
			r.Use(middleware.Auth(users))
			r.Get("/users", admin.Users)
			r.Post("/users", admin.CreateUser)
			r.Get("/users/{id}/breakdown", admin.UserBreakdown)
			r.Get("/users/{id}/summary", admin.UserUsageSummary)
			r.Get("/users/{id}", admin.UserProjects)
			r.Get("/summary", admin.Summary)
		})
	}
	if staticHandler != nil {
		router.NotFound(staticHandler.ServeHTTP)
	}

	return router
}

func minCLIVersionHeader(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if version != "" {
				w.Header().Set("X-Aiusage-Min-CLI-Version", version)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func newLogger(env string) *slog.Logger {
	if env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
