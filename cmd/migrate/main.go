package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/Tokuchi61/Novascans/internal/platform/config"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("config load failed", "error", err)
		os.Exit(1)
	}

	command := "status"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if err := run(context.Background(), cfg, command); err != nil {
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("migration command failed", "command", command, "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg config.Config, command string) error {
	db, err := sql.Open("pgx", cfg.DB.ConnectionString(cfg.DB.Name))
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	migrationsDir := filepath.Clean(filepath.Join("db", "migrations"))

	switch command {
	case "up":
		return goose.UpContext(ctx, db, migrationsDir)
	case "down":
		return goose.DownContext(ctx, db, migrationsDir)
	case "status":
		return goose.StatusContext(ctx, db, migrationsDir)
	default:
		return fmt.Errorf("unsupported command %q, expected one of: up, down, status", command)
	}
}
