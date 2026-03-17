package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
	"github.com/Tokuchi61/Novascans/internal/platform/metrics"
	platformmiddleware "github.com/Tokuchi61/Novascans/internal/platform/middleware"
	moduleshared "github.com/Tokuchi61/Novascans/internal/platform/module"
)

type routeOptions struct {
	Config    config.Config
	Logger    *slog.Logger
	Metrics   *metrics.Registry
	Modules   []moduleshared.Module
	ReadyFunc func(context.Context) error
}

func buildRouter(options routeOptions) http.Handler {
	router := chi.NewRouter()

	router.Use(platformmiddleware.RequestID)
	router.Use(platformmiddleware.RealIP)
	router.Use(platformmiddleware.Recover(options.Logger))
	router.Use(platformmiddleware.Timeout(options.Config.HTTP.ReadTimeout))
	router.Use(platformmiddleware.Logger(options.Logger))
	router.Use(platformmiddleware.Metrics(options.Metrics))

	router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		platformhttp.WriteError(w, platformhttp.NotFound("resource not found", nil))
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		platformhttp.WriteError(w, platformhttp.MethodNotAllowed("method not allowed", nil))
	})

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		platformhttp.WriteData(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if options.ReadyFunc != nil {
			if err := options.ReadyFunc(r.Context()); err != nil {
				platformhttp.WriteError(w, platformhttp.ServiceUnavailable("service not ready", err))
				return
			}
		}

		platformhttp.WriteData(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	if options.Config.Metrics.Enabled && options.Metrics != nil {
		router.Handle("/metrics", options.Metrics)
	}

	router.Route("/api/v1", func(r chi.Router) {
		for _, mod := range options.Modules {
			mod.RegisterRoutes(r)
		}
	})

	return router
}
