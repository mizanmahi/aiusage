package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mizanmahi/aiusage/server/internal/config"
	"github.com/mizanmahi/aiusage/server/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.Env)
	slog.SetDefault(logger)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: newRouter(logger),
	}

	logger.Info("server starting", "addr", server.Addr, "env", cfg.Env)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func newRouter(logger *slog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, "ok")
	})

	return router
}

func newLogger(env string) *slog.Logger {
	if env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
