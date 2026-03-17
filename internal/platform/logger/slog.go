package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func New(cfg config.LogConfig) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLevel(cfg.Level),
		}),
	)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
