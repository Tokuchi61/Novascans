package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	novapp "github.com/Tokuchi61/Novascans/internal/app"
	"github.com/Tokuchi61/Novascans/internal/platform/config"
	platformhttp "github.com/Tokuchi61/Novascans/internal/platform/http"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("config load failed", "error", err)
		os.Exit(1)
	}

	runtime, err := novapp.Bootstrap(context.Background(), cfg)
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("bootstrap failed", "error", err)
		os.Exit(1)
	}

	server := platformhttp.NewServer(cfg.HTTP, runtime.Router)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			runtime.Logger.Error("server shutdown failed", "error", err)
		}

		if err := runtime.Close(); err != nil {
			runtime.Logger.Error("runtime close failed", "error", err)
		}
	}()

	runtime.Logger.Info("http server starting", "addr", cfg.HTTP.Address())

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		runtime.Logger.Error("http server stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
